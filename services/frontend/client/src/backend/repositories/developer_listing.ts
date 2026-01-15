import { backendRequest } from "../backend_request";
import type { CursorPaginatedResult } from "../paginator";
import type { HardwareSpecification, ListingAccessMode } from "../types";
import type { Listing } from "./listing";

export interface DeveloperListing extends Listing {
    accessMode?: ListingAccessMode;
}

export interface UpdateListingRequest {
    id: string;
    title?: string;
    description?: string;
    hardware?: HardwareSpecification;
    price?: number;
    accessMode?: ListingAccessMode;
    documentation_markdown?: string;
}

export interface ListingGithubSetup {
    repositoryURL: string;
}

export interface AttachOrUpdateGithubSetupRequest {
    repositoryURL: string;
    accessToken: string;
}

export interface CreateScreenshotResponse {
    uploadUrl: string;
}

export const DeveloperListingAPI = {
    /**
     * Fetch a developer listing by id
     */
    async get(listingId: string): Promise<DeveloperListing> {
        return backendRequest<DeveloperListing>(
            `/developer/listings/${listingId}`,
            "GET"
        );
    },

    /**
     * Create a new developer listing
     */
    async create(
        title: string,
        description: string
    ): Promise<DeveloperListing> {
        return backendRequest<DeveloperListing>(
            "/developer/listings",
            "POST",
            { title, description }
        );
    },


    /**
     * Update an existing developer listing
     */
    async update(
        updateListingRequest: UpdateListingRequest
    ): Promise<DeveloperListing> {
        return backendRequest<DeveloperListing>(
            "/developer/listings",
            "PUT",
            updateListingRequest
        );
    },

    async list(cursor?: string): Promise<CursorPaginatedResult<DeveloperListing>> {
        const params = cursor ? `?cursor=${cursor}` : "";
        return backendRequest<CursorPaginatedResult<DeveloperListing>>(`/developer/listings${params}`, "GET");
    },
    /**
     * Delete a developer listing
     */
    async delete(listingId: string): Promise<void> {
        return backendRequest<void>(
            "/developer/listings",
            "DELETE",
            { id: listingId }
        );
    },

    screenshots: {
        /**
         * Requests a signed URL to upload a screenshot.
         * The backend will return a MinIO/S3 URL that allows a direct PUT request.
         */
        async createUploadUrl(listingId: string): Promise<CreateScreenshotResponse> {
            return backendRequest<CreateScreenshotResponse>(
                `/developer/listings/${listingId}/screenshots`,
                "POST"
            );
        },

        /**
         * Deletes a screenshot from both the database and object storage.
         */
        async delete(listingId: string, screenshotId: string): Promise<void> {
            return backendRequest<void>(
                `/developer/listings/${listingId}/screenshots/${screenshotId}`,
                "DELETE"
            );
        }
    },

    githubSetup: {
        /**
         * Get GitHub setup for a listing
         */
        async get(listingId: string): Promise<ListingGithubSetup> {
            return backendRequest<ListingGithubSetup>(
                `/developer/listings/${listingId}/setup/github`,
                "GET"
            );
        },

        /**
         * Attach or update GitHub setup for a listing
         */
        async attachOrUpdate(
            listingId: string,
            repositoryURL: string,
            accessToken: string
        ): Promise<ListingGithubSetup> {
            const attachOrUpdateGithubSetupRequest: AttachOrUpdateGithubSetupRequest = {
                repositoryURL,
                accessToken
            };

            return backendRequest<ListingGithubSetup>(
                `/developer/listings/${listingId}/setup/github`,
                "POST",
                attachOrUpdateGithubSetupRequest
            );
        },

        /**
         * Remove GitHub setup from a listing
         */
        async remove(listingId: string): Promise<void> {
            return backendRequest<void>(
                `/developer/listings/${listingId}/setup/github`,
                "DELETE"
            );
        }
    }
};