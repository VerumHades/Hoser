import backend_constants from "./backend_constants";

async function backend_request(endpoint: string, type: string, body: any) {
    let response = await fetch(`${backend_constants.address}${endpoint}`, {
        method: type,
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(body),
        credentials: "include"
    });

    return {
        response,
        json: await response.json()
    }
}


export const ListingAccessModes = {
    0: {label: 'Private', description: "Private, only available to you."},
    1: {label: 'Public', description: "Available to everyone."},
}

interface DeveloperListingUpdateData {
    id?: string,
    title?: string,
    description?: string,
    accessMode?: keyof typeof ListingAccessModes
}

interface APIResult {
    ok: boolean,
    json: any
}

export const API = {
    developer: {
        deleteListing: async (id: string): Promise<APIResult> => {
            let request = (await backend_request("/developer/listing", "DELETE", {id}));

            return {
                ok: request.response.ok,
                json: request.json
            }
        },
        editListing: async (id: string, data: DeveloperListingUpdateData) => {
            let request = (await backend_request("/developer/listing", "PUT", {...data, id}))

            return {
                ok: request.response.ok,
                json: request.json
            }
        },
        createListing: async (title: string = "My New Listing", description: string = "This is a description of my listing.") => {
            const data = {
                title,
                description
            };

            let request = (await backend_request(`/developer/listing`, "POST", data))

            return {
                ok: request.response.ok,
                json: request.json
            }
        },
    }
}