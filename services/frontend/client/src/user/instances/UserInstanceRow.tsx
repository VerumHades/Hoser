import type { ApiInstance } from "../../backend"
import { formatBytes } from "../../components/common"
import { ContractBadge } from "./InstanceContractBadge"
import { StatusBadge } from "./InstanceStatusBadge"

type UserInstanceRowProps = {
    instance: ApiInstance
}

export function UserInstanceRow({ instance }: UserInstanceRowProps) {
    const hardwareSpecification = instance.hardwareSpecification

    return (
        <tr className="hover:bg-slate-50 dark:hover:bg-slate-800 transition cursor-pointer">
            <td className="px-4 py-2 font-mono text-sm text-slate-700 dark:text-slate-300">
                {instance.id}
            </td>

            <td className="px-4 py-2 font-medium text-slate-900 dark:text-slate-100">
                {instance.listingId}
            </td>

            <td className="px-4 py-2 space-x-2">
                <StatusBadge state={instance.state} />
                <ContractBadge contractState={instance.contractState} />
            </td>

            <td className="px-4 py-2 text-sm text-slate-700 dark:text-slate-300 space-x-2">
                <span>CPU: {hardwareSpecification?.cpu ?? "—"}</span>
                <span>Mem: {formatBytes(hardwareSpecification?.ramBytes)}</span>
                <span>Disk: {formatBytes(hardwareSpecification?.diskBytes)}</span>
            </td>

            <td className="px-4 py-2 font-mono text-sm text-slate-700 dark:text-slate-300">
                {instance.billingId}
            </td>
        </tr>
    )
}
