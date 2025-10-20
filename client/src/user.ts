import backend_constants from "./backend_constants"


export type User = {
    Username: string
    IsDeveloper: boolean
}

export async function getUser(): Promise<User | undefined> {
    try {
        const response = await fetch(`${backend_constants.address}/user/data`, {
            method: "GET",
            credentials: "include"
        })
        if (response.status === 200) {
            return await response.json() as User
        }
    } catch {
        // handle errors
    }
}