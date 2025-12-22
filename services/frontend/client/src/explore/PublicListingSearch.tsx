import backend_constants from "../backend_constants";
import { PriceTag } from "../user/developer/prices/PriceTag";
import Query from "../components/querying/Query";
import TableRow from "../components/table/TableRow";
import Table from "../components/table/Table";
import { type Listing } from "../backend";
import { useNavigate } from "react-router-dom";
import { gotoListing } from "./PublicListingView";

export default function PublicListingSearch() {
    const navigate = useNavigate();

    const queryBuilder = (query: string) => {
        return { q: query };
    };

    return (
        <div className="w-full h-full max-w-7xl mx-auto px-4 py-6 flex flex-col gap-4 overflow-hidden">
            <Query
                endpoint={`${backend_constants.address}/listings`}
                queryBuilder={queryBuilder}
                bodyBuilder={(items?: Listing[]) => (
                    <Table className="flex flex-col gap-3 overflow-y-auto">
                        {items?.map((item) => (
                            <TableRow
                                key={item.id}
                                onClick={() => gotoListing(navigate, item)}
                                className="
                                    rounded-xl
                                    border
                                    border-slate-200
                                    dark:border-slate-800
                                    bg-white
                                    dark:bg-slate-900
                                    p-4
                                    hover:border-slate-400
                                    dark:hover:border-slate-600
                                    hover:shadow-md
                                    transition
                                    cursor-pointer
                                "
                            >
                                <div className="flex flex-col justify-between gap-2 min-w-0">
                                    <h3 className="text-lg font-semibold text-slate-900 dark:text-slate-200 truncate">
                                        {item.title}
                                    </h3>
                                    <p className="text-xs text-slate-600 dark:text-slate-500">
                                        By {item.author}
                                    </p>
                                </div>

                                <p className="text-sm text-slate-700 dark:text-slate-400 line-clamp-2">
                                    {item.description}
                                </p>

                                <div className="flex items-center justify-end">
                                    {item.price && (
                                        <PriceTag
                                            prices={[item.price]}
                                            rate="one time"
                                        />
                                    )}
                                </div>
                            </TableRow>
                        ))}
                    </Table>
                )}
            />
        </div>
    );
}
