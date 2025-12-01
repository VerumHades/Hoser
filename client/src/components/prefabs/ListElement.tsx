import { ChevronRight, Edit } from "lucide-react";
import { useEffect } from "react";

interface ListElementProps {
    children?: React.ReactNode
    onClick?: () => void
}

export default function ListElement({children, onClick}: ListElementProps) {
    useEffect(() => {
        document.title = "GHost";
    }, []);
    
    return (
        <div
            onClick={onClick}
            className="
            
            flex flex-col items-center shadow-lg text-gray-800 dark:text-gray-300 p-5 w-full max-h-50
                bg-slate-200 hover:bg-slate-300 dark:bg-gray-900 dark:hover:bg-gray-800 transition-all hover:translate-x-2"
        >
            {children}
        </div>
    );
}
