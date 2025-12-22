import { useCallback, useState } from "react";
import { useNavigate } from "react-router-dom";
import backend_constants from "../backend_constants";
import { API, type Listing } from "../backend";
import Query from "../components/querying/Query";
import { CollectionView } from "../components/view/CollectionView";
import { FlowSwitch } from "../components/navigation/Flow";

type ListingViewMode = "cards" | "table";

import { PublicListingRow } from "./PublicListinRow";
import { PublicListingCard } from "./PublicListingCard";

interface PublicListingSearchProps {
    onSelect?: (listing: Listing) => void;
}

export default function PublicListingSearch({ onSelect }: PublicListingSearchProps) {
    const navigate = useNavigate();
    const [viewMode, setViewMode] = useState<ListingViewMode>("cards");

    const renderListing = useCallback(
        (listings: Listing[] | undefined) => (
            <CollectionView
                items={listings}
                viewMode={viewMode}
                renderTableRow={(listing) => (
                    <PublicListingRow
                        key={listing.id}
                        listing={listing}
                        onSelect={onSelect ? () => onSelect(listing) : () => navigate(`/listing/${listing.id}`)}
                    />
                )}
                renderCard={(listing) => (
                    <PublicListingCard
                        key={listing.id}
                        listing={listing}
                        onSelect={onSelect ? () => onSelect(listing) : () => navigate(`/listing/${listing.id}`)}
                    />
                )}
                emptyState={<p className="text-sm text-slate-500">No listings found</p>}
            />
        ),
        [viewMode, onSelect, navigate]
    );

    const queryBuilder = useCallback((queryText: string) => ({ q: queryText }), []);

    return (
        <div className="flex flex-col w-full h-full max-w-7xl mx-auto px-4 py-6 gap-4">
            <div className="flex items-center justify-between mb-4">
                <FlowSwitch direction="next">
                    {/* Optionally attach a callback for creating a listing if needed */}
                    Refresh
                </FlowSwitch>
                <div className="flex gap-2">
                    <button
                        onClick={() => setViewMode("cards")}
                        className={viewMode === "cards" ? "font-semibold" : ""}
                    >
                        Cards
                    </button>
                    <button
                        onClick={() => setViewMode("table")}
                        className={viewMode === "table" ? "font-semibold" : ""}
                    >
                        Table
                    </button>
                </div>
            </div>

            <Query<Listing[]>
                endpoint={`${backend_constants.address}/listings`}
                queryBuilder={queryBuilder}
                className="flex-1 min-h-0 overflow-hidden"
                render={renderListing}
            />
        </div>
    );
}
