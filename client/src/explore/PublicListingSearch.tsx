import { useState } from "react";
import backend_constants from "../backend_constants";
import { PriceTag } from "../user/developer/prices/PriceTag";
import { TitleAndDescription } from "../components/prefabs/TitleAndDescription";
import { Flow, FlowSwitch, FlowTopBackWrapper } from "../components/navigation/Flow";
import { API, BillingFrequency, BillingFrequencyName } from "../backend";
import { useNavigate } from "react-router-dom";
import Query from "../components/querying/Query";
import TableRow from "../components/table/TableRow";
import Table from "../components/table/Table";

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


interface PriceDisplayProps {
    item: SearchItem;
}

function PriceDisplay({ item }: PriceDisplayProps) {
    return <div className="flex flex-col gap-3">
        {
            item.prices && Object.values(item.prices).map(price => <PriceTag
                prices={[price.currency]} rate={BillingFrequencyName[price.type as BillingFrequency]} />)
        }
    </div>
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
                    <PriceDisplay item={item}></PriceDisplay>
                </div>
            </div>
            <div className="flex flex-col">

            </div>
        </div>
    );
};

export default function PublicListingSearch() {
    const [item, setItem] = useState<undefined | SearchItem>(undefined);


    const queryBuilder = (query: string) => { return { q: query } }

    return <div className="w-full h-full max-w-7xl justify-center flex flex-col items-center">
        <Flow>
            <Query
                bodyBuilder={(items?: SearchItem[]) =>
                    <Table className="overflow-y-auto">
                        {items && items.map(item => <FlowSwitch direction="next" onClick={() => setItem(item)}>
                            <TableRow>
                                <div className="flex flex-col justify-between h-full">
                                    <h3 className="font-semibold text-lg text-slate-800 dark:text-slate-300">{item.title}</h3>
                                    <p className="text-xs mt-1 text-slate-700 dark:text-slate-500">By {item.author}</p>
                                </div>
                                <p className="text-sm text-slate-500">{item.description}</p>
                                <PriceDisplay item={item}></PriceDisplay>
                            </TableRow>
                        </FlowSwitch>)}
                    </Table>
                }
                queryBuilder={queryBuilder}
                endpoint={`${backend_constants.address}/rentals/public`}
            />
            <FlowTopBackWrapper>
                {item ? <PublicListingDisplay item={item}></PublicListingDisplay> : <></>}
            </FlowTopBackWrapper>
            <div></div>
        </Flow>
    </div>
}
