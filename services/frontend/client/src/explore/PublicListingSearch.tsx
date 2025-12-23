import { useCallback, useState } from "react";
import { useNavigate } from "react-router-dom";
import backend_constants from "../backend_constants";
import { API, type Listing } from "../backend";
import Query from "../components/querying/Query";

import { PublicListingRow } from "./List/PublicListingRow";
import { PublicListingCard } from "./List/PublicListingCard";
import { CollectionViewContainer, type CollectionViewMode } from "../components/view/CollectionViewContainer";

interface PublicListingSearchProps {
    onSelect?: (listing: Listing) => void;
}

export default function PublicListingSearch({ onSelect }: PublicListingSearchProps) {
    const navigate = useNavigate();
    const [viewMode, setViewMode] = useState<CollectionViewMode>("cards")
    
    const renderListing = useCallback(
        (listings: Listing[] | undefined) => (
            <CollectionViewContainer
                items={listings}
                onViewModeChange={(mode) => setViewMode(mode)}
                renderTableRow={(listing) => (
                    <PublicListingRow
                        key={listing.id}
                        listing={listing}
                        onSelect={onSelect ? () => onSelect(listing) : () => navigate(`/listing/${listing.id}`)} />
                )}
                renderCard={(listing) => (
                    <PublicListingCard
                        key={listing.id}
                        listing={listing}
                        onSelect={onSelect ? () => onSelect(listing) : () => navigate(`/listing/${listing.id}`)} />
                )}
                emptyState={<p className="text-sm text-slate-500">No listings found</p>} viewMode={viewMode}          />
        ),
        [viewMode, onSelect, navigate]
    );

    const queryBuilder = useCallback((queryText: string) => ({ q: queryText }), []);

    return (
        <Query<Listing[]>
            endpoint={`${backend_constants.address}/listings`}
            queryBuilder={queryBuilder}
            className="w-full max-w-4xl min-h-0 overflow-hidden mt-5 px-6"
            render={renderListing}
        />
    );
}
