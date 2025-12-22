import { useNavigate } from "react-router-dom";
import backend_constants from "../../backend_constants";

import Query from "../../components/querying/Query";
import Table from "../../components/table/Table";
import TableRow from "../../components/table/TableRow";
import { gotoListing } from "../../explore/PublicListingView";
import { useState } from "react";
import { type Listing } from "../../backend";
import { ExternalLink, Play } from "lucide-react";

export default function UserLibrary() {
    const queryBuilder = (query: string) => ({ q: query });
    const navigate = useNavigate();
    const [loadingInstanceId, setLoadingInstanceId] = useState<string | null>(null);

    const handleCreateInstance = (listingId: string) => {
        navigate(`/dashboard/create-instance/${listingId}`);
    };

    return (
        <Query
            className="w-full h-full flex flex-col items-center"
            queryBuilder={queryBuilder}
            endpoint={`${backend_constants.address}/user/library`}
            bodyBuilder={(entries?: Listing[]) => (
                <Table>
                    {entries?.map((entry) => (
                        <TableRow key={entry.id} className="flex justify-between items-center px-4 py-2 hover:bg-gray-50 rounded">
                            <div className="flex-1 cursor-pointer" onClick={() => gotoListing(navigate, entry)}>
                                <h3 className="font-semibold text-lg">{entry.title}</h3>
                                <p className="text-sm text-gray-600">{entry.description}</p>
                            </div>
                            <div className="flex space-x-2">
                                <button
                                    className="p-2 rounded hover:bg-gray-200"
                                    onClick={() => gotoListing(navigate, entry)}
                                >
                                    <ExternalLink size={20}/>
                                </button>
                                <button
                                    className="p-2 rounded hover:bg-gray-200 disabled:opacity-50"
                                    onClick={() => handleCreateInstance(entry.id)}
                                    disabled={loadingInstanceId === entry.id}
                                >
                                    <Play size={20} />
                                </button>
                            </div>
                        </TableRow>
                    ))}
                </Table>
            )}
        />
    );
}
