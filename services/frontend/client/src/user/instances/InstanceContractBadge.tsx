import type { ApiInstance } from "../../backend"
import clsx from "clsx"

type ContractBadgeProps = {
    contractState: ApiInstance["contractState"]
}

export function ContractBadge({ contractState }: ContractBadgeProps) {
    const className = clsx(
        "inline-flex items-center px-3 py-1 rounded-full text-xs font-semibold tracking-wide",
        contractState === "ACTIVE" && "bg-emerald-100 text-emerald-800 dark:bg-emerald-900/40 dark:text-emerald-300",
        contractState === "INACTIVE" && "bg-gray-200 text-gray-700 dark:bg-gray-700 dark:text-gray-300",
        contractState === "UNKNOWN" && "bg-yellow-100 text-yellow-800 dark:bg-yellow-900/40 dark:text-yellow-300"
    )

    return <span className={className}>{contractState.replaceAll("_", " ")}</span>
}
