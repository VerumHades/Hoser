import backend_constants from "./constants";

export async function backendRequest<T>(endpoint: string, method: string, body?: unknown): Promise<T> {
    const requestInit: RequestInit = {
        method,
        headers: { "Content-Type": "application/json" },
        credentials: "include",
    };

    if (body) requestInit.body = JSON.stringify(body);

    try {
        const response = await fetch(`${backend_constants.address}${endpoint}`, requestInit);
        return await response.json();
    } catch (err) {
        throw err;
    }
}