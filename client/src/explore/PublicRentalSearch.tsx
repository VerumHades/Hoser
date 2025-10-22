import { useEffect, useState, useCallback } from "react";
import backend_constants from "../backend_constants";
import debounce from "lodash.debounce"
import Search from "../Search";

interface SearchItem {
    ID: string | number;
    Title: string;
    Description: string;
    icon?: string;
    image?: string;
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

export default function PublicRentalSearch() {
    const itemBuilder = (item: SearchItem) => {
        return <div
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
    }

    const queryBuilder = (query: string) => {return {q: query}}

    return <>
        <Search itemBuilder={itemBuilder} queryBuilder={queryBuilder} endpoint={`${backend_constants.address}/rentals/public`}>
            
        </Search>
    </> 
}
