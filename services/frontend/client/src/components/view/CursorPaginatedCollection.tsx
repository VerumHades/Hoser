import { useCallback, useEffect, useState } from "react";
import { CursorPaginator, type CursorPaginatedResult } from "../../backend/paginator";
import { CollectionViewContainer, type CollectionViewMode } from "./CollectionViewContainer";

interface CursorPaginatedCollectionProps<TItem> {
    fetchPage: (cursor?: string) => Promise<CursorPaginatedResult<TItem> | null>;
    renderItemRow: (item: TItem) => React.ReactNode;
    renderItemCard: (item: TItem) => React.ReactNode;
    emptyState?: React.ReactNode;
}

/**
 * Manages cursor-based navigation with support for direct page jumping.
 */
export function CursorPaginatedCollection<TItem>({
    fetchPage,
    renderItemRow,
    renderItemCard,
    emptyState
}: CursorPaginatedCollectionProps<TItem>) {
    const [viewMode, setViewMode] = useState<CollectionViewMode>("cards");
    const [paginator, setPaginator] = useState<CursorPaginator<TItem>>();
    const [currentItems, setCurrentItems] = useState<TItem[]>([]);
    const [currentPage, setCurrentPage] = useState(0);
    const [isLoading, setIsLoading] = useState(false);

    /**
     * Synchronizes local state with the current paginator position.
     */
    const updateCollectionState = useCallback((items: TItem[] | null, paginatorInstance: CursorPaginator<TItem>) => {
        if (items) {
            setCurrentItems(items);
            setCurrentPage(paginatorInstance.getCurrentPageIndex());
        }
    }, []);

    /**
     * Handles the logic of moving the paginator to a specific index.
     */
    const navigateToPage = useCallback(async (targetPageIndex: number) => {
        if (!paginator || isLoading) return;

        setIsLoading(true);
        try {
            const items = await paginator.goToPage(targetPageIndex);
            updateCollectionState(items, paginator);
        } finally {
            setIsLoading(false);
        }
    }, [paginator, isLoading, updateCollectionState]);

    /**
     * Initializes the paginator instance and loads the first page.
     */
    useEffect(() => {
        const paginatorInstance = new CursorPaginator<TItem>(fetchPage);
        setPaginator(paginatorInstance);

        setIsLoading(true);
        paginatorInstance
            .getCurrentPage()
            .then((items) => updateCollectionState(items, paginatorInstance))
            .finally(() => setIsLoading(false));
    }, [fetchPage, updateCollectionState]);

    return (
        <CollectionViewContainer
            items={currentItems}
            viewMode={viewMode}
            onViewModeChange={setViewMode}
            renderTableRow={renderItemRow}
            renderCard={renderItemCard}
            emptyState={emptyState}
            currentPage={currentPage}
            totalPages={paginator?.pageCount()}
            onPageChange={navigateToPage}
            isLoading={isLoading}
        />
    );
}