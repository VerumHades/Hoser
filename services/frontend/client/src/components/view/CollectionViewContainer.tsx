import React, { useEffect, useRef, useState, type ReactNode } from "react";
import { 
    LayoutGrid, 
    Table as TableIcon, 
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
                {visiblePages.map((pageIndex) => (
                    <PageButton
                        key={pageIndex}
                        index={pageIndex}
                        isActive={pageIndex === currentPage}
                        isLoading={isLoading}
                        onClick={onPageChange}
                    />
                ))}
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
 * Individual page number button.
 */
function PageButton({ 
    index, 
    isActive, 
    isLoading, 
    onClick 
}: { 
    index: number; 
    isActive: boolean; 
    isLoading: boolean; 
    onClick: (idx: number) => void 
}) {
    const activeStyles = "bg-slate-900 text-white shadow-lg shadow-slate-200 dark:bg-white dark:text-slate-900 dark:shadow-none";
    const inactiveStyles = "text-slate-500 hover:bg-slate-100 dark:text-slate-400 dark:hover:bg-slate-800";

    return (
        <button
            disabled={isLoading}
            onClick={() => onClick(index)}
            className={`flex h-9 min-w-[36px] items-center justify-center rounded-xl text-sm font-semibold transition-all ${isActive ? activeStyles : inactiveStyles}`}
        >
            {index + 1}
        </button>
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
    const buttonBase = "relative flex h-8 w-10 items-center justify-center rounded-lg transition-all";
    const activeStyles = "bg-white text-slate-900 shadow-sm dark:bg-slate-700 dark:text-white";
    const inactiveStyles = "text-slate-500 hover:text-slate-700 dark:text-slate-400";

    return (
        <div className="flex rounded-xl bg-slate-100 p-1 dark:bg-slate-800/50">
            <button
                onClick={() => onViewModeChange("table")}
                className={`${buttonBase} ${viewMode === "table" ? activeStyles : inactiveStyles}`}
            >
                <TableIcon size={16} />
            </button>
            <button
                onClick={() => onViewModeChange("cards")}
                className={`${buttonBase} ${viewMode === "cards" ? activeStyles : inactiveStyles}`}
            >
                <LayoutGrid size={16} />
            </button>
        </div>
    );
}

/**
 * Main container with a clean, focused UI.
 */
export function CollectionViewContainer<ItemType>(props: CollectionViewContainerProps<ItemType>) {
    const scrollContainerReference = useRef<HTMLDivElement>(null);
    const [isContentScrollable, setIsContentScrollable] = useState(false);

    useEffect(() => {
        scrollToTop(scrollContainerReference);
        checkScrollability(scrollContainerReference, setIsContentScrollable);
    }, [props.currentPage, props.items, props.viewMode]);

    useEffect(() => {
        const observer = createResizeObserver(scrollContainerReference, setIsContentScrollable);
        return () => observer.disconnect();
    }, []);

    const showPagination = (props.totalPages ?? 0) > 1 && !!props.onPageChange;

    return (
        <div className="flex h-full w-full flex-col bg-white dark:bg-slate-950">
            <header className="flex h-20 items-center justify-between border-b border-slate-100 px-8 dark:border-slate-900">
                <div className="flex items-center gap-8">
                    <h2 className="text-lg font-bold tracking-tight text-slate-900 dark:text-white">
                        {props.header}
                    </h2>
                    {showPagination && <div className="hidden md:block"><PaginationLayout {...props} /></div>}
                </div>
                <ViewModeToggle viewMode={props.viewMode} onViewModeChange={props.onViewModeChange} />
            </header>

            <main ref={scrollContainerReference} className="flex-1 overflow-y-auto scroll-smooth">
                <div className="mx-auto max-w-7xl p-8">
                    <CollectionView {...props} />
                    {showPagination && isContentScrollable && (
                        <div className="mt-12 flex justify-center py-12">
                            <PaginationLayout {...props} />
                        </div>
                    )}
                </div>
            </main>
        </div>
    );
}

/**
 * Helper to check if the element has vertical overflow.
 */
function checkScrollability(
    ref: React.RefObject<HTMLDivElement>, 
    setScrollable: (val: boolean) => void
) {
    if (ref.current) {
        const { scrollHeight, clientHeight } = ref.current;
        setScrollable(scrollHeight > clientHeight);
    }
}

/**
 * Smooth scrolls container to the top.
 */
function scrollToTop(ref: React.RefObject<HTMLDivElement>) {
    ref.current?.scrollTo({ top: 0, behavior: "smooth" });
}

/**
 * Sets up a ResizeObserver to detect content height changes.
 */
function createResizeObserver(
    ref: React.RefObject<HTMLDivElement>, 
    setScrollable: (val: boolean) => void
) {
    const observer = new ResizeObserver(() => checkScrollability(ref, setScrollable));
    if (ref.current) observer.observe(ref.current);
    return observer;
}

/**
 * Structural wrapper for the pagination component.
 */
function PaginationLayout<ItemType>(props: CollectionViewContainerProps<ItemType>) {
    return (
        <Pagination
            currentPage={props.currentPage ?? 0}
            totalPages={props.totalPages ?? 1}
            onPageChange={props.onPageChange!}
            isLoading={props.isLoading ?? false}
        />
    );
}