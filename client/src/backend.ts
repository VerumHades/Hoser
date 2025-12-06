import backend_constants from "./backend_constants";

async function backend_request<T>(endpoint: string, type: string, body?: unknown): Promise<APIResult<T>> {
    const request_body: RequestInit = {
        method: type,
        headers: {
            "Content-Type": "application/json",
        },
        credentials: "include"
    }
    if (body)
        request_body.body = JSON.stringify(body);

    const response = await fetch(`${backend_constants.address}${endpoint}`, request_body);

    let json = undefined
    try {
        json = await response.json();
    }
    catch (err) {
        console.error(err)
    }

    return {
        response,
        ok: response.ok,
        json
    }
}


export const ListingAccessModes = {
    0: { label: 'Private', description: "Private, only available to you." },
    1: { label: 'Public', description: "Available to everyone." },
}

export const BillingFrequency = {
  OneTime: 0,
  Monthly: 1,
  Yearly: 2,
} as const;

export const BillingFrequencyName: Record<BillingFrequency, string> = {
  0: "OneTime",
  1: "Monthly",
  2: "Yearly",
};

export type BillingFrequency =
  (typeof BillingFrequency)[keyof typeof BillingFrequency];

export interface CurrencyRequest {
    value: number;      // corresponds to Go Value
    name: string;       // corresponds to Go Name
    short: string;      // corresponds to Go Short
}

export interface PricingEntry {
    id: string;          // "0", "1", etc
    type: 0 | 1;         // 0 = OneTime, 1 = Monthly
    currency: CurrencyRequest;
}

export interface HardwareUpdate {
    cpu?: number;   // corresponds to Go CPU
    ram?: number;   // corresponds to Go RAM (ramBytes)
    disk?: number;  // corresponds to Go Disk (diskBytes)
}

export interface ListingRequest {
    id: string;
    title?: string;
    description?: string;
    accessMode?: number;
    hardware?: HardwareUpdate;
}

export interface Listing {
    id: string;
    title?: string;
    description?: string;
    author: string;
    accessMode?: number;
    prices: PricingEntry[];
    hardware?: HardwareUpdate;
}

export type User = {
    username: string
    isDeveloper: boolean
}

interface APIResult<T> {
    ok: boolean,
    response: Response,
    json: T
}

export const API = {
    user: {
        async getData(): Promise<User | undefined> {
            const response = await backend_request("/user/data", "GET");
            return response.ok ? (response.json as User) : undefined
        },
        async rentListing(id: string) {
            return await backend_request("/user/rent", "POST", { id });
        }
    },
    developer: {
        listing: {
            async get(id: string): Promise<APIResult<Listing>> {
                return await backend_request("/developer/listing", "GET", {id});
            },
            async delete(id: string): Promise<APIResult<unknown>> {
                return await backend_request("/developer/listing", "DELETE", {id});
            },
            async update(data: ListingRequest) {
                return await backend_request("/developer/listing", "PUT", data)
            },
            async create(title: string = "My New Listing", description: string = "This is a description of my listing.") {
                const data = {
                    title,
                    description
                };

                return await backend_request(`/developer/listing`, "POST", data)
            },
            prices: {
                async delete(listing_id: string, price_id: string) {
                    return await backend_request("/developer/listing/price", "DELETE", {
                        listingId: listing_id,
                        pricingId: price_id,
                    })
                },
                async update(listing_id: string, price: PricingEntry) {
                    return await backend_request("/developer/listing/price", "PUT", {
                        listingId: listing_id,
                        pricingId: price.id,
                        ...price
                    })
                },
            }
        }


    }
}