// src/components/DashboardApp.tsx
import { useState } from "react";
import backend_constants from "../../backend_constants";
import DeveloperListingDisplay from "./DeveloperListing";
import { API, type DeveloperListing } from "../../backend";

import { Flow, FlowSwitch } from "../../components/navigation/Flow";
import Query from "../../components/querying/Query";
import Table from "../../components/table/Table";
import TableRow from "../../components/table/TableRow";
import { Earth, EarthLock } from "lucide-react";


export default function DeveloperListings() {
    const [listing, setListing] = useState<undefined | DeveloperListing>(undefined);

    const queryBuilder = (query: string) => { return { q: query } }

    const createListing = async () => {
        const response = await API.developer.listing.create()
        if (response.ok) setListing(response.json)
    }

    return <Flow className="flex flex-col flex-1 min-h-0 p-3">
        <>
            <FlowSwitch direction="next" onClick={() => createListing()}>
                New Listing
            </FlowSwitch>
            <Query
                bodyBuilder={(listings?: DeveloperListing[]) =>
                    <Table className="flex-1 overflow-y-auto min-h-0">
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
                className="flex-1 min-h-0"
                endpoint={`${backend_constants.address}/developer/listings`}
            />
        </>

        {listing ? <DeveloperListingDisplay listing={listing} onShouldClose={() => setListing(undefined)}></DeveloperListingDisplay> : <></>}
    </Flow>
}
