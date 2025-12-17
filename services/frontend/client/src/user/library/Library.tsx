import { useNavigate } from "react-router-dom";
import backend_constants from "../../backend_constants";

import Query from "../../components/querying/Query";
import Table from "../../components/table/Table";
import TableRow from "../../components/table/TableRow";
import type { Listing } from "../../backend";
import { gotoListing } from "../../explore/PublicListingView";

export default function UserLibrary() {
    const queryBuilder = (query: string) => { return { q: query } }
    const navigate = useNavigate()
    
    return <div className="w-full h-full flex flex-col items-center">
        <Query
            bodyBuilder={(entries?: Listing[]) =>
                <Table>
                    {entries && entries.map((entry) =>
                        <TableRow onClick={() => gotoListing(navigate,entry)}>
                            <h3 className="font-semibold text-lg">{entry.title}</h3>
                            <p className="text-sm text-gray-600">{entry.description}</p>
                        </TableRow >
                    )}
                </Table>
            }
            queryBuilder={queryBuilder}
            endpoint={`${backend_constants.address}/user/library`}
        ></Query>
    </div>
}
