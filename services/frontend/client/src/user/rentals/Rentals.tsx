// src/components/DashboardApp.tsx
import { useState } from "react";
import backend_constants from "../../backend_constants";

import type { Rental } from "./RentalDisplay";
import UserRentalDisplay from "./RentalDisplay";
import Query from "../../components/querying/Query";
import { Flow, FlowSwitch, FlowTopBackWrapper } from "../../components/navigation/Flow";
import Table from "../../components/table/Table";
import TableRow from "../../components/table/TableRow";

export default function UserRentalList() {
    const [rental, setRental] = useState<undefined | Rental>(undefined);


    const queryBuilder = (query: string) => { return { q: query } }

    return <div className="w-full h-full flex flex-col items-center">
        <Flow>
            <Query
                bodyBuilder={(rentals?: Rental[]) =>
                    <Table>
                        {rentals && rentals.map((rental) =>
                            <FlowSwitch direction="next" onClick={() => setRental(rental)}>
                                <TableRow>
                                    <h3 className="font-semibold text-lg">{rental.title}</h3>
                                    <p className="text-sm text-gray-600">{rental.description}</p>
                                </TableRow >
                            </FlowSwitch>
                        )}
                    </Table>
                }
                queryBuilder={queryBuilder}
                endpoint={`${backend_constants.address}/listings`}
            >
            </Query>
            {rental && <FlowTopBackWrapper><UserRentalDisplay rental={rental}></UserRentalDisplay></FlowTopBackWrapper>}
        </Flow>
    </div>
}
