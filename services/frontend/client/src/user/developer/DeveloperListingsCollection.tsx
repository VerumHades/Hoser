import { useCallback, useState } from "react"
import backend_constants from "../../backend_constants"
import { API, type DeveloperListing } from "../../backend"
import Query from "../../components/querying/Query"
import { CollectionView } from "../../components/view/CollectionView"
import { DeveloperListingRow } from "./DeveloperListingRow"
import { DeveloperListingCard } from "./DeveloperListingCard"
import { FlowSwitch } from "../../components/navigation/Flow"

type ListingViewMode = "cards" | "table"

interface DeveloperListingsCollectionProps {
    onSelect: (listing: DeveloperListing) => void
}

export default function DeveloperListingsCollection({ onSelect }: DeveloperListingsCollectionProps) {
    const [viewMode, setViewMode] = useState<ListingViewMode>("cards")

    const renderListing = useCallback(
        (listings: DeveloperListing[] | undefined) => (
            <CollectionView
                items={listings}
                viewMode={viewMode}
                renderTableRow={(listing) => <DeveloperListingRow key={listing.id} listing={listing} onSelect={onSelect} />}
                renderCard={(listing) => <DeveloperListingCard key={listing.id} listing={listing} onSelect={onSelect} />}
                emptyState={<p className="text-sm text-slate-500">No listings found</p>}
            />
        ),
        [viewMode, onSelect]
    )

    const queryBuilder = useCallback((queryText: string) => ({ q: queryText }), [])

    const [creating, setCreating] = useState(false)
    const createListing = async () => {
        setCreating(true)
        const response = await API.developer.listing.create()
        if (response.ok) {
            onSelect(response.json)
        }
        setCreating(false)
    }

    return (
        <div className="flex flex-col h-full">
            <div className="flex items-center justify-between mb-4">
                <FlowSwitch direction="next" onClick={createListing}>
                    New Listing
                </FlowSwitch>
                <div className="flex gap-2">
                    <button onClick={() => viewMode !== "cards" && viewMode} className={viewMode === "cards" ? "font-semibold" : ""}>Cards</button>
                    <button onClick={() => viewMode !== "table" && viewMode} className={viewMode === "table" ? "font-semibold" : ""}>Table</button>
                </div>
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
