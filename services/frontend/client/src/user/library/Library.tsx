import { useState, useCallback } from "react"
import { useNavigate } from "react-router-dom"
import { CollectionViewContainer } from "../../components/view/CollectionViewContainer"
import { UserLibraryRow } from "./List/UserLibraryRow"
import { UserLibraryCard } from "./List/UserLibraryCard"
import { CursorPaginatedCollection } from "../../components/view/CursorPaginatedCollection"
import type { Listing } from "../../backend/repositories/listing"
import { UserAPI } from "../../backend/repositories/user"

type LibraryViewMode = "cards" | "table"

export default function UserLibrary() {
    const navigate = useNavigate()

    const queryBuilder = useCallback((queryText: string) => ({ q: queryText }), [])

    const handleCreateInstance = useCallback((listingId: string) => {
        navigate(`/dashboard/create-instance/${listingId}`)
    }, [navigate])

    return (
        <div className="flex-1 min-h-0 flex flex-col items-center">
            <div className="w-full max-w-4xl min-h-0 overflow-hidden mt-5 px-6">
                <CursorPaginatedCollection<Listing>
                    fetchPage={cursor => UserAPI.library.listBatch(cursor)}
                    renderItemRow={entry => (
                        <UserLibraryRow
                            key={entry.id}
                            listing={entry}
                            onCreateInstance={handleCreateInstance}
                        />
                    )}
                    renderItemCard={entry => (
                        <UserLibraryCard
                            key={entry.id}
                            listing={entry}
                            onCreateInstance={handleCreateInstance}
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
    )
}
