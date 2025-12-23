import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import {
    API,
    type HardwareSpecification,
    type Listing,
    type ApiBillingAccount,
} from "../../backend";
import SelectBox from "../../components/input/SelectBox";
import HardwareSettings from "../developer/hardware/HardwareSettings";
import CreateBillingAccount from "../billing/CreateBillingAccount";

export default function CreateInstance() {
    const { listingId } = useParams<{ listingId: string }>();
    const navigate = useNavigate();

    const [listing, setListing] = useState<Listing | null>(null);
    const [billingAccounts, setBillingAccounts] = useState<ApiBillingAccount[] | null>(null);
    const [selectedBillingAccount, setSelectedBillingAccount] = useState<string | null>(null);
    const [loading, setLoading] = useState(false);

    const [hardwareSpecification, setHardwareSpecification] =
        useState<HardwareSpecification>({
            cpu: 4,
            ramBytes: 8 * 1024 * 1024 * 1024,
            diskBytes: 100 * 1024 * 1024 * 1024,
        });

    useEffect(() => {
        if (!listingId) return;

        API.listing.get(listingId).then((response) => {
            if (response.ok) {
                setListing(response.json);
            } else {
                alert("Failed to load listing.");
                navigate("/dashboard/library");
            }
        });
    }, [listingId, navigate]);

    useEffect(() => {
        API.user.billing.list().then((response) => {
            if (response.ok) {
                setBillingAccounts(response.json);
                if (response.json.length === 1) {
                    setSelectedBillingAccount(response.json[0].id);
                }
            } else {
                alert("Failed to fetch billing accounts.");
            }
        });
    }, []);

    const calculateTotalPrice = (): number => {
        if (!listing) return 0;

        const basePrice = listing.price?.amount ?? 0;
        const ramPrice = (hardwareSpecification.ramBytes / 1024 ** 3) * 5;
        const diskPrice = (hardwareSpecification.diskBytes / 1024 ** 3) * 0.1;

        return basePrice + ramPrice + diskPrice;
    };

    const launchInstance = async () => {
        if (!listingId || !selectedBillingAccount) {
            alert("Please select a billing account.");
            return;
        }

        try {
            setLoading(true);

            const response = await API.user.instances.launch(
                listingId,
                selectedBillingAccount,
                hardwareSpecification
            );

            if (response.ok) {
                navigate("/dashboard/library");
            } else {
                alert("Failed to launch instance.");
            }
        } finally {
            setLoading(false);
        }
    };

    if (!listing || !billingAccounts) {
        return <p>Loading...</p>;
    }

    const billingOptions: Record<string, { label: string }> = {};
    billingAccounts.forEach((account) => {
        if(account.status != "active") return;
        billingOptions[account.id] = { label: account.id };
    });

    return (
        <div className="max-w-3xl mx-auto p-6 space-y-6">
            <h1 className="text-2xl font-bold">{listing.title}</h1>
            <p className="text-gray-700">{listing.description}</p>

            <div className="space-y-2">
                <label className="font-semibold">Select Billing Account</label>
                <SelectBox
                    options={billingOptions}
                    onEmpty={() => 
                        <CreateBillingAccount/>
                        }
                    defaultValue={selectedBillingAccount ?? undefined}
                    onSelected={setSelectedBillingAccount}
                />
            </div>

            <HardwareSettings
                initialSpec={hardwareSpecification}
                onChange={setHardwareSpecification}
            />

            <div className="space-y-1">
                <h2 className="font-semibold">Price Summary</h2>
                <p>Base Price: ${listing.price?.amount ?? 0}</p>
                <p>Hardware Price: ${(calculateTotalPrice() - (listing.price?.amount ?? 0)).toFixed(2)}</p>
                <p className="font-bold">Total: ${calculateTotalPrice().toFixed(2)}</p>
            </div>

            <button
                className="px-4 py-2 bg-blue-600 text-white rounded disabled:opacity-50"
                disabled={loading}
                onClick={launchInstance}
            >
                {loading ? "Launching..." : "Launch Instance"}
            </button>
        </div>
    );
}
