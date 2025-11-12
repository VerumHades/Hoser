import backend_constants from "./backend_constants"


export type User = {
    username: string
    isDeveloper: boolean
}

export function clearCachedUser() {
    localStorage.removeItem("user")
}

export async function getUser(): Promise<User | undefined> {
    let localUser = localStorage.getItem("user") 

    if(localUser != null) 
        return JSON.parse(localUser) as User;

    try {
        const response = await fetch(`${backend_constants.address}/user/data`, {
            method: "GET",
            credentials: "include"
        })
        if (response.status === 200) {
            let user = await response.json() as User
            localStorage.setItem("user", JSON.stringify(user))
            return user
        }
        localStorage.setItem("user", "undefined")
    } catch {
        
    }
}