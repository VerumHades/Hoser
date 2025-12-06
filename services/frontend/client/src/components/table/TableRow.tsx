import type { Clickable, HasChildren, HasClassname } from "../common";
import { Children } from "react";

export default function TableRow({children, onClick, className}: HasChildren & Clickable & HasClassname){
    return <div className={"flex md:flex-row gap-4 flex-col p-4 transition-all dark:hover:bg-indigo-800 hover:bg-indigo-300 " + className} onClick={onClick}>
        {Children.toArray(children).map(child => 
            <div className="flex-1">
                {child}
            </div>
        )}
    </div>
}