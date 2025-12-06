import { useNavigate } from "react-router-dom"

export default function Logo() {
    let navigate = useNavigate()

    return <div className="flex flex-row items-center" onClick={() => navigate("/")}>
        <img src="/ghost_logo.svg" className="h-8 w-auto" alt="Logo" />
        <span className="text-2xl font-bold text-gray-900 dark:text-white">
            gHost
        </span>
    </div>
}