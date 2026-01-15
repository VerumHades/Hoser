import { BookmarkPlus, BookmarkCheck, LogIn, CheckCircle, Loader2, ShoppingCart } from "lucide-react";
import { TitleAndDescription } from "../../components/prefabs/TitleAndDescription";
import type { Listing } from "../../backend/repositories/listing";

interface HeaderProps {
    listing: Listing;
    isAuthenticated: boolean;
    isSaved: boolean;
    isOwned: boolean;
    isProcessing: boolean;
    isActionLoading: { library: boolean; purchase: boolean };
    onSave: () => void;
    onPurchase: () => void;
    onLogin: () => void;
}

export function ListingHeader({ 
    listing, isAuthenticated, isSaved, isOwned, isProcessing, isActionLoading, onSave, onPurchase, onLogin 
}: HeaderProps) {
    const formatCurrency = (amt: number) => 
        new Intl.NumberFormat("de-DE", { style: "currency", currency: "EUR" }).format(amt / 100);

    return (
        <section className="flex flex-col md:flex-row justify-between gap-6 border-b border-slate-200 dark:border-slate-800 pb-8">
            <TitleAndDescription
                title={listing.title ?? "Untitled"}
                description={listing.description ?? "No description available."}
                titleClassname="text-4xl font-extrabold tracking-tight"
                descriptionClassname="max-w-2xl text-slate-600 dark:text-slate-400"
            />

            <div className="flex flex-col items-start md:items-end gap-4 min-w-[280px]">
                <div className="text-4xl font-black text-slate-900 dark:text-white">{formatCurrency(listing.price)}</div>
                <div className="flex items-center gap-2 w-full md:w-auto">
                    <button 
                        onClick={onSave} 
                        disabled={isActionLoading.library} 
                        className={`flex items-center gap-2 px-4 py-2.5 rounded-xl border border-slate-200 dark:border-slate-800 hover:bg-slate-50 dark:hover:bg-slate-900 transition-all ${isSaved ? "text-blue-600 dark:text-blue-600 bg-blue-50 dark:bg-blue-900/20" : "dark:text-slate-50"}`}
                    >
                        {isSaved ? <BookmarkCheck className="w-4 h-4" /> : <BookmarkPlus className="w-4 h-4" />}
                        <span className="text-sm font-semibold">{isSaved ? "Saved" : "Save"}</span>
                    </button>

                    {!isAuthenticated ? (
                        <button onClick={onLogin} className="flex-1 md:flex-none flex items-center gap-2 px-6 py-2.5 bg-slate-900 dark:bg-white text-white dark:text-slate-900 rounded-xl font-bold hover:opacity-90">
                            <LogIn className="w-4 h-4" /> Sign in to Buy
                        </button>
                    ) : (
                        <button 
                            onClick={onPurchase} 
                            disabled={isActionLoading.purchase || isOwned || isProcessing}
                            className={`flex-1 md:flex-none flex items-center justify-center gap-2 px-6 py-2.5 rounded-xl font-bold transition-all ${isOwned ? "bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400" : isProcessing ? "bg-amber-100 text-amber-700" : "bg-blue-600 text-white hover:bg-blue-700 shadow-lg shadow-blue-500/20"}`}
                        >
                            {isOwned ? <><CheckCircle className="w-4 h-4" /> Owned</> : 
                             isProcessing ? <><Loader2 className="w-4 h-4 animate-spin" /> Processing</> : 
                             <><ShoppingCart className="w-4 h-4" /> {isActionLoading.purchase ? "Authorizing..." : "Purchase Now"}</>}
                        </button>
                    )}
                </div>
            </div>
        </section>
    );
}