
import { backendRequest } from "../backend_request";
import type { CursorPaginatedResult } from "../paginator";
import type { DeveloperListing } from "./developer_listing";
import type { Listing } from "./listing";

export type User = {
    username: string;
    isDeveloper: boolean;
};


export const UserAPI = {
    async getData(): Promise<User> {
        return await backendRequest<User>("/user/data", "GET");
    },
    library: {
        async add(id: string) { return backendRequest("/user/library", "POST", { id }); },
        async delete(id: string) { return backendRequest("/user/library", "DELETE", { id }); },
        async check(listingId: string): Promise<boolean> { 
            const response = await backendRequest(`/user/library/${listingId}`, "GET") as {hasListing: boolean}; 
            return response.hasListing
        },
        async listBatch(cursor?: string): Promise<CursorPaginatedResult<Listing>> {
            const params = cursor ? `?cursor=${cursor}` : "";
            return backendRequest<CursorPaginatedResult<Listing>>(`/user/library${params}`, "GET");
        },
    },
};