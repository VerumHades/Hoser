import { useCallback, useEffect, useState } from "react";
import { CursorPaginator, type CursorPaginatedResult } from "../../backend/paginator";
import { CollectionViewContainer, type CollectionViewMode } from "./CollectionViewContainer";

interface CursorPaginatedCollectionProps<TItem> {
    fetchPage: (cursor?: string) => Promise<CursorPaginatedResult<TItem>>;
    renderItemRow: (item: TItem) => React.ReactNode;
    renderItemCard: (item: TItem) => React.ReactNode;
    emptyState?: React.ReactNode;
}

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

    useEffect(() => {
        const paginatorInstance = new CursorPaginator<TItem>(fetchPage);
        setPaginator(paginatorInstance);

        setIsLoading(true);
        paginatorInstance
            .getCurrentPage()
            .then(items => {
                setCurrentItems(items);
                setCurrentPage(paginatorInstance.getCurrentPageIndex());
            })
            .finally(() => setIsLoading(false));
    }, [fetchPage]);

    const goNextPage = useCallback(async () => {
        if (!paginator) return;

        setIsLoading(true);
        const items = await paginator.nextPage();
        if (items) {
            setCurrentItems(items);
            setCurrentPage(paginator.getCurrentPageIndex());
        }
        setIsLoading(false);
    }, [paginator]);

    const goPreviousPage = useCallback(async () => {
        if (!paginator) return;

        setIsLoading(true);
        const items = await paginator.prevPage();
        if (items) {
            setCurrentItems(items);
            setCurrentPage(paginator.getCurrentPageIndex());
        }
        setIsLoading(false);
    }, [paginator]);

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
            onPageChange={pageIndex => {
                if (pageIndex < currentPage) {
                    goPreviousPage();
                } else if (pageIndex > currentPage) {
                    goNextPage();
                }
            }}
            isLoading={isLoading}
        />
    );
}