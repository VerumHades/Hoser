import { ChevronRight, Image } from "lucide-react";
import { useEffect } from "react";

interface ListElementProps {
    children?: React.ReactNode
    icon?: string
    onClick?: () => void
}

export default function ListElement({children, icon, onClick}: ListElementProps) {
    useEffect(() => {
        document.title = "GHost";
    }, []);
    
    return (
        <div
            className="flex items-center gap-4 p-4 shadow-lg text-gray-800 dark:text-gray-10
                bg-slate-200 hover:bg-slate-300 dark:bg-gray-900 dark:hover:bg-gray-800 transition-all hover:translate-x-2 mx-2 mb-5"
        >
            {
                icon ? <img
                    src={icon}
                    alt=""
                    className="w-12 h-12 object-cover rounded-xl"
                /> : <Image size={50} />
            }
            <div className="flex flex-row justify-between items-center w-full" onClick={onClick}>
                <div className="flex-1">
                    {children}
                </div>
                <div>
                    <ChevronRight size={24}></ChevronRight>
                </div>
            </div>
        </div>
    );
}
