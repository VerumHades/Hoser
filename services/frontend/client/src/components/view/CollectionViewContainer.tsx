import React from "react"
import { LayoutGrid, Table as TableIcon, ChevronsLeft, ChevronsRight, ChevronLeft, ChevronRight } from "lucide-react"
import { CollectionView } from "./CollectionView"

export type CollectionViewMode = "table" | "cards"

type CollectionViewContainerProps<ItemType> = {
    items?: ItemType[]
    viewMode: CollectionViewMode
    onViewModeChange: (nextViewMode: CollectionViewMode) => void
    renderTableRow: (item: ItemType) => React.ReactNode
    renderCard: (item: ItemType) => React.ReactNode
    emptyState?: React.ReactNode
    currentPage?: number
    totalPages?: number
    onPageChange?: (page: number) => void
    isLoading?: boolean
}

interface PaginationProps {
    currentPage: number;
    totalPages: number;
    isLoading: boolean;
    onPageChange: (pageIndex: number) => void;
}

/**
 * Calculates the start of the window to keep the current page centered.
 */
function calculateWindowStart(currentPage: number, totalPages: number, windowSize: number): number {
    const halfWindow = Math.floor(windowSize / 2);
    const idealStart = currentPage - halfWindow;
    const maxStart = totalPages - windowSize;
    return Math.max(0, Math.min(idealStart, maxStart));
}

/**
 * Generates an array of sequential page numbers centered around the current page.
 */
function generateCenterWindow(currentPage: number, totalPages: number): number[] {
    const windowSize = Math.min(5, totalPages);
    const startPage = calculateWindowStart(currentPage, totalPages, windowSize);
    
    return Array.from({ length: windowSize }, (_, index) => startPage + index);
}

/**
 * Determines the styling for a page button based on its active state.
 */
function getPageButtonClassName(isActive: boolean): string {
    const baseStyles = "px-3 py-1 rounded-md transition-colors min-w-[40px] flex justify-center items-center";
    const activeStyles = "bg-blue-600 text-white font-bold";
    const inactiveStyles = "bg-gray-200 dark:bg-gray-700 hover:bg-gray-300 dark:hover:bg-gray-600 text-gray-900 dark:text-gray-100";
    
    return `${baseStyles} ${isActive ? activeStyles : inactiveStyles}`;
}

/**
 * Renders an individual page navigation button.
 */
function PageButton({ 
    pageIndex, 
    isActive, 
    isDisabled, 
    onClick 
}: { 
    pageIndex: number, 
    isActive: boolean, 
    isDisabled: boolean, 
    onClick: (index: number) => void 
}) {
    return (
        <button
            type="button"
            onClick={() => onClick(pageIndex)}
            disabled={isDisabled}
            className={getPageButtonClassName(isActive)}
        >
            {pageIndex + 1}
        </button>
    );
}

/**
 * Renders a navigation icon button for jumping or stepping through pages.
 */
function IconButton({ 
    icon: Icon, 
    isDisabled, 
    onClick 
}: { 
    icon: any, 
    isDisabled: boolean, 
    onClick: () => void 
}) {
    return (
        <button
            type="button"
            onClick={onClick}
            disabled={isDisabled}
            className="p-2 rounded-md bg-gray-200 dark:bg-gray-700 hover:bg-gray-300 dark:hover:bg-gray-600 disabled:opacity-30 transition-opacity"
        >
            <Icon size={16} />
        </button>
    );
}

/**
 * Renders the jump controls (±10) and step controls (±1).
 */
function NavigationControls({ 
    currentPage, 
    totalPages, 
    isLoading, 
    onPageChange 
}: PaginationProps) {
    const jumpAmount = 10;
    
    return (
        <>
            <IconButton 
                icon={ChevronsLeft} 
                isDisabled={currentPage < jumpAmount || isLoading} 
                onClick={() => onPageChange(currentPage - jumpAmount)} 
            />
            <IconButton 
                icon={ChevronLeft} 
                isDisabled={currentPage === 0 || isLoading} 
                onClick={() => onPageChange(currentPage - 1)} 
            />
        </>
    );
}

/**
 * Pagination control with a centered sliding window and distal jump buttons.
 */
export function Pagination({
    currentPage,
    totalPages,
    isLoading,
    onPageChange
}: PaginationProps) {
    const visiblePages = generateCenterWindow(currentPage, totalPages);
    const jumpAmount = 10;

    return (
        <div className="flex items-center justify-center gap-2 mt-4">
            <NavigationControls 
                currentPage={currentPage} 
                totalPages={totalPages} 
                isLoading={isLoading} 
                onPageChange={onPageChange} 
            />

            {visiblePages.map((pageIndex) => (
                <PageButton
                    key={pageIndex}
                    pageIndex={pageIndex}
                    isActive={pageIndex === currentPage}
                    isDisabled={isLoading}
                    onClick={onPageChange}
                />
            ))}

            <IconButton 
                icon={ChevronRight} 
                isDisabled={currentPage >= totalPages - 1 || isLoading} 
                onClick={() => onPageChange(currentPage + 1)} 
            />
            <IconButton 
                icon={ChevronsRight} 
                isDisabled={currentPage > totalPages - 1 - jumpAmount || isLoading} 
                onClick={() => onPageChange(currentPage + jumpAmount)} 
            />
        </div>
    );
}

/**
 * Renders the view mode selection buttons.
 */
function ViewModeToggle({ 
    viewMode, 
    onViewModeChange 
}: { 
    viewMode: CollectionViewMode, 
    onViewModeChange: (mode: CollectionViewMode) => void 
}) {
    const activeClass = "bg-gray-200 dark:bg-gray-600";
    const inactiveClass = "text-gray-500 hover:bg-gray-100 dark:hover:bg-gray-700";

    return (
        <div className="flex justify-end gap-2">
            <button
                type="button"
                onClick={() => onViewModeChange("table")}
                className={`p-2 rounded-md transition-colors ${viewMode === "table" ? activeClass : inactiveClass}`}
            >
                <TableIcon size={18} />
            </button>
            <button
                type="button"
                onClick={() => onViewModeChange("cards")}
                className={`p-2 rounded-md transition-colors ${viewMode === "cards" ? activeClass : inactiveClass}`}
            >
                <LayoutGrid size={18} />
            </button>
        </div>
    );
}

/**
 * Collection view container with centralized pagination and view mode toggling.
 */
export function CollectionViewContainer<ItemType>({
    items,
    viewMode,
    onViewModeChange,
    renderTableRow,
    renderCard,
    emptyState,
    currentPage = 0,
    totalPages = 1,
    onPageChange,
    isLoading,
}: CollectionViewContainerProps<ItemType>) {
    const hasMultiplePages = totalPages > 1 && onPageChange;

    return (
        <div className="flex flex-col gap-4 w-full h-full min-h-0">
            <div className="flex justify-between items-center">
                <div className="flex-1">
                    {hasMultiplePages && (
                        <Pagination
                            currentPage={currentPage}
                            totalPages={totalPages}
                            onPageChange={onPageChange}
                            isLoading={isLoading ?? false}
                        />
                    )}
                </div>
                <ViewModeToggle 
                    viewMode={viewMode} 
                    onViewModeChange={onViewModeChange} 
                />
            </div>

            <CollectionView
                items={items}
                viewMode={viewMode}
                renderTableRow={renderTableRow}
                renderCard={renderCard}
                emptyState={emptyState}
            />
        </div>
    );
}