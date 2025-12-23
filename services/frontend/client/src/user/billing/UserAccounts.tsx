import { useEffect, useState } from "react"
import { useNavigate, Routes, Route } from "react-router-dom"
import { API, type ApiBillingAccount } from "../../backend"
import BillingAccountView from "./BillingAccountView"

export default function UserBilling() {
    const [accounts, setAccounts] = useState<ApiBillingAccount[]>([])
    const [loading, setLoading] = useState(false)
    const [creating, setCreating] = useState(false)
    const [provider, setProvider] = useState("")
    const [providerId, setProviderId] = useState("")
    const navigate = useNavigate()

    const fetchAccounts = async () => {
        setLoading(true)
        const result = await API.user.billing.list()
        if (result.ok) setAccounts(result.json)
        setLoading(false)
    }

    const handleCreate = async () => {
        if (!provider || !providerId) return
        setCreating(true)
        await API.user.billing.create(provider, providerId)
        setProvider("")
        setProviderId("")
        await fetchAccounts()
        setCreating(false)
    }

    const handleSuspend = async (id: string) => {
        await API.user.billing.suspend(id)
        await fetchAccounts()
    }

    const handleClose = async (id: string) => {
        await API.user.billing.close(id)
        await fetchAccounts()
    }

    useEffect(() => {
        fetchAccounts()
    }, [])

    return (
        <Routes>
            <Route
                index
                element={
                    <div className="w-full h-full flex flex-col items-center p-4">
                        {/* Add new account */}
                        <div className="w-full max-w-4xl flex flex-col md:flex-row gap-2 mb-6">
                            <input
                                type="text"
                                placeholder="Payment Provider"
                                className="border rounded p-2 flex-1 focus:ring-2 focus:ring-blue-500 focus:outline-none"
                                value={provider}
                                onChange={(e) => setProvider(e.target.value)}
                            />
                            <input
                                type="text"
                                placeholder="Provider Account ID"
                                className="border rounded p-2 flex-1 focus:ring-2 focus:ring-blue-500 focus:outline-none"
                                value={providerId}
                                onChange={(e) => setProviderId(e.target.value)}
                            />
                            <button
                                className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 disabled:opacity-50 transition"
                                disabled={creating}
                                onClick={handleCreate}
                            >
                                {creating ? "Creating..." : "Create"}
                            </button>
                        </div>

                        {/* Accounts list */}
                        {loading ? (
                            <p className="text-slate-500">Loading...</p>
                        ) : accounts.length === 0 ? (
                            <p className="text-slate-500">No billing accounts found.</p>
                        ) : (
                            <div className="w-full max-w-4xl flex flex-col gap-4">
                                {accounts.map((acct) => (
                                    <div
                                        key={acct.id}
                                        className="flex flex-col md:flex-row md:items-center md:justify-between border rounded-lg p-4 hover:shadow-md transition cursor-pointer bg-white dark:bg-slate-900 dark:border-slate-800"
                                        onClick={() => navigate(`/dashboard/billing/${acct.id}`)}
                                    >
                                        <div className="space-y-1 text-slate-800 dark:text-slate-200">
                                            <p><span className="font-semibold">ID:</span> {acct.id}</p>
                                            <p><span className="font-semibold">Status:</span> {acct.status}</p>
                                            <p><span className="font-semibold">Provider:</span> {acct.paymentProvider}</p>
                                            <p><span className="font-semibold">Provider Account ID:</span> {acct.providerAccountId}</p>
                                            <p><span className="font-semibold">Created:</span> {new Date(acct.createdAt * 1000).toLocaleString()}</p>
                                        </div>

                                        <div className="flex gap-2 mt-2 md:mt-0">
                                            {acct.status === "active" && (
                                                <button
                                                    className="bg-yellow-500 text-white px-3 py-1 rounded hover:bg-yellow-600 transition"
                                                    onClick={(e) => { e.stopPropagation(); handleSuspend(acct.id); }}
                                                >
                                                    Suspend
                                                </button>
                                            )}
                                            {acct.status !== "closed" && (
                                                <button
                                                    className="bg-red-600 text-white px-3 py-1 rounded hover:bg-red-700 transition"
                                                    onClick={(e) => { e.stopPropagation(); handleClose(acct.id); }}
                                                >
                                                    Close
                                                </button>
                                            )}
                                        </div>
                                    </div>
                                ))}
                            </div>
                        )}
                    </div>
                }
            />

            <Route path=":id" element={<BillingAccountView />} />
        </Routes>
    )
}
