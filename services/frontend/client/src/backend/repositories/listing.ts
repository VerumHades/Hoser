import { backendRequest } from "../backend_request";
import type { HardwareSpecification } from "../types";
import type { CursorPaginatedResult } from "../paginator";

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
    async get(id: string) { return backendRequest<Listing>(`/listings/${id}`, "GET"); },

    async search(cursor?: string): Promise<CursorPaginatedResult<Listing>> {
        const params = cursor ? `?cursor=${cursor}` : "";
        return backendRequest<CursorPaginatedResult<Listing>>(`/search/listings${params}`, "GET");
    },

    async getScreenshotReadURL(listingID: string, screenshotID: string){
        return backendRequest<GetScreenshotResponse>(
            `/listings/${listingID}/screenshots/${screenshotID}`,
            "GET"
        )
    },
};
