import backend_constants from "./constants";
import { HttpError } from "./http_error";

export async function backendRequest<T>(endpoint: string, method: string, body?: unknown): Promise<T | null> {
    const requestInit: RequestInit = {
        method,
        headers: { "Content-Type": "application/json" },
        credentials: "include",
    };

    if (body) requestInit.body = JSON.stringify(body);

    const response = await fetch(`${backend_constants.address}${endpoint}`, requestInit);

    try {
        const json = await response.json();

        if(!response.ok) 
            throw new HttpError(response.status, json.error);

        return json
    }
    catch(err) {
        console.log(err)
        return null;
    }

}