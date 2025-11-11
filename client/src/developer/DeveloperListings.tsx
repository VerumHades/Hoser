// src/components/DashboardApp.tsx
import { useState } from "react";
import backend_constants from "../backend_constants";
import { ChevronLeft, Plus } from "lucide-react";
import DeveloperListingDisplay, { type Listing } from "./DeveloperListing";
import ElementList from "../components/ElementList";
import { API } from "../backend";
import ListElement from "../components/prefabs/ListElement";


export default function DeveloperListings() {
    const [listing, setListing] = useState<undefined | Listing>(undefined);

    const itemViewBuilder = (listing: Listing) => 
        <DeveloperListingDisplay listing={listing} onShouldClose={() => setListing(undefined)}></DeveloperListingDisplay>

    const itemBuilder = (item: Listing) => {
        return <ListElement onClick={() => setListing(item)}>
            <h3 className="font-semibold text-lg">{item.Title}</h3>
            <p className="text-sm text-gray-600">{item.Description}</p>
        </ListElement>
    }

    const queryBuilder = (query: string) => { return { q: query } }

    const createListing = async () => {
        const response = await API.developer.createListing()
        if(response.ok) setListing(response.json)
    }

    return <div className="w-full h-full flex flex-col items-center">
        <ElementList
            element={listing}
            itemBuilder={itemBuilder} 
            queryBuilder={queryBuilder} 
            endpoint={`${backend_constants.address}/developer/listings`}
            itemViewBuilder={itemViewBuilder}
            onResetElement={() => setListing(undefined)}
        >

        </ElementList>
        <div className="flex flex-row justify-end w-full">
            <div className="p-5 bg-slate-300 hover:bg-slate-400 dark:bg-gray-700 rounded-2xl dark:hover:bg-gray-600 transition-all" onClick={createListing}>
                <Plus size={32} />
            </div>
        </div>
    </div>
}
