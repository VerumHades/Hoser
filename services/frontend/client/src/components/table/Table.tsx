import type { HasChildren, HasClassname } from "../common";



export default function Table({children, className} : HasChildren & HasClassname){
    return <div className={
            `dark:[&>*:nth-child(odd)]:bg-slate-800 
            dark:[&>*:nth-child(even)]:bg-slate-900 
            [&>*:nth-child(odd)]:bg-slate-100
            [&>*:nth-child(even)]:bg-slate-200 
            ` + className}>
        {children}
    </div>
}