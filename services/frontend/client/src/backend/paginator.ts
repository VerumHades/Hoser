export interface CursorPaginatedResult<T> {
    items: T[];
    cursor?: string;
}

/**
 * Generic stateful paginator for cursor-based APIs.
 * Stores history of pages and supports dynamic page navigation.
 */
export class CursorPaginator<T> {
    private fetchFunction: (cursor?: string) => Promise<CursorPaginatedResult<T>>;
    
    private pageCursors: (string | undefined)[] = [undefined]; // history of cursors for each page
    private currentPageIndex = 0;

    constructor(fetchFunction: (cursor?: string) => Promise<CursorPaginatedResult<T>>) {
        this.fetchFunction = fetchFunction;
    }

    /**
     * Fetches and caches the page at the current index.
     */
    async getCurrentPage(): Promise<T[]> {
        const cursor = this.pageCursors[this.currentPageIndex];
        const { items, cursor: nextCursor } = await this.fetchFunction(cursor);

        // store next cursor in history if going forward
        if (this.pageCursors.length === this.currentPageIndex + 1) {
            this.pageCursors.push(nextCursor);
        }

        return items;
    }

    /** Move to the next page */
    async nextPage(): Promise<T[] | null> {
        if (this.pageCursors[this.currentPageIndex + 1] === undefined) {
            return null; // no next page
        }
        this.currentPageIndex++;
        return this.getCurrentPage();
    }

    /** Move to the previous page */
    async prevPage(): Promise<T[] | null> {
        if (this.currentPageIndex === 0) return null;
        this.currentPageIndex--;
        return this.getCurrentPage();
    }

    /** Jump to a specific page index */
    async goToPage(pageIndex: number): Promise<T[] | null> {
        if (pageIndex < 0 || pageIndex >= this.pageCursors.length) return null;
        this.currentPageIndex = pageIndex;
        return this.getCurrentPage();
    }

    pageCount(): number {
        return this.pageCursors.length
    }
    /** Returns the current page index (0-based) */
    getCurrentPageIndex(): number {
        return this.currentPageIndex;
    }
}
