import backend_constants from "./backend_constants";

export async function backend_request(endpoint: string, type: string, body: any) {
    return await fetch(`${backend_constants.address}${endpoint}`, {
        method: type,
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(body),
        credentials: "include" 
    });
}
