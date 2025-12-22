import { useParams } from "react-router";
import { useEffect, useState } from "react";
import { API, type ApiBillingAccount, type ApiPayment } from "../../backend";

export default function BillingAccountView() {
    const { id } = useParams<{ id: string }>();
    const [account, setAccount] = useState<ApiBillingAccount | null>(null);
    const [payments, setPayments] = useState<ApiPayment[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchAccountAndPayments = async () => {
            if (!id) return;

            const accRes = await API.user.billing.get(id);
            if (accRes.ok) setAccount(accRes.json);

            const paymentsRes = await API.user.billing.payments.list(id);
            if (paymentsRes.ok) setPayments(paymentsRes.json);

            setLoading(false);
        };
        fetchAccountAndPayments();
    }, [id]);

    if (loading) return <p className="text-center py-10 text-gray-500">Loading...</p>;
    if (!account) return <p className="text-center py-10 text-red-500">Billing account not found.</p>;

    return (
        <div className="w-full max-w-5xl mx-auto p-6 flex flex-col gap-8">
            {/* Billing Account Info */}
            <section className="bg-white dark:bg-gray-900 border border-slate-200 dark:border-gray-800 rounded-xl p-6 shadow-sm">
                <h2 className="text-2xl font-semibold mb-4">Billing Account {account.id}</h2>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <InfoRow label="Status" value={account.status} />
                    <InfoRow label="Provider" value={account.paymentProvider} />
                    <InfoRow label="Provider Account ID" value={account.providerAccountId} />
                    <InfoRow label="Created" value={new Date(account.createdAt * 1000).toLocaleString()} />
                </div>
            </section>

            {/* Payments List */}
            <section className="bg-white dark:bg-gray-900 border border-slate-200 dark:border-gray-800 rounded-xl p-6 shadow-sm flex flex-col gap-4">
                <h3 className="text-xl font-semibold mb-4">Payments</h3>
                {payments.length === 0 ? (
                    <p className="text-gray-500">No payments found.</p>
                ) : (
                    <ul className="flex flex-col gap-4">
                        {payments.map((payment) => (
                            <li key={payment.id} className="border border-slate-200 dark:border-gray-800 rounded-lg p-4 hover:shadow-md transition-shadow">
                                <InfoRow label="ID" value={payment.id} />
                                <InfoRow label="Amount" value={`${payment.amount} ${payment.currency}`} />
                                <InfoRow label="Status" value={payment.status} />
                                <InfoRow label="Kind" value={payment.kind} />
                                {payment.metadata && (
                                    <div className="mt-2">
                                        <span className="font-semibold">Metadata:</span>
                                        <pre className="bg-gray-100 dark:bg-gray-800 p-2 rounded mt-1 overflow-x-auto">
                                            {JSON.stringify(payment.metadata, null, 2)}
                                        </pre>
                                    </div>
                                )}
                            </li>
                        ))}
                    </ul>
                )}
            </section>
        </div>
    );
}

function InfoRow({ label, value }: { label: string; value: string | number }) {
    return (
        <div className="flex justify-between">
            <span className="font-medium text-slate-700 dark:text-slate-400">{label}:</span>
            <span className="text-slate-900 dark:text-slate-200">{value}</span>
        </div>
    );
}
