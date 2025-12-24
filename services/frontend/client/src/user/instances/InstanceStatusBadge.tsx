import type { ApiInstance } from "../../backend"
import clsx from "clsx"

type StatusBadgeProps = {
    state: ApiInstance["state"]
}

export function StatusBadge({ state }: StatusBadgeProps) {
    const className = clsx(
        "inline-flex items-center px-3 py-1 rounded-full text-xs font-semibold tracking-wide",
        state === "RUNNING" && "bg-emerald-100 text-emerald-800 dark:bg-emerald-900/40 dark:text-emerald-300",
        state === "STOPPED" && "bg-slate-200 text-slate-800 dark:bg-slate-700 dark:text-slate-200",
        state === "UPDATING_HARDWARE" && "bg-blue-100 text-blue-800 dark:bg-blue-900/40 dark:text-blue-300",
        state === "ERROR" && "bg-red-100 text-red-800 dark:bg-red-900/40 dark:text-red-300",
        state === "BUILDING" && "bg-cyan-100 text-cyan-800 dark:bg-cyan-900/40 dark:text-cyan-300",
        state === "UNKNOWN" && "bg-gray-200 text-gray-700 dark:bg-gray-700 dark:text-gray-300"
    )

    return <span className={className}>{state.replaceAll("_", " ")}</span>
}
