import { useEffect, useState } from "react"
import { formatBytes } from "../../components/common"
import { ConsolePrompt } from "../../components/view/ConsolePrompt"

type UserInstanceCardProps = {
    instance: ApiInstance
}

function HardwareRow({ label, value }: { label: string; value: string }) {
    return (
        <div className="flex justify-between text-sm">
            <span className="text-slate-500 dark:text-slate-400">{label}</span>
            <span className="font-mono text-slate-800 dark:text-slate-200">{value}</span>
        </div>
    )
}

/**
 * Provides a colored banner for instance or contract states.
 */
function StateBanner({ label, state }: { label: string; state: string }) {
    const stateColors: Record<string, string> = {
        BUILDING: "bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-300",
        RUNNING: "bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-300",
        STOPPED: "bg-slate-100 text-slate-800 dark:bg-slate-800/50 dark:text-slate-300",
        UPDATING_HARDWARE: "bg-amber-100 text-amber-800 dark:bg-amber-900/20 dark:text-amber-300",
        ERROR: "bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-300",
        UNKNOWN: "bg-gray-100 text-gray-800 dark:bg-gray-800/50 dark:text-gray-300",
        INACTIVE: "bg-gray-100 text-gray-800 dark:bg-gray-800/50 dark:text-gray-300",
        ACTIVE: "bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-300",
    }

    const colorClass = stateColors[state] ?? stateColors.UNKNOWN

    return (
        <div className={`rounded-lg px-3 py-1 text-xs font-semibold ${colorClass}`}>
            {label}: {state}
        </div>
    )
}

export function UserInstanceCard({ instance }: UserInstanceCardProps) {
    const hardwareSpecification = contract.hardwareSpecification
    const desiredHardwareSpecification = contract.desiredHardwareSpecification
    const hasHardwareDrift =
        desiredHardwareSpecification !== undefined &&
        JSON.stringify(desiredHardwareSpecification) !== JSON.stringify(hardwareSpecification)

    const [deploymentState, setDeploymentState] = useState<DeployedInstanceState | null>(null)
    const [loadingState, setLoadingState] = useState(true)

    useEffect(() => {
        let isMounted = true

        async function fetchState() {
            try {
                setLoadingState(true)
                const response = await API.user.instances.state.get(contract.id)
                if (isMounted) setDeploymentState(response.json)
            } catch (err) {
                console.error("Failed to fetch deployment state:", err)
            } finally {
                if (isMounted) setLoadingState(false)
            }
        }

        fetchState()
        const interval = setInterval(fetchState, 5000)
        return () => {
            isMounted = false
            clearInterval(interval)
        }
    }, [contract.id])

    return (
        <div className="rounded-2xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 p-5 shadow-sm hover:shadow-lg transition flex flex-col gap-4">
            
            {/* Header: Status and Contract */}
            <div className="flex flex-col sm:flex-row justify-between gap-2">
                <div className="flex gap-2">
                    <StateBanner label="Instance" state={contract.state} />
                    <StateBanner label="Contract" state={contract.contractState} />
                </div>
                {hasHardwareDrift && (
                    <div className="rounded-lg bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-800 px-3 py-1 text-xs font-medium text-amber-800 dark:text-amber-300 text-center">
                        Hardware update pending
                    </div>
                )}
            </div>

            {/* Deployment State & Errors */}
            {deploymentState && (
                <div className="flex flex-col gap-2 p-3 rounded-lg border border-slate-200 dark:border-slate-800 bg-slate-50 dark:bg-slate-800/30">
                    {deploymentState.stateMessage && (
                        <div className="text-sm font-medium text-slate-700 dark:text-slate-300">
                            {deploymentState.stateMessage}
                        </div>
                    )}
                    {deploymentState.state === "ErrorState" && deploymentState.error && (
                        <ConsolePrompt errorText={deploymentState.error} />
                    )}
                </div>
            )}

            {/* Hardware specifications */}
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 rounded-xl bg-slate-50 dark:bg-slate-800/50 p-4">
                <HardwareRow
                    label="CPU"
                    value={hardwareSpecification?.cpu !== undefined ? `${hardwareSpecification.cpu} cores` : "—"}
                />
                <HardwareRow label="Memory" value={formatBytes(hardwareSpecification?.ramBytes)} />
                <HardwareRow label="Disk" value={formatBytes(hardwareSpecification?.diskBytes)} />
            </div>

            {/* Billing info */}
            <div className="flex justify-between text-xs text-slate-500 dark:text-slate-400">
                <span>Billing Account</span>
                <span className="font-mono">{contract.billingId}</span>
            </div>
        </div>
    )
}
