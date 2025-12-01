

interface TableProps {
    children?: React.ReactNode
}

export default function Table({children} : TableProps){
    return <div className="[&>*:nth-child(odd)]:bg-slate-800 [&>*:nth-child(even)]:bg-slate-900">
        {children}
    </div>
}