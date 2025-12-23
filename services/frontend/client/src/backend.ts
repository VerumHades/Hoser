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

export interface Money {
    amount: number;      // corresponds to Go Value
    code: string;      // corresponds to Go Short
}

export interface PricingEntry {
    id: string;          // "0", "1", etc
    type: 0 | 1;         // 0 = OneTime, 1 = Monthly
    currency: Money;
}

export interface HardwareSpecification {
    cpu?: number;   // corresponds to Go CPU
    ramBytes?: number;   // corresponds to Go RAM (ramBytes)
    diskBytes?: number;  // corresponds to Go Disk (diskBytes)
}

export interface ListingRequest {
    id: string;
    title?: string;
    description?: string;
    accessMode?: number;
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
    accessMode?: number;
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

export interface ApiBillingAccount {
    id: string;
    status: "active" | "suspended" | "closed";
    paymentProvider: string;
    providerAccountId: string;
    createdAt: number; // Unix timestamp
}


export type ApiPaymentKind = "one_time" | "subscription" | "unknown";
export type ApiPaymentStatus = "pending" | "paid" | "failed" | "refunded";

export interface ApiPayment {
    id: string;
    billingAccountId: string;
    amount: number;
    currency: string;
    status: ApiPaymentStatus;
    kind: ApiPaymentKind;
    createdAt: number; // unix timestamp
    paidAt?: number;
}

// =================== TYPES ===================
export interface ApiInstance {
    id: string;
    listingId: string;
    billingId: string;
    state: string;
    hardwareSpecification: HardwareSpecification;
}

export interface LaunchInstanceRequest {
    listingId: string;
    billingAccountId: string;
    hardwareSpecification: HardwareSpecification;
}

export interface UpdateHardwareRequest {
    hardwareSpecification: HardwareSpecification;
}

export type HasInLibraryResponse = {
    hasListing: boolean;
};

// =================== EXTENDED API ===================
export const API = {
    user: {
        async getData(): Promise<User | undefined> {
            const response = await backend_request("/user/data", "GET");
            return response.ok ? (response.json as User) : undefined
        },
        library:{
            async add(id: string) {
                return await backend_request("/user/library", "POST", { id });
            },
            async delete(id: string) {
                return await backend_request("/user/library", "DELETE", { id });
            },
            async check(listingId: string): Promise<APIResult<HasInLibraryResponse>> {
                return await backend_request(`/user/library/${listingId}/exists`, "GET");
            }
        },
        billing: {
            async list(): Promise<APIResult<ApiBillingAccount[]>> {
                return await backend_request("/user/billing", "GET");
            },
            async create(paymentProvider: string, providerAccountID: string): Promise<APIResult<ApiBillingAccount>> {
                const data = { paymentProvider, providerAccountID };
                return await backend_request("/user/billing", "POST", data);
            },
            async get(id: string): Promise<APIResult<ApiBillingAccount>> {
                return await backend_request(`/user/billing/${id}`, "GET");
            },
            async suspend(id: string): Promise<APIResult<ApiBillingAccount>> {
                return await backend_request(`/user/billing/${id}/suspend`, "POST");
            },
            async close(id: string): Promise<APIResult<ApiBillingAccount>> {
                return await backend_request(`/user/billing/${id}/close`, "POST");
            },
            payments: {
                async list(billingAccountId: string): Promise<APIResult<ApiPayment[]>> {
                    return await backend_request(`/user/billing/${billingAccountId}/payments`, "GET");
                },
                async get(billingAccountId: string, paymentId: string): Promise<APIResult<ApiPayment>> {
                    return await backend_request(`/user/billing/${billingAccountId}/payment/${paymentId}`, "GET");
                },
                async getMetadata(billingAccountId: string, paymentId: string): Promise<APIResult<any>> {
                    return await backend_request(`/user/billing/${billingAccountId}/payment/${paymentId}/metadata`, "GET");
                }
            }
        },
        instances: {
            async launch(listingId: string, billingAccountId: string, hardwareSpecification: HardwareSpecification): Promise<APIResult<ApiInstance>> {
                const data: LaunchInstanceRequest = { listingId, billingAccountId, hardwareSpecification };
                return await backend_request(`/user/instances`, "POST", data);
            },
            async updateHardware(instanceId: string, hardwareSpecification: HardwareSpecification): Promise<APIResult<ApiInstance>> {
                const data: UpdateHardwareRequest = { hardwareSpecification };
                return await backend_request(`/user/instances/${instanceId}/hardware`, "PATCH", data);
            },
            async get(instanceId: string): Promise<APIResult<ApiInstance>> {
                return await backend_request(`/user/instances/${instanceId}`, "GET");
            },
            async listByBilling(billingId: string): Promise<APIResult<ApiInstance[]>> {
                return await backend_request(`/user/billing/${billingId}/instances`, "GET");
            },
            async list(): Promise<APIResult<ApiInstance[]>> {
                return await backend_request(`/user/instances`, "GET");
            }
        }
    },
    listing: {
        async get(id: string): Promise<APIResult<Listing>> {
            return await backend_request(`/listing/${id}`, "GET");
        },
    },
    developer: {
        listing: {
            async get(id: string): Promise<APIResult<DeveloperListing>> {
                return await backend_request("/developer/listing/" + id, "GET");
            },
            async delete(id: string): Promise<APIResult<unknown>> {
                return await backend_request("/developer/listing", "DELETE", {id});
            },
            async update(data: ListingRequest) {
                return await backend_request("/developer/listing", "PUT", data)
            },
            async create(title: string = "My New Listing", description: string = "This is a description of my listing.") {
                const data = { title, description };
                return await backend_request<DeveloperListing>(`/developer/listing`, "POST", data)
            },
        }
    }
};