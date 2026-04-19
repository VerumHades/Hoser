
export interface NumericRange {
    min?: number;
    max?: number;
}

export interface DateRange {
    from?: string; // ISO 8601 format
    to?: string;   // ISO 8601 format
}

export interface ListingSearchQuery {
    text?: string;
    cpu?: NumericRange;
    ramBytes?: NumericRange;
    diskBytes?: NumericRange;
    price?: NumericRange;
    createdAt?: DateRange;
    authorId?: string;
}

/**
 * buildSearchParams serializes the search query object into a URL-safe string.
 */
export function buildSearchParams(query: ListingSearchQuery, cursor?: string): string {
    const urlParams = new URLSearchParams();

    if (cursor) urlParams.append("cursor", cursor);
    if (query.text) urlParams.append("q", query.text);
    if (query.authorId) urlParams.append("author_id", query.authorId);

    // Flatten numeric ranges into the query format expected by the backend
    appendNumericRange(urlParams, query.cpu, "cpu");
    appendNumericRange(urlParams, query.ramBytes, "ram");
    appendNumericRange(urlParams, query.diskBytes, "disk");
    appendNumericRange(urlParams, query.price, "price");
    
    // Flatten date ranges
    appendDateRange(urlParams, query.createdAt);

    return urlParams.toString();
}

function appendNumericRange(params: URLSearchParams, range: NumericRange | undefined, prefix: string): void {
    if (!range) return;
    if (range.min !== undefined) params.append(`${prefix}_min`, range.min.toString());
    if (range.max !== undefined) params.append(`${prefix}_max`, range.max.toString());
}

function appendDateRange(params: URLSearchParams, range: DateRange | undefined): void {
    if (!range) return;
    if (range.from) params.append("date_from", range.from);
    if (range.to) params.append("date_to", range.to);
}