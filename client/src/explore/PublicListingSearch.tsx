import { useState } from "react";
import backend_constants from "../backend_constants";
import ElementList from "../components/ElementList";
import ListElement from "../components/prefabs/ListElement";
import { PriceTag } from "../components/prefabs/PriceTag";
import { TitleAndDescription } from "../components/prefabs/TitleAndDescription";
import { Flow, FlowSwitch, FlowTopBackWrapper } from "../components/Flow";
import Search from "../components/Search";
import { API } from "../backend";
import { useNavigate } from "react-router-dom";
interface CurrencyValue {
    Name: string,
    Short: string,
    Value: number
}

interface SearchItem {
    ID: string;
    Title: string;
    Description: string;
    MonthlyHardwarePrice?: CurrencyValue
    MonthlySubscriptionPrice?: CurrencyValue
    SinglePurchasePrice?: CurrencyValue
    author: string;
}

/*
<div
                            key={item.ID}
                            className="flex items-center gap-4 p-4 rounded-md shadow-lg bg-gray-900 mt-3"
                        >
                            <img
                                src={item.icon || item.image || "/placeholder.svg"}
                                alt=""
                                className="w-12 h-12 object-cover rounded-xl"
                            />
                            <div>
                                <h3 className="font-semibold text-lg">{item.Title}</h3>
                                <p className="text-sm text-gray-600">{item.Description}</p>
                                <p className="text-xs text-gray-400 mt-1">By {item.author}</p>
                            </div>
                        </div>

                    `${backend_constants.address}/rentals/public?${params}`
*/

interface PublicListingDisplayProps {
    item: SearchItem;
    onClick?: () => void;
}


function SinglePurchasePriceTag({ item, onClick }: PublicListingDisplayProps) {
    return item.SinglePurchasePrice ?
        <PriceTag onClick={onClick} prices={[{
            name: "One Time",
            currency_short: item.SinglePurchasePrice.Short,
            value: item.SinglePurchasePrice.Value
        }]}></PriceTag> : null
}
function MonthlySubscriptionPriceTag({ item, onClick }: PublicListingDisplayProps) {
    return item.MonthlySubscriptionPrice && item.MonthlyHardwarePrice ?
        <PriceTag onClick={onClick} prices={[{
            name: "Author",
            currency_short: item.MonthlySubscriptionPrice.Short,
            value: item.MonthlySubscriptionPrice.Value
        }, {
            name: "Hardware",
            currency_short: item.MonthlyHardwarePrice.Short,
            value: item.MonthlyHardwarePrice.Value
        }]} rate="monthly"></PriceTag> : null
}

interface PriceDisplayProps {
    onClickMonthly?: () => void,
    onClickOneTime?: () => void,
    item: SearchItem
}
function PriceDisplay({ item, onClickMonthly, onClickOneTime }: PriceDisplayProps) {
    return <div className="flex flex-row items-center text-slate-500 dark:text-slate-400 text-sm whitespace-nowrap md:mt-0 mt-5">
        <SinglePurchasePriceTag item={item} onClick={onClickOneTime}></SinglePurchasePriceTag>
        {
            item.SinglePurchasePrice &&
            item.MonthlySubscriptionPrice &&
            item.MonthlyHardwarePrice && <label className="mx-2">or</label>
        }
        <MonthlySubscriptionPriceTag item={item} onClick={onClickMonthly}></MonthlySubscriptionPriceTag>
    </div>
}

function PublicListingDisplay({ item }: PublicListingDisplayProps) {
    const navigate = useNavigate()
    const handleRent = () => {
        API.user.rentListing("" + item.ID)
        navigate("/account", { state: { dashpage: "My Rentals" } })
    }

    return (
        <div className="w-full h-full flex flex-col p-6">
            <div className="flex flex-row justify-between items-end">
                <TitleAndDescription
                    title={item.Title}
                    description={item.Description}
                    titleClassname="text-7xl">

                </TitleAndDescription>
                <div className="text-slate-500 dark:text-slate-400 text-sm">
                    <label className="">Purchase:</label>
                    <PriceDisplay item={item} onClickMonthly={handleRent} onClickOneTime={handleRent}></PriceDisplay>
                </div>
            </div>
            <div className="flex flex-col">

            </div>
        </div>
    );
};

export default function PublicListingSearch() {
    const [item, setItem] = useState<undefined | SearchItem>(undefined);

    const itemBuilder = (item: SearchItem) => {

        return <FlowSwitch direction="next">
            <ListElement onClick={() => setItem(item)}>
                <div className="flex flex-col md:flex-row w-full">
                    <div className="flex-1">
                        <h3 className="font-semibold text-lg">{item.Title}</h3>
                        <p className="text-sm">{item.Description}</p>
                        <p className="text-xs mt-1">By {item.author}</p>
                    </div>
                    <PriceDisplay item={item}></PriceDisplay>
                </div>
            </ListElement>
        </FlowSwitch>
    }

    const queryBuilder = (query: string) => { return { q: query } }

    return <div className="w-full h-full flex flex-col items-center bg-slate-100 dark:bg-slate-950">
        <Flow>
            <Search itemBuilder={itemBuilder} queryBuilder={queryBuilder} endpoint={`${backend_constants.address}/rentals/public`} />
            <FlowTopBackWrapper>
                {item ? <PublicListingDisplay item={item}></PublicListingDisplay> : <></>}
            </FlowTopBackWrapper>
            <div></div>
        </Flow>
    </div>
}
