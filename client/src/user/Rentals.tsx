// src/components/DashboardApp.tsx
import { useState } from "react";
import backend_constants from "../backend_constants";
import { ChevronLeft } from "lucide-react";
import ElementList from "../components/ElementList";
import ListElement from "../components/prefabs/ListElement";
import type { Rental } from "./RentalDisplay";
import UserRentalDisplay from "./RentalDisplay";

export default function UserRentalList() {
    const [rental, setRental] = useState<undefined | Rental>(undefined);

    const itemViewBuilder = (rental: Rental) => <div className="flex flex-col w-full h-full">
        <div className="my-3 py-2 flex flex-row items-center hover:bg-gray-700 rounded-md transition-all" onClick={() => setRental(undefined)} >
            <ChevronLeft size={32} ></ChevronLeft>
            <label>Back</label>
        </div>
        <div className="flex-1">
            <UserRentalDisplay rental={rental}></UserRentalDisplay>
        </div>
    </div>

    const itemBuilder = (item: Rental) => {
        return <ListElement onClick={() => setRental(item)}>
            <h3 className="font-semibold text-lg">{item.Title}</h3>
            <p className="text-sm text-gray-600">{item.Description}</p>
        </ListElement>
    }

    const queryBuilder = (query: string) => { return { q: query } }

    return <div className="w-full h-full flex flex-col items-center">
        <ElementList
            element={rental}
            itemBuilder={itemBuilder}
            queryBuilder={queryBuilder}
            endpoint={`${backend_constants.address}/user/rentals`}
            itemViewBuilder={itemViewBuilder}
        >

        </ElementList>
    </div>
}
