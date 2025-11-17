import backend_constants from "./backend_constants";

async function backend_request(endpoint: string, type: string, body?: any): Promise<APIResult> {

    let request_body: RequestInit = {
        method: type,
        headers: {
            "Content-Type": "application/json",
        },
        credentials: "include"
    }
    if (body)
        request_body.body = JSON.stringify(body);

    const response = await fetch(`${backend_constants.address}${endpoint}`, request_body);

    return {
        response,
        ok: response.ok,
        json: await response.json()
    }
}


export const ListingAccessModes = {
    0: { label: 'Private', description: "Private, only available to you." },
    1: { label: 'Public', description: "Available to everyone." },
}

interface DeveloperListingUpdateData {
    id?: string,
    title?: string,
    description?: string,
    accessMode?: keyof typeof ListingAccessModes
}


export type User = {
    username: string
    isDeveloper: boolean
}

interface APIResult {
    ok: boolean,
    response: Response,
    json: any
}

export const API = {
    user: {
        async getData(): Promise<User | undefined> {
            let response = await backend_request("/user/data", "GET");
            return response.ok ? (response.json as User) : undefined
        }
    },
    developer: {
        async deleteListing(id: string): Promise<APIResult> {
            return await backend_request("/developer/listing", "DELETE", { id });
        },
        async editListing(id: string, data: DeveloperListingUpdateData) {
            return await backend_request("/developer/listing", "PUT", { ...data, id })
        },
        async createListing(title: string = "My New Listing", description: string = "This is a description of my listing.") {
            const data = {
                title,
                description
            };
            return await backend_request(`/developer/listing`, "POST", data)
        },
    }
}