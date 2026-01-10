import backend_constants from "./constants";
import { HttpError } from "./http_error";

export async function backendRequest<T>(endpoint: string, method: string, body?: unknown): Promise<T> {
    const requestInit: RequestInit = {
        method,
        headers: { "Content-Type": "application/json" },
        credentials: "include",
    };

    if (body) requestInit.body = JSON.stringify(body);

    const response = await fetch(`${backend_constants.address}${endpoint}`, requestInit);
    const json = await response.json();

    if(response.status != 200) 
        throw new HttpError(response.status, json.error);
    
    return json
}