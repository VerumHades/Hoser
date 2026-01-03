import type { ApiInstance } from "../../backend"
import { formatBytes } from "../../components/common"
import { ContractBadge } from "./InstanceContractBadge"
import { StatusBadge } from "./InstanceStatusBadge"

type UserInstanceRowProps = {
    instance: ApiInstance
}

export function UserInstanceRow({ instance }: UserInstanceRowProps) {
    const hardwareSpecification = contract.hardwareSpecification

    return (
        <tr className="hover:bg-slate-50 dark:hover:bg-slate-800 transition cursor-pointer">
            <td className="px-4 py-2 font-mono text-sm text-slate-700 dark:text-slate-300">
                {contract.id}
            </td>

            <td className="px-4 py-2 font-medium text-slate-900 dark:text-slate-100">
                {contract.listingId}
            </td>

            <td className="px-4 py-2 space-x-2">
                <StatusBadge state={contract.state} />
                <ContractBadge contractState={contract.contractState} />
            </td>

            <td className="px-4 py-2 text-sm text-slate-700 dark:text-slate-300 space-x-2">
                <span>CPU: {hardwareSpecification?.cpu ?? "—"}</span>
                <span>Mem: {formatBytes(hardwareSpecification?.ramBytes)}</span>
                <span>Disk: {formatBytes(hardwareSpecification?.diskBytes)}</span>
            </td>

            <td className="px-4 py-2 font-mono text-sm text-slate-700 dark:text-slate-300">
                {contract.billingId}
            </td>
        </tr>
    )
}
