export interface CursorPaginatedResult<T> {
    items: T[];
    cursor?: string;
}

/**
 * A stateful paginator that optimizes navigation by using cached cursors.
 * Proactively discovers future cursors based on a lookahead parameter.
 */
export class CursorPaginator<T> {
    private readonly fetchFunction: (cursor?: string) => Promise<CursorPaginatedResult<T> | null>;
    private readonly pageCursors: (string | undefined)[] = [undefined];
    private currentPageIndex: number = 0;
    private readonly lookahead: number;

    constructor(
        fetchFunction: (cursor?: string) => Promise<CursorPaginatedResult<T> | null>,
        lookahead: number = 0
    ) {
        this.fetchFunction = fetchFunction;
        this.lookahead = lookahead;
    }

    /**
     * Retrieves the items for the current page index and triggers lookahead discovery.
     */
    async getCurrentPage(): Promise<T[]> {
        const result = await this.fetchByCursor(this.pageCursors[this.currentPageIndex]);

        if (!result) return [];

        this.updateCache(this.currentPageIndex + 1, result.cursor);
        this.runLookaheadDiscovery();

        return result.items;
    }

    async nextPage(): Promise<T[] | null> {
        if (!this.canNavigateForward()) return null;

        this.currentPageIndex++;
        return this.getCurrentPage();
    }

    async prevPage(): Promise<T[] | null> {
        if (this.isAtStart()) return null;

        this.currentPageIndex--;
        return this.getCurrentPage();
    }

    async goToPage(targetPageIndex: number): Promise<T[] | null> {
        if (targetPageIndex < 0) return null;

        if (this.isIndexCached(targetPageIndex)) {
            this.currentPageIndex = targetPageIndex;
            return this.getCurrentPage();
        }

        return this.discoverPagesUntilTarget(targetPageIndex);
    }

    pageCount(): number {
        return this.pageCursors.length;
    }

    getCurrentPageIndex():number {
        return this.currentPageIndex;
    }

    private async fetchByCursor(cursor?: string): Promise<CursorPaginatedResult<T> | null> {
        return await this.fetchFunction(cursor);
    }

    private async runLookaheadDiscovery(): Promise<void> {
        for (let i = 0; i < this.lookahead; i++) {
            const lastKnownIndex = this.pageCursors.length - 1;
            const lastCursor = this.pageCursors[lastKnownIndex];

            if (lastCursor === undefined && lastKnownIndex > 0) break;

            const result = await this.fetchByCursor(lastCursor);
            if (!result || !result.cursor) break;

            this.updateCache(lastKnownIndex + 1, result.cursor);
        }
    }

    private updateCache(index: number, nextCursor?: string): void {
        if (index === this.pageCursors.length && nextCursor !== undefined) {
            this.pageCursors.push(nextCursor);
        }
    }

    private async discoverPagesUntilTarget(targetPageIndex: number): Promise<T[] | null> {
        while (this.currentPageIndex < targetPageIndex) {
            const result = await this.nextPage();
            if (!result) return null;
        }
        return this.getCurrentPage();
    }

    private isIndexCached(index: number): boolean {
        return index < this.pageCursors.length;
    }

    private canNavigateForward(): boolean {
        return this.isIndexCached(this.currentPageIndex + 1) || 
               this.pageCursors[this.currentPageIndex] !== undefined;
    }

    private isAtStart(): boolean {
        return this.currentPageIndex === 0;
    }
}