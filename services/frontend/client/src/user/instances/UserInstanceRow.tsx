import type { ApiInstance } from "../../backend"
import clsx from "clsx"

type UserInstanceRowProps = {
    instance: ApiInstance
}

function StatusBadge({ status }: { status: string }) {
    const color = clsx(
        "px-2 py-1 rounded text-xs font-semibold",
        status === "running" && "bg-green-100 text-green-800",
        status === "stopped" && "bg-gray-200 text-gray-800 dark:bg-gray-700 dark:text-gray-200",
        status === "error" && "bg-red-100 text-red-800"
    )
    return <span className={color}>{status.toUpperCase()}</span>
}

export function UserInstanceRow({ instance }: UserInstanceRowProps) {
    return (
        <tr className="hover:bg-gray-50 dark:hover:bg-slate-800 transition cursor-pointer w-full">
            <td className="px-4 py-2 font-mono text-sm">{instance.id}</td>
            <td className="px-4 py-2 font-medium">{instance.listingTitle || instance.listingId}</td>
            <td className="px-4 py-2">
                <StatusBadge status={instance.state} />
            </td>
            <td className="px-4 py-2 text-sm">
                CPU: {instance.hardwareSpecification?.cpu}, Mem: {instance.hardwareSpecification?.ramBytes}, Storage: {instance.hardwareSpecification?.diskBytes} {instance.hardwareSpecification?.gpu ? `, GPU: ${instance.hardwareSpecification.gpu}` : ""}
            </td>
            <td className="px-4 py-2 font-mono text-sm">{instance.billingId}</td>
        </tr>
    )
}
