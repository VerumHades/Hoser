
import backend_constants from "../backend_constants";
import { PriceTag } from "../user/developer/prices/PriceTag";

//import { useNavigate } from "react-router-dom";
import Query from "../components/querying/Query";
import TableRow from "../components/table/TableRow";
import Table from "../components/table/Table";

import {type Listing } from "../backend";
import { useNavigate } from "react-router-dom";
import { gotoListing } from "./PublicListingView";


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

export default function PublicListingSearch() {
    const queryBuilder = (query: string) => { return { q: query } }
    const navigate = useNavigate()

    return <div className="w-full h-full max-w-7xl justify-center flex flex-col items-stretch overflow-hidden">
        <Query
            bodyBuilder={(items?: Listing[]) =>
                <Table className="overflow-y-auto">
                    {items && items.map(item => 
                        <TableRow onClick={() => gotoListing(navigate,item)}>
                            <div className="flex flex-col justify-between h-full">
                                <h3 className="font-semibold text-lg text-slate-800 dark:text-slate-300">{item.title}</h3>
                                <p className="text-xs mt-1 text-slate-700 dark:text-slate-500">By {item.author}</p>
                            </div>
                            <p className="text-sm text-slate-500">{item.description}</p>
                            {item.price && <PriceTag prices={[item.price]} rate={"one time"} />}
                        </TableRow>
                    )}
                </Table>
            }
            queryBuilder={queryBuilder}
            endpoint={`${backend_constants.address}/listings`}
        />
    </div>
}
