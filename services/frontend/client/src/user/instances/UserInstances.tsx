import { useEffect, useState, useCallback } from "react"
import { API, type ApiInstance } from "../../backend"
import { CollectionViewContainer } from "../../components/view/CollectionViewContainer"
import TableRow from "../../components/table/TableRow"
import { UserInstanceRow } from "./UserInstanceRow"
import { UserInstanceCard } from "./UserInstanceCard"

type InstanceViewMode = "table" | "cards"

export default function UserInstances() {
    const [instances, setInstances] = useState<ApiInstance[] | null>(null)
    const [loading, setLoading] = useState(true)
    const [viewMode, setViewMode] = useState<InstanceViewMode>("cards")

    useEffect(() => {
        const fetchInstances = async () => {
            setLoading(true)
            const response = await API.user.instances.list()
            if (response.ok) {
                setInstances(response.json)
            } else {
                alert("Failed to load instances: " + JSON.stringify(response.error))
            }
            setLoading(false)
        }
        fetchInstances()
    }, [])

    const renderInstances = useCallback(
        (items: ApiInstance[] | undefined) => (
            <CollectionViewContainer
                items={items}
                viewMode={viewMode}
                onViewModeChange={setViewMode}
                renderTableRow={(instance) => (
                    <UserInstanceRow instance={instance}></UserInstanceRow>
                )}
                renderCard={(instance) => (
                    <UserInstanceCard key={instance.id} instance={instance} />
                )}
                emptyState={<p className="text-sm text-slate-500">No instances found.</p>}
            />
        ),
        [viewMode]
    )

    if (loading) return <p className="text-center mt-6">Loading instances...</p>

    return (
        <div className="flex-1 mx-auto p-6">
            {renderInstances(instances ?? [])}
        </div>
    )
}
