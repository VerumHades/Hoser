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
    name: string,
    short: string,
    value: number
}

interface HardwareDTO {
    cpu: number;
    ram: number;
    disk: number;
}

interface ListingPrices {
    singlePurchase?: CurrencyValue;
    monthlySubscription?: CurrencyValue;
    monthlyHardware?: CurrencyValue;
}

interface SearchItem {
    id: string;
    title: string;
    description: string;
    images?: string[];
    hardware?: HardwareDTO;
    prices?: ListingPrices;
    author?: string;
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
    const price = item.prices?.singlePurchase;
    if (!price) return null;

    return (
        <PriceTag
            onClick={onClick}
            prices={[{
                name: "One Time",
                currency_short: price.short,
                value: price.value
            }]}
        />
    );
}

function MonthlySubscriptionPriceTag({ item, onClick }: PublicListingDisplayProps) {
    const subscription = item.prices?.monthlySubscription;
    const hardware = item.prices?.monthlyHardware;

    if (!subscription || !hardware) return null;

    return (
        <PriceTag
            onClick={onClick}
            prices={[
                {
                    name: "Author",
                    currency_short: subscription.short,
                    value: subscription.value
                },
                {
                    name: "Hardware",
                    currency_short: hardware.short,
                    value: hardware.value
                }
            ]}
            rate="monthly"
        />
    );
}

interface PriceDisplayProps {
    onClickMonthly?: () => void;
    onClickOneTime?: () => void;
    item: SearchItem;
}

function PriceDisplay({ item, onClickMonthly, onClickOneTime }: PriceDisplayProps) {
    const hasOneTime = !!item.prices?.singlePurchase;
    const hasMonthly = !!item.prices?.monthlySubscription && !!item.prices?.monthlyHardware;

    return (
        <div className="flex flex-row items-center text-slate-500 dark:text-slate-400 text-sm whitespace-nowrap md:mt-0 mt-5">
            {hasOneTime && <SinglePurchasePriceTag item={item} onClick={onClickOneTime} />}
            {hasOneTime && hasMonthly && <label className="mx-2">or</label>}
            {hasMonthly && <MonthlySubscriptionPriceTag item={item} onClick={onClickMonthly} />}
        </div>
    );
}

interface PublicListingDisplayProps {
    item: SearchItem;
}


function PublicListingDisplay({ item }: PublicListingDisplayProps) {
    const navigate = useNavigate()
    const handleRent = () => {
        API.user.rentListing("" + item.id)
        navigate("/account", { state: { dashpage: "My Rentals" } })
    }

    return (
        <div className="w-full h-full flex flex-col p-6">
            <div className="flex flex-row justify-between items-end">
                <TitleAndDescription
                    title={item.title}
                    description={item.description}
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
        console.log(item)
        return <FlowSwitch direction="next">
            <ListElement onClick={() => setItem(item)}>
                <div className="flex flex-col md:flex-row w-full">
                    <div className="flex-1">
                        <h3 className="font-semibold text-lg">{item.title}</h3>
                        <p className="text-sm">{item.description}</p>
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
