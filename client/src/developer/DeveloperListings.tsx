// src/components/DashboardApp.tsx
import { useState } from "react";
import backend_constants from "../backend_constants";
import Search from "../components/Search";
import { ChevronLeft, ChevronRight, Plus } from "lucide-react";
import DeveloperListingDisplay from "./DeveloperListing";

interface SearchItem {
    ID: string | number;
    Title: string;
    Description: string;
    icon?: string;
    image?: string;
    Author: string;
}

export default function DeveloperListings() {
    const [listing, setListing] = useState<undefined | SearchItem>(undefined);

    if (listing) return <div className="flex flex-col w-full h-full">
        <div className="my-3 py-2 flex flex-row items-center hover:bg-gray-700 rounded-md transition-all" onClick={() => setListing(undefined)} >
            <ChevronLeft size={32} ></ChevronLeft>
            <label>Back</label>
        </div>
        <div className="flex-1">
            <DeveloperListingDisplay listing={listing}></DeveloperListingDisplay>
        </div>
    </div>

    const itemBuilder = (item: SearchItem) => {
        return <div
            key={item.ID}
            className="flex items-center gap-4 p-4 shadow-lg bg-gray-900 hover:bg-gray-800 transition-all hover:translate-x-2 mx-2"
        >
            <img
                src={item.icon || item.image || "/placeholder.svg"}
                alt=""
                className="w-12 h-12 object-cover rounded-xl"
            />
            <div className="flex flex-row justify-between items-center w-full" onClick={() => setListing(item)}>
                <div>
                    <h3 className="font-semibold text-lg">{item.Title}</h3>
                    <p className="text-sm text-gray-600">{item.Description}</p>
                </div>
                <div>
                    <ChevronRight size={24}></ChevronRight>
                </div>
            </div>
        </div>
    }

    const queryBuilder = (query: string) => { return { q: query } }

    const createListing = () => {
        const data = {
            title: "My New Listing",
            description: "This is a description of my listing."
        };

        fetch(`${backend_constants.address}/developer/listing`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            credentials: "include", // ✅ sends cookies/session along with the request
            body: JSON.stringify(data)
        })
            .then(response => {
                if (!response.ok) {
                    throw new Error(`HTTP error! status: ${response.status}`);
                }
                return response.json(); // parse JSON response
            })
            .then(result => {
                setListing(result)
            })
            .catch(error => {
                console.error("Error submitting listing:", error);
            });
    }

    return <div className="w-full h-full flex flex-col items-center">
        <Search itemBuilder={itemBuilder} queryBuilder={queryBuilder} endpoint={`${backend_constants.address}/developer/listings`}>

        </Search>
        <div className="flex flex-row justify-end w-full">
            <div className="p-5 dark:bg-gray-700 rounded-2xl dark:hover:bg-gray-600 transition-all" onClick={createListing}>
                <Plus size={32} />
            </div>
        </div>
    </div>
}
