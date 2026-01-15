import MainNavbar from "../navigation/MainNavbar"

interface WithNavbarProps {
    children: React.ReactNode
}


export default function WithNavbar({ children }: WithNavbarProps) {
    return <div className="absolute flex flex-col inset-0 justify-start">
        <MainNavbar></MainNavbar>
        <div className="absolute w-screen top-16 flex flex-col items-center  bg-slate-50 dark:bg-gray-950">
            {children}
        </div>
    </div>
}