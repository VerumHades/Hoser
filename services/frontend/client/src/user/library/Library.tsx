import { useState, useCallback } from "react"
import { useNavigate } from "react-router-dom"
import backend_constants from "../../backend_constants"

import Query from "../../components/querying/Query"
import type { Listing } from "../../backend"

import { CollectionViewContainer } from "../../components/view/CollectionViewContainer"
import { UserLibraryRow } from "./List/UserLibraryRow"
import { UserLibraryCard } from "./List/UserLibraryCard"

type LibraryViewMode = "cards" | "table"

export default function UserLibrary() {
    const navigate = useNavigate()
    const [viewMode, setViewMode] = useState<LibraryViewMode>("cards")

    const queryBuilder = useCallback((queryText: string) => ({ q: queryText }), [])

    const handleCreateInstance = useCallback((listingId: string) => {
        navigate(`/dashboard/create-instance/${listingId}`)
    }, [navigate])

    const renderLibrary = useCallback(
        (entries: Listing[] | undefined) => (
            <CollectionViewContainer
                items={entries}
                viewMode={viewMode}
                onViewModeChange={setViewMode}
                renderTableRow={(entry: Listing) => (
                    <UserLibraryRow
                        key={entry.id}
                        listing={entry}
                        onCreateInstance={handleCreateInstance}
                    />
                )}
                renderCard={(entry: Listing) => (
                    <UserLibraryCard
                        key={entry.id}
                        listing={entry}
                        onCreateInstance={handleCreateInstance}
                    />
                )}
                emptyState={
                    <p className="text-sm text-slate-500">
                        Your library is empty
                    </p>
                }
            />
        ),
        [viewMode, handleCreateInstance]
    )

    return (
        <div className="flex flex-1 flex-col p-6 gap-4 h-full">
            <Query<Listing[]>
                className="flex-1 min-h-0"
                endpoint={`${backend_constants.address}/user/library`}
                queryBuilder={queryBuilder}
                render={renderLibrary}
            />
        </div>
    )
}
