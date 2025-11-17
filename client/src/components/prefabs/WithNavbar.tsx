import Navbar from "../../Navbar"

interface WithNavbarProps {
    children: React.ReactNode
}


export default function WithNavbar({ children }: WithNavbarProps) {
    return <div className="absolute flex flex-col inset-0 justify-start">
        <Navbar></Navbar>
        <div className="absolute inset-0 top-16 flex flex-col items-center">
            {children}
        </div>
    </div>
}