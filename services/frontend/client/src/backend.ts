import backend_constants from "./backend_constants";

async function backend_request<T>(endpoint: string, method: string, body?: unknown): Promise<APIResult<T>> {
    const requestInit: RequestInit = {
        method,
        headers: { "Content-Type": "application/json" },
        credentials: "include",
    };

    if (body) requestInit.body = JSON.stringify(body);

    try {
        const response = await fetch(`${backend_constants.address}${endpoint}`, requestInit);
        let json: T | undefined;
        try {
            json = await response.json();
        } catch {}

        return { response, ok: response.ok, json };
    } catch (err) {
        throw err;
    }
}

// =================== ENUMS & TYPES ===================
export const ListingAccessModes = {
    Private: 0,
    Public: 1,
} as const;

export type ListingAccessMode = (typeof ListingAccessModes)[keyof typeof ListingAccessModes];

export interface Money {
    amount: number;
    code: string;
}

export interface HardwareSpecification {
    cpu?: number;
    ramBytes?: number;
    diskBytes?: number;
}

export interface ListingRequest {
    id: string;
    title?: string;
    description?: string;
    accessMode?: ListingAccessMode;
    hardware?: HardwareSpecification;
}

export interface Listing {
    id: string;
    title?: string;
    description?: string;
    author: string;
    price: Money;
    hardware?: HardwareSpecification;
}

export interface DeveloperListing extends Listing {
    accessMode?: ListingAccessMode;
}

export type User = {
    username: string;
    isDeveloper: boolean;
};

interface APIResult<T> {
    ok: boolean;
    response: Response;
    json?: T;
}

// =================== NEW HARDWARE RATE TYPES ===================
export interface ApiHardwareRate {
    id: string;
    resourceType: "CPU" | "RAM" | "DISK";
    costInCents: number;
    validFromTime: string;
}

export interface ApiHardwareRatesResponse {
    rates: ApiHardwareRate[];
}

// =================== LISTING TYPES ===================
export interface ApiListingBase {
    id: string;
    title?: string;
    description?: string;
    price: number;
    hardware?: HardwareSpecification;
}

export interface ApiDeveloperListing extends ApiListingBase {
    accessMode?: ListingAccessMode;
}

export interface AddListingRequest {
    title: string;
    description: string;
}

export interface UpdateListingRequest {
    id: string;
    title?: string;
    description?: string;
    hardware?: HardwareSpecification;
    price?: number;
    accessMode?: ListingAccessMode;
}

// =================== GITHUB SETUP ===================
export interface ListingGithubSetup {
    repositoryURL: string;
}

export interface AttachOrUpdateGithubSetupRequest {
    repositoryURL: string;
    accessToken: string;
}

// =================== API CLIENT ===================
export const API = {
    hardware: {
        async getRates(): Promise<APIResult<ApiHardwareRatesResponse>> {
            return backend_request("/hardware/rates", "GET");
        },
    },

    user: {
        async getData(): Promise<User | undefined> {
            const response = await backend_request<User>("/user/data", "GET");
            return response.ok ? response.json : undefined;
        },
        library: {
            async add(id: string) {
                return backend_request("/user/library", "POST", { id });
            },
            async delete(id: string) {
                return backend_request("/user/library", "DELETE", { id });
            },
            async check(listingId: string): Promise<APIResult<{ hasListing: boolean }>> {
                return backend_request(`/user/library/${listingId}/exists`, "GET");
            },
            async listBatch(cursor?: string) {
                const params = cursor ? `?cursor=${cursor}` : "";
                return backend_request<{ items: DeveloperListing[]; nextCursor?: string }>(`/user/library${params}`, "GET");
            },
        },
    },

    listing: {
        async get(id: string): Promise<APIResult<Listing>> {
            return backend_request(`/listings/${id}`, "GET");
        },
        async listBatch(cursor?: string) {
            const params = cursor ? `?cursor=${cursor}` : "";
            return backend_request<{ items: Listing[]; nextCursor?: string }>(`/listing${params}`, "GET");
        },
    },

    developer: {
        listing: {
            async get(id: string): Promise<APIResult<DeveloperListing>> {
                return backend_request(`/developer/listings/${id}`, "GET");
            },
            async create(title: string, description: string) {
                return backend_request<DeveloperListing>("/developer/listings", "POST", { title, description });
            },
            async update(data: UpdateListingRequest) {
                return backend_request("/developer/listings", "PUT", data);
            },
            async delete(id: string) {
                return backend_request("/developer/listings", "DELETE", { id });
            },
            setup: {
                async get(listingId: string): Promise<APIResult<ListingGithubSetup>> {
                    return backend_request(`/developer/listings/${listingId}/setup/github`, "GET");
                },
                async attachOrUpdate(listingId: string, repositoryURL: string, accessToken: string) {
                    const payload: AttachOrUpdateGithubSetupRequest = { repositoryURL, accessToken };
                    return backend_request(`/developer/listings/${listingId}/setup/github`, "POST", payload);
                },
                async remove(listingId: string) {
                    return backend_request(`/developer/listings/${listingId}/setup/github`, "DELETE");
                },
            },
        },
    },
};
