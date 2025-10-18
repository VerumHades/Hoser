import backend_constants from "./backend_constants"


type User = {
    Username: string
}

let user: User | undefined = undefined

try {
    let response = await fetch(`${backend_constants.address}/user/data`, {
        method: "GET",
        credentials: "include"
    })

    if (response.status == 200) {
        user = await response.json() as User
    }
}
catch (err) {

}

console.log(user)
export default user
