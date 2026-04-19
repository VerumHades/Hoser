import { backendRequest } from "../backend_request";
import type { HardwareSpecification } from "../types";
import type { CursorPaginatedResult } from "../paginator";
import { HttpError } from "../http_error";
import { buildSearchParams, type ListingSearchQuery } from "../utils/query";

export interface Listing {
    id: string;
    title?: string;
    description?: string;
    author: string;
    price: number;
    recommended_hardware?: HardwareSpecification;
    screenshotIds?: string[]; // Add this
    documentation_markdown: string;
}

interface GetScreenshotResponse {
    readUrl: string,
}

export const ListingAPI = {
    async get(id: string) { 
        return backendRequest<Listing>(`/listings/${id}`, "GET"); 
    },

    /**
     * search executes a complex query against the listing index.
     * @param query - The filter and text parameters.
     * @param cursor - The pagination pointer.
     */
    async search(query: ListingSearchQuery, cursor?: string): Promise<CursorPaginatedResult<Listing> | null> {
        const params = buildSearchParams(query, cursor);
        try {
            return backendRequest<CursorPaginatedResult<Listing>>(`/search/listings?${params}`, "GET");
        }
        catch(error){
            return null
        }
        
    },

    async getScreenshotReadURL(listingID: string, screenshotID: string) {
        return backendRequest<GetScreenshotResponse>(
            `/listings/${listingID}/screenshots/${screenshotID}`,
            "GET"
        );
    },
};