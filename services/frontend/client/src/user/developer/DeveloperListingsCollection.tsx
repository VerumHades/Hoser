import { useCallback, useState } from "react"
import backend_constants from "../../backend_constants"
import { API, type DeveloperListing } from "../../backend"
import Query from "../../components/querying/Query"
import { CollectionViewContainer } from "../../components/view/CollectionViewContainer"
import { DeveloperListingRow } from "./List/DeveloperListingRow"
import { DeveloperListingCard } from "./List/DeveloperListingCard"
import { FlowSwitch } from "../../components/navigation/Flow"

type ListingViewMode = "cards" | "table"

interface DeveloperListingsCollectionProps {
    onSelect: (listing: DeveloperListing) => void
}

export default function DeveloperListingsCollection({ onSelect }: DeveloperListingsCollectionProps) {
    const [viewMode, setViewMode] = useState<ListingViewMode>("cards")
    const [creating, setCreating] = useState(false)

    const renderListing = useCallback(
        (listings: DeveloperListing[] | undefined) => (
            <CollectionViewContainer
                items={listings}
                viewMode={viewMode}
                onViewModeChange={setViewMode}
                renderTableRow={(listing) => (
                    <DeveloperListingRow
                        key={listing.id}
                        listing={listing}
                        onSelect={() => onSelect(listing)}
                    />
                )}
                renderCard={(listing) => (
                    <DeveloperListingCard
                        key={listing.id}
                        listing={listing}
                        onSelect={() => onSelect(listing)}
                    />
                )}
                emptyState={<p className="text-sm text-slate-500">No listings found</p>}
            />
        ),
        [viewMode, onSelect]
    )

    const queryBuilder = useCallback((queryText: string) => ({ q: queryText }), [])

    const createListing = async () => {
        setCreating(true)
        const response = await API.developer.listing.create()
        if (response.ok) {
            onSelect(response.json)
        }
        setCreating(false)
    }

    return (
        <div className="flex flex-1 flex-col h-full min-h-0 p-6">
            <div className="flex items-center justify-between mb-4">
                <FlowSwitch direction="next" onClick={createListing}>
                    New Listing
                </FlowSwitch>
            </div>

            <Query<DeveloperListing[]>
                endpoint={`${backend_constants.address}/developer/listings`}
                queryBuilder={queryBuilder}
                className="flex-1 min-h-0 overflow-hidden"
                render={renderListing}
            />
        </div>
    )
}
