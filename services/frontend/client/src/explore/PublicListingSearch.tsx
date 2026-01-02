import { useNavigate } from "react-router-dom";
import { PublicListingRow } from "./List/PublicListingRow";
import { PublicListingCard } from "./List/PublicListingCard";
import { ListingAPI, type Listing } from "../backend/repositories/listing";
import { CursorPaginatedCollection } from "../components/view/CursorPaginatedCollection";

interface PublicListingSearchProps {
    onSelect?: (listing: Listing) => void;
}

export default function PublicListingSearch({
    onSelect
}: PublicListingSearchProps) {
    const navigate = useNavigate();

    return (
        <div className="w-full max-w-4xl min-h-0 overflow-hidden mt-5 px-6">
            <CursorPaginatedCollection<Listing>
                fetchPage={cursor => ListingAPI.search(cursor)}
                renderItemRow={listing => (
                    <PublicListingRow
                        key={listing.id}
                        listing={listing}
                        onSelect={
                            onSelect
                                ? () => onSelect(listing)
                                : () => navigate(`/listing/${listing.id}`)
                        }
                    />
                )}
                renderItemCard={listing => (
                    <PublicListingCard
                        key={listing.id}
                        listing={listing}
                        onSelect={
                            onSelect
                                ? () => onSelect(listing)
                                : () => navigate(`/listing/${listing.id}`)
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
    );
}