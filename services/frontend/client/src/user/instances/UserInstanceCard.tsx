import type { ApiInstance } from "../../backend"
import clsx from "clsx"

type UserInstanceCardProps = {
    instance: ApiInstance
}

function StatusBadge({ status }: { status: string }) {
    const color = clsx(
        "px-2 py-1 rounded text-xs font-semibold inline-block",
        status === "running" && "bg-green-100 text-green-800",
        status === "stopped" && "bg-gray-200 text-gray-800 dark:bg-gray-700 dark:text-gray-200",
        status === "error" && "bg-red-100 text-red-800"
    )
    return <span className={color}>{status.toUpperCase()}</span>
}

export function UserInstanceCard({ instance }: UserInstanceCardProps) {
    return (
        <div className="border rounded-xl p-4 shadow-sm hover:shadow-md transition bg-white dark:bg-slate-900 dark:border-slate-800 flex flex-col gap-3">
            <div className="flex justify-between items-center">
                <h3 className="font-semibold text-lg truncate">{instance.listingTitle || instance.listingId}</h3>
                <StatusBadge status={instance.state} />
            </div>

            <div className="font-mono text-sm text-slate-700 dark:text-slate-300 space-y-1">
                <p><strong>Instance ID:</strong> {instance.id}</p>
                <p><strong>Billing Account:</strong> {instance.billingId}</p>
            </div>

            <div className="grid grid-cols-2 gap-2 text-sm text-slate-700 dark:text-slate-300">
                <p><strong>CPU:</strong> {instance.hardwareSpecification?.cpu}</p>
                <p><strong>Memory:</strong> {instance.hardwareSpecification?.memory} MB</p>
                <p><strong>Storage:</strong> {instance.hardwareSpecification?.storage} GB</p>
                <p><strong>GPU:</strong> {instance.hardwareSpecification?.gpu || "None"}</p>
            </div>
        </div>
    )
}
