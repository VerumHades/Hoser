import { useNavigate } from "react-router-dom";
import backend_constants from "../../backend_constants";

import Query from "../../components/querying/Query";
import Table from "../../components/table/Table";
import TableRow from "../../components/table/TableRow";
import { gotoListing } from "../../explore/PublicListingView";
import { useEffect, useState } from "react";
import { API, type HardwareSpecification, type Listing, type ApiBillingAccount } from "../../backend";
import SelectBox from "../../components/input/SelectBox";

export default function UserLibrary() {
    const queryBuilder = (query: string) => ({ q: query });
    const navigate = useNavigate();
    const [loadingInstanceId, setLoadingInstanceId] = useState<string | null>(null);
    const [billingAccounts, setBillingAccounts] = useState<ApiBillingAccount[] | null>(null);
    const [selectedBillingAccount, setSelectedBillingAccount] = useState<string | null>(null);
    const [showSelectBoxFor, setShowSelectBoxFor] = useState<string | null>(null);

    // Default hardware spec for demo purposes
    const defaultHardwareSpec: HardwareSpecification = {
        cpu: "4 cores",
        memory: 8192,
        storage: 100,
        gpu: undefined
    };

    // Fetch billing accounts on mount
    useEffect(() => {
        const fetchBillingAccounts = async () => {
            const response = await API.user.billing.list();
            if (response.ok) {
                if (response.json.length === 0) {
                    alert("You need a billing account to launch an instance.");
                    navigate("/user/billing");
                } else {
                    setBillingAccounts(response.json);
                    if (response.json.length === 1) {
                        setSelectedBillingAccount(response.json[0].id);
                    }
                }
            } else {
                alert("Failed to fetch billing accounts: " + JSON.stringify(response.error));
            }
        };
        fetchBillingAccounts();
    }, [navigate]);

    const launchInstance = async (listingId: string) => {
        if (!billingAccounts || billingAccounts.length === 0) {
            alert("No billing account available.");
            return;
        }

        if (!selectedBillingAccount) {
            // If multiple accounts, show SelectBox modal for this listing
            setShowSelectBoxFor(listingId);
            return;
        }

        try {
            setLoadingInstanceId(listingId);
            const response = await API.user.instances.launch(listingId, selectedBillingAccount, defaultHardwareSpec);
            if (response.ok) {
                alert("Instance launched successfully!");
            } else {
                alert("Failed to launch instance: " + JSON.stringify(response.error));
            }
        } finally {
            setLoadingInstanceId(null);
        }
    };

    const handleBillingAccountSelected = (billingAccountId: string) => {
        setSelectedBillingAccount(billingAccountId);
        if (showSelectBoxFor) {
            launchInstance(showSelectBoxFor);
            setShowSelectBoxFor(null);
        }
    };

    return (
        <div className="w-full h-full flex flex-col items-center">
            <Query
                bodyBuilder={(entries?: Listing[]) => (
                    <Table>
                        {entries && entries.map((entry) => (
                            <TableRow key={entry.id} className="flex flex-col space-y-2">
                                <div
                                    className="cursor-pointer"
                                    onClick={() => gotoListing(navigate, entry)}
                                >
                                    <h3 className="font-semibold text-lg">{entry.title}</h3>
                                    <p className="text-sm text-gray-600">{entry.description}</p>
                                </div>
                                <button
                                    className="px-3 py-1 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50"
                                    disabled={loadingInstanceId === entry.id || !billingAccounts}
                                    onClick={() => launchInstance(entry.id)}
                                >
                                    {loadingInstanceId === entry.id ? "Launching..." : "Launch Instance"}
                                </button>

                                {/* Show SelectBox if this listing needs a billing account selection */}
                                {showSelectBoxFor === entry.id && billingAccounts && billingAccounts.length > 1 && (
                                    <SelectBox
                                        options={Object.fromEntries(
                                            billingAccounts.map(b => [b.id, { label: b.id, description: b.status }])
                                        )}
                                        onSelected={(id) => handleBillingAccountSelected(id)}
                                    />
                                )}
                            </TableRow>
                        ))}
                    </Table>
                )}
                queryBuilder={queryBuilder}
                endpoint={`${backend_constants.address}/user/library`}
            />
        </div>
    );
}
