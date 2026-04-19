import React, { useEffect, useRef, type ReactNode } from "react";
import { 
    LayoutGrid, 
    Table as TableIcon, 
    ChevronsLeft, 
    ChevronsRight, 
    ChevronLeft, 
    ChevronRight 
} from "lucide-react";
import { CollectionView } from "./CollectionView";

export type CollectionViewMode = "table" | "cards";

type CollectionViewContainerProps<ItemType> = {
    items?: ItemType[];
    viewMode: CollectionViewMode;
    onViewModeChange: (nextViewMode: CollectionViewMode) => void;
    renderTableRow: (item: ItemType) => React.ReactNode;
    renderCard: (item: ItemType) => React.ReactNode;
    emptyState?: React.ReactNode;
    currentPage?: number;
    totalPages?: number;
    onPageChange?: (page: number) => void;
    isLoading?: boolean;
    header?: ReactNode;
};

/**
 * Renders a stylized icon button for pagination navigation.
 */
function IconButton({ 
    icon: Icon, 
    isDisabled, 
    onClick 
}: { 
    icon: any; 
    isDisabled: boolean; 
    onClick: () => void 
}) {
    return (
        <button
            type="button"
            onClick={onClick}
            disabled={isDisabled}
            className="group flex h-9 w-9 items-center justify-center rounded-xl border border-slate-200 bg-white text-slate-600 transition-all hover:border-slate-300 hover:bg-slate-50 disabled:opacity-20 dark:border-slate-800 dark:bg-slate-900 dark:text-slate-400 dark:hover:border-slate-700 dark:hover:bg-slate-800"
        >
            <Icon size={16} className="transition-transform group-active:scale-90" />
        </button>
    );
}

/**
 * Generates an array of page numbers to display in the pagination bar.
 */
function generateCenterWindow(currentPage: number, totalPages: number): number[] {
    const windowSize = Math.min(5, totalPages);
    const halfWindow = Math.floor(windowSize / 2);
    const startPage = Math.max(0, Math.min(currentPage - halfWindow, totalPages - windowSize));
    
    return Array.from({ length: windowSize }, (_, index) => startPage + index);
}

/**
 * Full pagination component with modern squircle styling.
 */
export function Pagination({
    currentPage,
    totalPages,
    isLoading,
    onPageChange
}: {
    currentPage: number;
    totalPages: number;
    isLoading: boolean;
    onPageChange: (pageIndex: number) => void;
}) {
    const visiblePages = generateCenterWindow(currentPage, totalPages);

    return (
        <nav className="flex items-center gap-2" aria-label="Pagination">
            <IconButton 
                icon={ChevronLeft} 
                isDisabled={currentPage === 0 || isLoading} 
                onClick={() => onPageChange(currentPage - 1)} 
            />

            <div className="flex items-center gap-1.5 px-2">
                {visiblePages.map((pageIndex) => {
                    const isActive = pageIndex === currentPage;
                    return (
                        <button
                            key={pageIndex}
                            disabled={isLoading}
                            onClick={() => onPageChange(pageIndex)}
                            className={`flex h-9 min-w-[36px] items-center justify-center rounded-xl text-sm font-semibold transition-all
                                ${isActive 
                                    ? "bg-slate-900 text-white shadow-lg shadow-slate-200 dark:bg-white dark:text-slate-900 dark:shadow-none" 
                                    : "text-slate-500 hover:bg-slate-100 dark:text-slate-400 dark:hover:bg-slate-800"
                                }`}
                        >
                            {pageIndex + 1}
                        </button>
                    );
                })}
            </div>

            <IconButton 
                icon={ChevronRight} 
                isDisabled={currentPage >= totalPages - 1 || isLoading} 
                onClick={() => onPageChange(currentPage + 1)} 
            />
        </nav>
    );
}

/**
 * Segmented control for view switching.
 */
function ViewModeToggle({ 
    viewMode, 
    onViewModeChange 
}: { 
    viewMode: CollectionViewMode; 
    onViewModeChange: (mode: CollectionViewMode) => void 
}) {
    const isActive = (mode: CollectionViewMode) => viewMode === mode;
    const buttonBase = "relative flex h-8 w-10 items-center justify-center rounded-lg transition-all";

    return (
        <div className="flex rounded-xl bg-slate-100 p-1 dark:bg-slate-800/50">
            <button
                onClick={() => onViewModeChange("table")}
                className={`${buttonBase} ${isActive("table") ? "bg-white text-slate-900 shadow-sm dark:bg-slate-700 dark:text-white" : "text-slate-500 hover:text-slate-700 dark:text-slate-400"}`}
            >
                <TableIcon size={16} />
            </button>
            <button
                onClick={() => onViewModeChange("cards")}
                className={`${buttonBase} ${isActive("cards") ? "bg-white text-slate-900 shadow-sm dark:bg-slate-700 dark:text-white" : "text-slate-500 hover:text-slate-700 dark:text-slate-400"}`}
            >
                <LayoutGrid size={16} />
            </button>
        </div>
    );
}

/**
 * Main container with a clean, focused UI.
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
    header
}: CollectionViewContainerProps<ItemType>) {
    const scrollContainerRef = useRef<HTMLDivElement>(null);
    const hasMultiplePages = totalPages > 1 && !!onPageChange;

    useEffect(() => {
        scrollContainerRef.current?.scrollTo({ top: 0, behavior: "smooth" });
    }, [currentPage]);

    return (
        <div className="flex h-full w-full flex-col bg-white dark:bg-slate-950">
            <header className="flex h-20 items-center justify-between border-b border-slate-100 px-8 dark:border-slate-900">
                <div className="flex items-center gap-8">
                    <h2 className="text-lg font-bold tracking-tight text-slate-900 dark:text-white">
                        {header}
                    </h2>
                    
                    {hasMultiplePages && (
                        <div className="hidden md:block">
                            <Pagination
                                currentPage={currentPage}
                                totalPages={totalPages}
                                onPageChange={onPageChange}
                                isLoading={isLoading ?? false}
                            />
                        </div>
                    )}
                </div>

                <ViewModeToggle 
                    viewMode={viewMode} 
                    onViewModeChange={onViewModeChange} 
                />
            </header>

            <main 
                ref={scrollContainerRef}
                className="flex-1 overflow-y-auto scroll-smooth"
            >
                <div className="mx-auto max-w-7xl p-8">
                    <CollectionView
                        items={items}
                        viewMode={viewMode}
                        renderTableRow={renderTableRow}
                        renderCard={renderCard}
                        emptyState={emptyState}
                    />
                    
                    {hasMultiplePages && (
                        <div className="mt-12 flex justify-center py-12">
                            <Pagination
                                currentPage={currentPage}
                                totalPages={totalPages}
                                onPageChange={onPageChange}
                                isLoading={isLoading ?? false}
                            />
                        </div>
                    )}
                </div>
            </main>
        </div>
    );
}