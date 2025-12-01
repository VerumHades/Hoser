// src/components/DashboardApp.tsx
import { useState } from "react";
import backend_constants from "../backend_constants";
import ListElement from "../components/prefabs/ListElement";
import type { Rental } from "./RentalDisplay";
import UserRentalDisplay from "./RentalDisplay";
import Query from "../components/querying/Query";
import { Flow, FlowSwitch } from "../components/Flow";

export default function UserRentalList() {
    const [rental, setRental] = useState<undefined | Rental>(undefined);


    const queryBuilder = (query: string) => { return { q: query } }

    return <div className="w-full h-full flex flex-col items-center">
        <Flow>
            <Query
                bodyBuilder={(rentals?: Rental[]) =>
                    <>
                        {
                            rentals && rentals.map((rental) =>
                                <ListElement onClick={() => setRental(rental)}>
                                    <FlowSwitch direction="next">
                                        <h3 className="font-semibold text-lg">{rental.title}</h3>
                                        <p className="text-sm text-gray-600">{rental.description}</p>
                                    </FlowSwitch>
                                </ListElement>
                            )
                        }
                    </>
                }
                queryBuilder={queryBuilder}
                endpoint={`${backend_constants.address}/user/rentals`}
            >
            </Query>
            {rental && <UserRentalDisplay rental={rental}></UserRentalDisplay>}
        </Flow>
    </div>
}
