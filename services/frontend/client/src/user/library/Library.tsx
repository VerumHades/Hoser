import { useState } from "react"
import { useNavigate } from "react-router-dom"
import backend_constants from "../../backend_constants"

import Query from "../../components/querying/Query"
import type { Listing } from "../../backend"

import { CollectionView } from "../../components/view/CollectionView"
import { UserLibraryRow } from "./UserLibraryRow"
import { UserLibraryCard } from "./UserLibraryCard"

type LibraryViewMode = "cards" | "table"

export default function UserLibrary() {
    const navigate = useNavigate()
    const [viewMode, setViewMode] =
        useState<LibraryViewMode>("cards")

    const queryBuilder = (queryText: string) => ({
        q: queryText
    })

    const handleCreateInstance = (listingId: string) => {
        navigate(`/dashboard/create-instance/${listingId}`)
    }

    return (
        <div className="flex flex-col gap-4 h-full">
            <div className="flex justify-end gap-2">
                <button onClick={() => setViewMode("cards")}>
                    Cards
                </button>
                <button onClick={() => setViewMode("table")}>
                    Table
                </button>
            </div>

            <Query<Listing[]>
                className="flex-1 min-h-0"
                endpoint={`${backend_constants.address}/user/library`}
                queryBuilder={queryBuilder}
                render={(entries) => (
                    <CollectionView
                        items={entries}
                        viewMode={viewMode}
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
                )}
            />
        </div>
    )
}
