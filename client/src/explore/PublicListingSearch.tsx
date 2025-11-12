import { useState } from "react";
import backend_constants from "../backend_constants";
import ElementList from "../components/ElementList";
import ListElement from "../components/prefabs/ListElement";
import { PriceTag } from "../components/prefabs/PriceTag";
import { TitleAndDescription } from "../components/prefabs/TitleAndDescription";

interface CurrencyValue {
    Name: string,
    Short: string,
    Value: number
}

interface SearchItem {
    ID: string | number;
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
}
function CreateSinglePurchasePriceTag(item: SearchItem) {
    return item.SinglePurchasePrice ?
        <PriceTag prices={[{
            name: "One Time",
            currency_short: item.SinglePurchasePrice.Short,
            value: item.SinglePurchasePrice.Value
        }]}></PriceTag> : undefined
}
function CreateMonthlySubscriptionPriceTag(item: SearchItem) {
    return item.MonthlySubscriptionPrice && item.MonthlyHardwarePrice ?
        <PriceTag prices={[{
            name: "Author",
            currency_short: item.MonthlySubscriptionPrice.Short,
            value: item.MonthlySubscriptionPrice.Value
        }, {
            name: "Hardware",
            currency_short: item.MonthlyHardwarePrice.Short,
            value: item.MonthlyHardwarePrice.Value
        }]} rate="monthly"></PriceTag> : <></>
}

function PublicListingDisplay({ item }: PublicListingDisplayProps) {
    const singlePurchaseTag = CreateSinglePurchasePriceTag(item)
    const monthlyPurchaseTag = CreateMonthlySubscriptionPriceTag(item)

    return (
        <div className="w-full h-full flex flex-col p-6 shadow-md rounded-lg">
            <div>
                <TitleAndDescription
                    title={item.Title}
                    description={item.Description}
                    titleClassname="text-7xl">
                    
                </TitleAndDescription>
            </div>
            <div className="flex flex-col">
                <h2 className="text-2xl font-bold mb-4 text-gray-900 dark:text-gray-100">
                    Purchase:
                </h2>
                <div className="flex flex-col">
                    { singlePurchaseTag ? <ListElement>
                        <div className="flex flex-row items-center">
                            <TitleAndDescription
                                className="flex-1"
                                title={"One time"}
                                description={"Pay once, rent only hardware"}>
                            </TitleAndDescription>
                            <div>{singlePurchaseTag}</div>
                        </div>
                    </ListElement> : <></>}
                    {monthlyPurchaseTag ? <ListElement>
                        <div className="flex flex-row items-center">
                            <TitleAndDescription
                                className="flex-1"
                                title={"Monthly Subscription"}
                                description={"Pay monthly rent, with hardware"}>
                            </TitleAndDescription>
                            <div>{monthlyPurchaseTag}</div>
                        </div>
                    </ListElement> : <></>}
                </div>
            </div>
        </div>
    );
};


export default function PublicListingSearch() {
    const [item, setItem] = useState<undefined | SearchItem>(undefined);

    const itemViewBuilder = (item: SearchItem) => <PublicListingDisplay item={item} />

    const itemBuilder = (item: SearchItem) => {
        const singlePurchaseTag = CreateSinglePurchasePriceTag(item)
        const monthlyPurchaseTag = CreateMonthlySubscriptionPriceTag(item)

        const showOr = singlePurchaseTag && monthlyPurchaseTag;

        return <ListElement onClick={() => setItem(item)}>
            <div className="flex flex-col md:flex-row w-full">
                <div className="flex-1">
                    <h3 className="font-semibold text-lg">{item.Title}</h3>
                    <p className="text-sm">{item.Description}</p>
                    <p className="text-xs mt-1">By {item.author}</p>
                </div>
                <div className="flex flex-row items-center text-slate-500 dark:text-slate-400 text-sm whitespace-nowrap md:mt-0 mt-5">
                    {singlePurchaseTag}{showOr ? <label className="mx-5">or</label> : <></>}{monthlyPurchaseTag}
                </div>
            </div>
        </ListElement>
    }

    const queryBuilder = (query: string) => { return { q: query } }

    return <div className="w-full h-full flex flex-col items-center bg-slate-100">
        <ElementList
            element={item}
            itemBuilder={itemBuilder}
            queryBuilder={queryBuilder}
            endpoint={`${backend_constants.address}/rentals/public`}
            itemViewBuilder={itemViewBuilder}
            onResetElement={() => setItem(undefined)}
        >

        </ElementList>
    </div>
}
