import MainNavbar from "../navigation/MainNavbar"

interface WithNavbarProps {
    children: React.ReactNode
}


export default function WithNavbar({ children }: WithNavbarProps) {
    return <div className="h-full flex flex-col">
        <MainNavbar></MainNavbar>
        <div className="flex flex-1 flex-col items-center  bg-slate-50 dark:bg-gray-950">
            {children}
        </div>
    </div>
}