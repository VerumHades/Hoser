// src/components/DashboardApp.tsx
import { useState } from "react";
import backend_constants from "../backend_constants";
import { ChevronLeft, Plus } from "lucide-react";
import DeveloperListingDisplay from "./DeveloperListing";
import { API, type Listing } from "../backend";
import ListElement from "../components/prefabs/ListElement";
import { Flow, FlowSwitch } from "../components/Flow";
import Search from "../components/Search";


export default function DeveloperListings() {
    const [listing, setListing] = useState<undefined | Listing>(undefined);

    const itemBuilder = (item: Listing) => {
        return <FlowSwitch direction="next">
            <ListElement onClick={() => setListing(item)}>
                <h3 className="font-semibold text-lg">{item.title}</h3>
                <p className="text-sm text-gray-600">{item.description}</p>
            </ListElement>
        </FlowSwitch>
    }

    const queryBuilder = (query: string) => { return { q: query } }

    const createListing = async () => {
        const response = await API.developer.createListing()
        if (response.ok) setListing(response.json)
    }

    return <div className="w-full h-full flex flex-col items-center min-h-0">
        <Flow>
            <Search itemBuilder={itemBuilder} queryBuilder={queryBuilder} endpoint={`${backend_constants.address}/developer/listings`} />
            {listing ? <DeveloperListingDisplay listing={listing} onShouldClose={() => setListing(undefined)}></DeveloperListingDisplay> : <></>}
        </Flow>

    </div>
}
