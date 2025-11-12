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

    const itemViewBuilder = (rental: Rental) => <UserRentalDisplay rental={rental}></UserRentalDisplay>

    const itemBuilder = (item: Rental) => {
        return <ListElement onClick={() => setRental(item)}>
            <h3 className="font-semibold text-lg">{item.title}</h3>
            <p className="text-sm text-gray-600">{item.description}</p>
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
            onResetElement={() => setRental(undefined)}
        >

        </ElementList>
    </div>
}
