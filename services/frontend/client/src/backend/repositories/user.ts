
import { backendRequest } from "../backend_request";
import type { CursorPaginatedResult } from "../paginator";
import type { Listing } from "./listing";

export type User = {
    username: string;
    isDeveloper: boolean;
};


export const UserAPI = {
    /**
     * Retrieves the profile and account data for the currently authenticated user.
     */
    async getData(): Promise<User> {
        return await backendRequest<User>("/user/data", "GET");
    },

    /**
     * Executes the purchase process for a listing, creating the necessary ledger transactions.
     */
    async purchaseListing(listingId: string): Promise<void> {
        return await backendRequest(`/user/listings/${listingId}/purchase`, "POST");
    },

    /**
     * Checks the ledger and settlement status to determine if the user owns the listing.
     */
    async isOwner(listingId: string): Promise<boolean> {
        const response = await backendRequest<{ isOwner: boolean }>(
            `/user/listings/${listingId}/ownership`, 
            "GET"
        );
        return response?.isOwner ?? false;
    },

    /**
     * Checks the ledger and settlement status to determine if the user owns the listing.
     */
    async isPurchaseProcessingOwner(listingId: string): Promise<boolean> {
        const response = await backendRequest<{ isProcessing: boolean }>(
            `/user/listings/${listingId}/processing`, 
            "GET"
        );
        return response?.isProcessing ?? false;
    },

    library: {
        /**
         * Adds a listing to the user's personal library collection.
         */
        async add(id: string) { 
            return backendRequest("/user/library", "POST", { id }); 
        },

        /**
         * Removes a listing from the user's personal library collection.
         */
        async delete(id: string) { 
            return backendRequest("/user/library", "DELETE", { id }); 
        },

        /**
         * Verifies if a specific listing exists within the user's library.
         */
        async check(listingId: string): Promise<boolean> { 
            const response = await backendRequest(`/user/library/${listingId}`, "GET") as { hasListing: boolean }; 
            return response.hasListing;
        },

        /**
         * Returns a paginated list of all listings currently in the user's library.
         */
        async listBatch(cursor?: string): Promise<CursorPaginatedResult<Listing>> {
            const parameters = cursor ? `?cursor=${cursor}` : "";
            return backendRequest<CursorPaginatedResult<Listing>>(`/user/library${parameters}`, "GET");
        },
    },
};