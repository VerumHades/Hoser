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
    /** Current page number (1-based) */
    currentPage?: number
    /** Total number of pages available */
    totalPages?: number
    /** Callback when the user changes the page */
    onPageChange?: (page: number) => void
    /** Whether a page fetch is in progress */
    isLoading?: boolean
}

type PaginationProps = {
    currentPage: number
    totalPages: number
    onPageChange: (page: number) => void
    isLoading?: boolean
}


/**
 * Page-based pagination with numbered buttons, ellipsis, and prev/next arrows.
 */
/**
 * Page-based pagination with numbered buttons, ellipsis, and prev/next arrows.
 */
function Pagination({ currentPage, totalPages, onPageChange, isLoading }: PaginationProps) {
    const createPageArray = () => {
        const pages: (number | "...")[] = []
        if (totalPages <= 7) {
            for (let i = 1; i <= totalPages; i++) pages.push(i)
        } else {
            pages.push(1)
            if (currentPage > 4) pages.push("...")
            const start = Math.max(2, currentPage - 1)
            const end = Math.min(totalPages - 1, currentPage + 1)
            for (let i = start; i <= end; i++) pages.push(i)
            if (currentPage < totalPages - 3) pages.push("...")
            pages.push(totalPages)
        }
        return pages
    }

    const pages = createPageArray()

    return (
        <div className="flex justify-center items-center gap-1 mt-2">
            <button
                onClick={() => currentPage > 1 && onPageChange(currentPage - 1)}
                disabled={currentPage === 1 || isLoading}
                className="px-3 py-1 rounded bg-gray-200 dark:bg-gray-700 hover:bg-gray-300 dark:hover:bg-gray-600 disabled:opacity-50"
            >
                &lt;
            </button>

            {pages.map((page, index) =>
                page === "..." ? (
                    <span key={index} className="px-2 py-1">
                        …
                    </span>
                ) : (
                    <button
                        key={page}
                        onClick={() => onPageChange(page as number)}
                        className={`px-3 py-1 rounded transition ${
                            page === currentPage
                                ? "bg-indigo-600 text-white"
                                : "bg-gray-200 dark:bg-gray-700 hover:bg-gray-300 dark:hover:bg-gray-600"
                        }`}
                        disabled={isLoading}
                    >
                        {page}
                    </button>
                )
            )}

            <button
                onClick={() => currentPage < totalPages && onPageChange(currentPage + 1)}
                disabled={currentPage === totalPages || isLoading}
                className="px-3 py-1 rounded bg-gray-200 dark:bg-gray-700 hover:bg-gray-300 dark:hover:bg-gray-600 disabled:opacity-50"
            >
                &gt;
            </button>
        </div>
    )
}



/**
 * Collection view container with explicit page-based pagination
 * and toggleable table/card views.
 */
export function CollectionViewContainer<ItemType>({
    items,
    viewMode,
    onViewModeChange,
    renderTableRow,
    renderCard,
    emptyState,
    currentPage = 1,
    totalPages = 1,
    onPageChange,
    isLoading,
}: CollectionViewContainerProps<ItemType>) {
    return (
        <div className="flex flex-col gap-3 w-full h-full min-h-0">
            <div className="flex justify-between">
                {/* Pagination Controls */}
                {totalPages > 1 && onPageChange && (
                    <Pagination
                        currentPage={currentPage}
                        totalPages={totalPages}
                        onPageChange={onPageChange}
                        isLoading={isLoading}
                    />
                )}
                {/* View Mode Toggle */}
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
            </div>
            

            {/* Collection View */}
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
