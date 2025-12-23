import { useEffect, useMemo, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import {
    API,
    type ApiBillingAccount,
    type ApiHardwareRatesResponse,
    type HardwareSpecification,
    type Listing,
} from "../../backend";
import SelectBox from "../../components/input/SelectBox";
import HardwareSettings from "../developer/hardware/HardwareSettings";
import CreateBillingAccount from "../billing/CreateBillingAccount";

export default function CreateInstance() {
    const { listingId } = useParams<{ listingId: string }>();
    const navigate = useNavigate();

    const [listing, setListing] = useState<Listing | null>(null);
    const [billingAccounts, setBillingAccounts] = useState<ApiBillingAccount[] | null>(null);
    const [selectedBillingAccountId, setSelectedBillingAccountId] = useState<string | null>(null);
    const [isLaunching, setIsLaunching] = useState<boolean>(false);
    const [rentalDurationMonths, setRentalDurationMonths] = useState<number>(1);
    const [hardwareRates, setHardwareRates] = useState<ApiHardwareRatesResponse | null>(null);

    const [hardwareSpecification, setHardwareSpecification] =
        useState<HardwareSpecification>({
            cpu: 4,
            ramBytes: 8 * 1024 * 1024 * 1024,
            diskBytes: 100 * 1024 * 1024 * 1024,
        });

    useEffect(() => {
        API.hardware.getRates("USD").then((response) => {
            if (!response.ok) {
                return;
            }

            setHardwareRates(response.json);
        });
    }, []);


    useEffect(() => {
        if (!listingId) {
            navigate("/dashboard/library");
            return;
        }

        API.listing.get(listingId).then((response) => {
            if (!response.ok) {
                navigate("/dashboard/library");
                return;
            }

            setListing(response.json);
        });
    }, [listingId, navigate]);

    useEffect(() => {
        API.user.billing.list().then((response) => {
            if (!response.ok) {
                return;
            }

            setBillingAccounts(response.json);

            if (response.json.length === 1) {
                setSelectedBillingAccountId(response.json[0].id);
            }
        });
    }, []);

    /**
     * Computes the full monthly price based on listing base price and hardware configuration.
     */
    const { monthlyPrice, totalPrice } = useMemo(() => {
        if (!listing || !hardwareRates) {
            return { monthlyPrice: 0, totalPrice: 0 };
        }

        const hoursPerMonth = 24 * 30;

        const cpuHourlyCost =
            hardwareRates.cpuCost.amount *
            hardwareSpecification.cpu *
            hoursPerMonth;

        const ramHourlyCost =
            hardwareRates.ramCost.amount *
            hardwareSpecification.ramBytes *
            hoursPerMonth;

        const diskHourlyCost =
            hardwareRates.diskCost.amount *
            hardwareSpecification.diskBytes *
            hoursPerMonth;

        const basePrice = listing.price?.amount ?? 0;

        const monthly =
            basePrice +
            cpuHourlyCost +
            ramHourlyCost +
            diskHourlyCost;

        return {
            monthlyPrice: monthly,
            totalPrice: monthly * rentalDurationMonths,
        };
    }, [
        listing,
        hardwareRates,
        hardwareSpecification,
        rentalDurationMonths,
    ]);

    /**
     * Launches a new instance using the selected configuration.
     */
    const launchInstance = async (): Promise<void> => {
        if (!listingId || !selectedBillingAccountId) {
            return;
        }

        setIsLaunching(true);

        try {
            const response = await API.user.instances.launch(
                listingId,
                selectedBillingAccountId,
                hardwareSpecification
            );

            if (response.ok) {
                navigate("/dashboard/library");
            }
        } finally {
            setIsLaunching(false);
        }
    };

    if (!listing || !billingAccounts || !hardwareRates) {
        return (
            <div className="min-h-screen flex items-center justify-center text-slate-500">
                Loading…
            </div>
        );
    }


    const billingAccountOptions: Record<string, { label: string }> = {};
    billingAccounts.forEach((account) => {
        if (account.status !== "active") {
            return;
        }

        billingAccountOptions[account.id] = {
            label: account.id,
        };
    });

    return (
        <div className="min-h-screen w-full bg-slate-50 dark:bg-slate-950 px-6 py-16">
            <div className="max-w-5xl mx-auto flex flex-col gap-12">
                <header className="flex flex-col gap-4">
                    <h1 className="text-3xl md:text-4xl font-bold tracking-tight text-slate-900 dark:text-slate-100">
                        Launch instance
                    </h1>

                    <p className="max-w-2xl text-lg text-slate-600 dark:text-slate-400">
                        Configure hardware and billing for{" "}
                        <span className="font-medium text-slate-900 dark:text-slate-200">
                            {listing.title}
                        </span>
                    </p>
                </header>

                <section className="grid grid-cols-1 md:grid-cols-3 gap-8">
                    <div className="md:col-span-2 flex flex-col gap-8">
                        <div className="p-6 rounded-xl bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800">
                            <h2 className="text-lg font-semibold text-slate-900 dark:text-slate-100 mb-4">
                                Billing account
                            </h2>

                            <SelectBox
                                options={billingAccountOptions}
                                defaultValue={selectedBillingAccountId ?? undefined}
                                onSelected={setSelectedBillingAccountId}
                                onEmpty={() => <CreateBillingAccount />}
                            />
                        </div>

                        <div className="p-6 rounded-xl bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800">
                            <h2 className="text-lg font-semibold text-slate-900 dark:text-slate-100 mb-2">
                                Rental duration
                            </h2>

                            <p className="text-sm text-slate-600 dark:text-slate-400 mb-4">
                                Choose how long you want to rent this hardware.
                            </p>

                            <select
                                value={rentalDurationMonths}
                                onChange={(event) => setRentalDurationMonths(Number(event.target.value))}
                                className="w-full px-4 py-3 rounded-lg bg-slate-50 dark:bg-slate-800 border border-slate-300 dark:border-slate-700 text-slate-900 dark:text-slate-100 focus:outline-none focus:ring-2 focus:ring-slate-900 dark:focus:ring-slate-100"
                            >
                                <option value={1}>1 month</option>
                                <option value={3}>3 months</option>
                                <option value={6}>6 months</option>
                                <option value={12}>12 months</option>
                            </select>
                        </div>

                        <div className="p-6 rounded-xl bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800">
                            <h2 className="text-lg font-semibold text-slate-900 dark:text-slate-100 mb-4">
                                Hardware configuration
                            </h2>

                            <HardwareSettings
                                initialSpec={hardwareSpecification}
                                onChange={setHardwareSpecification}
                            />
                        </div>
                    </div>

                    <aside className="flex flex-col gap-6">
                        <div className="p-6 rounded-xl bg-slate-100 dark:bg-slate-900 border border-slate-200 dark:border-slate-800">
                            <h2 className="text-lg font-semibold text-slate-900 dark:text-slate-100 mb-4">
                                Price summary
                            </h2>

                            <div className="flex flex-col gap-2 text-sm text-slate-600 dark:text-slate-400">
                                <div className="flex justify-between">
                                    <span>Monthly price</span>
                                    <span>
                                        {hardwareRates.cpuCost.currencyCode} {monthlyPrice.toFixed(2)}
                                    </span>

                                </div>

                                <div className="flex justify-between">
                                    <span>Duration</span>
                                    <span>{rentalDurationMonths} month{rentalDurationMonths > 1 ? "s" : ""}</span>
                                </div>

                                <div className="flex justify-between pt-2 border-t border-slate-300 dark:border-slate-700 font-medium text-slate-900 dark:text-slate-100">
                                    <span>Total</span>
                                    <span>
                                        {hardwareRates.cpuCost.currencyCode} {monthlyPrice.toFixed(2)}
                                    </span>

                                </div>
                            </div>
                        </div>

                        <button
                            onClick={launchInstance}
                            disabled={isLaunching || !selectedBillingAccountId}
                            className="w-full px-6 py-3 rounded-lg bg-slate-900 text-white dark:bg-slate-100 dark:text-slate-900 font-medium hover:opacity-90 transition disabled:opacity-50"
                        >
                            {isLaunching ? "Launching…" : "Launch instance"}
                        </button>
                    </aside>
                </section>
            </div>
        </div>
    );
}
