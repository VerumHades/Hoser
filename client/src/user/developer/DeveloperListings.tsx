// src/components/DashboardApp.tsx
import { useState } from "react";
import backend_constants from "../../backend_constants";
import DeveloperListingDisplay from "./DeveloperListing";
import { API, type Listing } from "../../backend";

import { Flow, FlowSwitch } from "../../components/navigation/Flow";
import Query from "../../components/querying/Query";
import Table from "../../components/table/Table";
import TableRow from "../../components/table/TableRow";
import { Earth, EarthLock } from "lucide-react";


export default function DeveloperListings() {
    const [listing, setListing] = useState<undefined | Listing>(undefined);

    const queryBuilder = (query: string) => { return { q: query } }

    const createListing = async () => {
        const response = await API.developer.listing.create()
        if (response.ok) setListing(response.json)
    }

    return <div className="w-full h-full flex flex-col items-center min-h-0">
        <Flow>
            <Query
                bodyBuilder={(listings?: Listing[]) =>
                    <Table>
                        {listings && listings.map((listing) =>
                            <FlowSwitch direction="next" onClick={() => setListing(listing)}>
                                <TableRow>
                                    <h3 className="font-semibold text-lg">{listing.title}</h3>
                                    <p className="text-sm text-gray-600">{listing.description}</p>
                                    <div className="flex flex-row justify-between">
                                        {listing.accessMode == 1 ? <>Public < Earth /></> : <>Private < EarthLock /></>}
                                    </div>
                                </TableRow>
                            </FlowSwitch>)
                        }
                    </Table>
                }
                queryBuilder={queryBuilder}
                endpoint={`${backend_constants.address}/developer/listings`}
            />
            {listing ? <DeveloperListingDisplay listing={listing} onShouldClose={() => setListing(undefined)}></DeveloperListingDisplay> : <></>}
        </Flow>
    </div>
}
