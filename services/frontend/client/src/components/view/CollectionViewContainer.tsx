import React from "react"
import { LayoutGrid, Table as TableIcon } from "lucide-react"
import { CollectionView } from "./CollectionView"

export type CollectionViewMode = "table" | "cards"

type CollectionViewContainerProps<ItemType> = {
    items?: ItemType[]
    viewMode: CollectionViewMode
    onViewModeChange: (nextViewMode: CollectionViewMode) => void
    renderTableRow: (item: ItemType) => React.ReactNode
    renderCard: (item: ItemType) => React.ReactNode
    emptyState?: React.ReactNode
}

/**
 * Presentational container that exposes explicit controls
 * for switching collection view modes.
 * State is owned by the parent.
 */
export function CollectionViewContainer<ItemType>({
    items,
    viewMode,
    onViewModeChange,
    renderTableRow,
    renderCard,
    emptyState,
}: CollectionViewContainerProps<ItemType>) {
    return (
        <div className="flex flex-col gap-3 w-full h-full min-h-0">
            <div className="flex justify-end gap-2">
                <button
                    type="button"
                    aria-label="Table view"
                    onClick={() => onViewModeChange("table")}
                    className={`p-2 rounded-md transition ${
                        viewMode === "table"
                            ? "bg-gray-200 text-gray-900"
                            : "text-gray-500 hover:bg-gray-100"
                    }`}
                >
                    <TableIcon size={18} />
                </button>

                <button
                    type="button"
                    aria-label="Card view"
                    onClick={() => onViewModeChange("cards")}
                    className={`p-2 rounded-md transition ${
                        viewMode === "cards"
                            ? "bg-gray-200 text-gray-900"
                            : "text-gray-500 hover:bg-gray-100"
                    }`}
                >
                    <LayoutGrid size={18} />
                </button>
            </div>

            <CollectionView
                items={items}
                viewMode={viewMode}
                renderTableRow={renderTableRow}
                renderCard={renderCard}
                emptyState={emptyState}
            />
        </div>
    )
}
