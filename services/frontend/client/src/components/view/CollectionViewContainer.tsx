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


/** 
 * Page-based pagination with numbered buttons, ellipsis, and prev/next arrows.
 */
interface PaginationProps {
    currentPage: number;
    totalPages: number;
    isLoading: boolean;
    onPageChange: (pageIndex: number) => void;
}

/**
 * Pagination control using zero-based page indices internally
 */
export function Pagination({
    currentPage,
    totalPages,
    isLoading,
    onPageChange
}: PaginationProps) {
    const displayPageNumber = currentPage + 1;
    const isFirstPage = currentPage === 0;
    const isLastPage = currentPage === totalPages - 1;

    return (
        <div className="flex items-center justify-center gap-2 mt-4">
            <button
                type="button"
                onClick={() => onPageChange(currentPage - 1)}
                disabled={isFirstPage || isLoading}
                className="px-3 py-1 rounded-md bg-gray-200 dark:bg-gray-700
                           hover:bg-gray-300 dark:hover:bg-gray-600
                           disabled:opacity-40 disabled:cursor-not-allowed"
            >
                &lt;
            </button>

            <span className="min-w-[2rem] text-center font-medium text-gray-900 dark:text-gray-100">
                {displayPageNumber}
            </span>

            <button
                type="button"
                onClick={() => onPageChange(currentPage + 1)}
                disabled={isLastPage || isLoading}
                className="px-3 py-1 rounded-md bg-gray-200 dark:bg-gray-700
                           hover:bg-gray-300 dark:hover:bg-gray-600
                           disabled:opacity-40 disabled:cursor-not-allowed"
            >
                &gt;
            </button>
        </div>
    );
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
                        isLoading={isLoading ?? false}
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
