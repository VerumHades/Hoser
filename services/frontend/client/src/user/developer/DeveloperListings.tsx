import { DeveloperListingAPI, type DeveloperListing } from "../../backend/repositories/developer_listing";
import { CursorPaginatedCollection } from "../../components/view/CursorPaginatedCollection";

import { useNavigate } from "react-router-dom"
import { DeveloperListingRow } from "./List/DeveloperListingRow";
import { DeveloperListingCard } from "./List/DeveloperListingCard";

export default function DeveloperListings() {
    const navigate = useNavigate()

    return (
        <div className="flex-1 min-h-0 flex flex-col items-center">
            <div className="w-full max-w-4xl min-h-0 overflow-hidden mt-5 px-6">
                <CursorPaginatedCollection<DeveloperListing>
                    fetchPage={cursor => DeveloperListingAPI.list(cursor)}
                    renderItemRow={listing => (
                        <DeveloperListingRow
                            key={listing.id}
                            listing={listing}
                            onSelect={
                                () => navigate(`/dashboard/developer/listing/${listing.id}`)
                            }
                        />
                    )}
                    renderItemCard={listing => (
                        <DeveloperListingCard
                            key={listing.id}
                            listing={listing}
                            onSelect={
                                () => navigate(`/dashboard/developer/listing/${listing.id}`)
                            }
                        />
                    )}
                    emptyState={
                        <p className="text-sm text-slate-500">
                            No listings found
                        </p>
                    }
                />
            </div>
        </div>
        
    );
}
