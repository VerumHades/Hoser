import { useEffect, useState } from "react";
import { API, type ApiInstance } from "../../backend";
import Table from "../../components/table/Table";
import TableRow from "../../components/table/TableRow";

/**
 * Page to list all user instances.
 * Displays instance ID, listing ID, billing account, hardware, and state.
 */
export default function UserInstances() {
    const [instances, setInstances] = useState<ApiInstance[] | null>(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchInstances = async () => {
            setLoading(true);
            const response = await API.user.instances.list();
            if (response.ok) {
                setInstances(response.json);
            } else {
                alert("Failed to load instances: " + JSON.stringify(response.error));
            }
            setLoading(false);
        };
        fetchInstances();
    }, []);

    if (loading) return <p>Loading instances...</p>;
    if (!instances || instances.length === 0) return <p>No instances found.</p>;

    return (
        <div className="max-w-5xl mx-auto p-6 space-y-4">
            <Table>
                {instances.map((instance) => (
                    <TableRow key={instance.id} className="flex flex-col md:flex-row justify-between items-start md:items-center px-4 py-2 hover:bg-gray-50 rounded">
                        <div className="flex-1 space-y-1">
                            <p><strong>Instance ID:</strong> {instance.id}</p>
                            <p><strong>Listing ID:</strong> {instance.listingId}</p>
                            <p><strong>Billing Account:</strong> {instance.billingId}</p>
                            <p><strong>State:</strong> {instance.state}</p>
                        </div>
                        <div className="mt-2 md:mt-0 flex flex-col space-y-1">
                            <p><strong>Hardware:</strong></p>
                            <p>CPU: {instance.hardwareSpecification?.cpu}</p>
                            <p>Memory: {instance.hardwareSpecification?.memory} MB</p>
                            <p>Storage: {instance.hardwareSpecification?.storage} GB</p>
                            <p>GPU: {instance.hardwareSpecification?.gpu || "None"}</p>
                        </div>
                    </TableRow>
                ))}
            </Table>
        </div>
    );
}
