import { useCallback, useState } from "react"
import { useNavigate } from "react-router-dom"
import { API, type Listing } from "../backend"
import Query from "../components/querying/Query"
import { CollectionViewContainer, type CollectionViewMode } from "../components/view/CollectionViewContainer"
import { PublicListingRow } from "./List/PublicListingRow"
import { PublicListingCard } from "./List/PublicListingCard"
import backend_constants from "../backend_constants"

interface PublicListingSearchProps {
    onSelect?: (listing: Listing) => void
}

export default function PublicListingSearch({ onSelect }: PublicListingSearchProps) {
    const navigate = useNavigate()
    const [viewMode, setViewMode] = useState<CollectionViewMode>("cards")

    const queryBuilder = 
        (queryText: string) => {
            return { q: queryText }
        }

    const render = useCallback((
                results: Listing[],
                goToNext?: () => void,
                goToPrev?: () => void,
                currentPage?: number,
                totalPages?: number
            ) =>
                <CollectionViewContainer
                    items={results}
                    viewMode={viewMode}
                    onViewModeChange={setViewMode}
                    renderTableRow={listing => (
                        <PublicListingRow
                            key={listing.id}
                            listing={listing}
                            onSelect={onSelect ? () => onSelect(listing) : () => navigate(`/listing/${listing.id}`)}
                        />
                    )}
                    renderCard={listing => (
                        <PublicListingCard
                            key={listing.id}
                            listing={listing}
                            onSelect={onSelect ? () => onSelect(listing) : () => navigate(`/listing/${listing.id}`)}
                        />
                    )}
                    emptyState={<p className="text-sm text-slate-500">No listings found</p>}
                    currentPage={currentPage}
                    totalPages={totalPages}
                    onPageChange={page => {
                        if (page < (currentPage ?? 0)) goToPrev?.()
                        else if (page > (currentPage ?? 0)) goToNext?.()
                    }}
                />
    , [onSelect, navigate])

    return (
        <Query<Listing>
            endpoint={`${backend_constants.address}/search/listings`}
            queryBuilder={queryBuilder}
            className="w-full max-w-4xl min-h-0 overflow-hidden mt-5 px-6"
            render={render}
        />

    )
}
