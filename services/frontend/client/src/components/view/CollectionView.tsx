
import React from "react"
import Table from "../table/Table"

type CollectionViewMode = "table" | "cards"

type CollectionViewProps<ItemType> = {
    items?: ItemType[]
    viewMode: CollectionViewMode
    renderTableRow: (item: ItemType) => React.ReactNode
    renderCard: (item: ItemType) => React.ReactNode
    emptyState?: React.ReactNode
}

export function CollectionView<ItemType>({
    items,
    viewMode,
    renderTableRow,
    renderCard,
    emptyState
}: CollectionViewProps<ItemType>) {
    if (!items || items.length === 0) {
        return <>{emptyState ?? null}</>
    }

    if (viewMode === "table") {
        return (
            <Table className="flex-1 overflow-y-auto min-h-0">
                {items.map((item, index) => (
                    <React.Fragment key={index}>
                        {renderTableRow(item)}
                    </React.Fragment>
                ))}
            </Table>
        )
    }

    return (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 overflow-y-auto">
            {items.map((item, index) => (
                <React.Fragment key={index}>
                    {renderCard(item)}
                </React.Fragment>
            ))}
        </div>
    )
}