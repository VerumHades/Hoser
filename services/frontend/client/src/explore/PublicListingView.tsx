import { useEffect, useState, type JSX } from "react";
import { useParams, type NavigateFunction } from "react-router-dom";
import { TitleAndDescription } from "../components/prefabs/TitleAndDescription";
import { BookmarkPlus, BookmarkCheck, ShoppingCart, CheckCircle, Loader2 } from "lucide-react";
import { ListingAPI, type Listing } from "../backend/repositories/listing";
import { UserAPI } from "../backend/repositories/user";
import { ListingScreenshot } from "../ListingScreenshot";

/**
 * Navigates the user to the detailed view of a specific listing.
 */
export function gotoListing(navigate: NavigateFunction, listing: Listing | undefined) {
    if (!listing) return;
    navigate("/listing/" + listing.id);
}

export function PublicListingView(): JSX.Element {
    const { id: listingIdentifier } = useParams<{ id: string }>();

    const [listing, setListing] = useState<Listing | null>(null);
    const [isLoading, setIsLoading] = useState<boolean>(true);
    const [hasError, setHasError] = useState<boolean>(false);
    
    const [isSavedInLibrary, setIsSavedInLibrary] = useState<boolean>(false);
    const [isOwnedByUser, setIsOwnedByUser] = useState<boolean>(false);
    const [isPurchaseProcessing, setIsPurchaseProcessing] = useState<boolean>(false);
    
    const [isLibraryActionLoading, setIsLibraryActionLoading] = useState<boolean>(false);
    const [isPurchaseActionLoading, setIsPurchaseActionLoading] = useState<boolean>(false);

    // Initial Data Load
    useEffect(() => {
        async function initializeListingView() {
            if (!listingIdentifier) {
                setHasError(true);
                setIsLoading(false);
                return;
            }

            try {
                const [listingResult, libraryStatus, ownershipStatus, processingStatus] = await Promise.all([
                    ListingAPI.get(listingIdentifier),
                    UserAPI.library.check(listingIdentifier),
                    UserAPI.isOwner(listingIdentifier),
                    UserAPI.isPurchaseProcessingOwner(listingIdentifier)
                ]);

                setListing(listingResult);
                setIsSavedInLibrary(libraryStatus);
                setIsOwnedByUser(ownershipStatus);
                setIsPurchaseProcessing(processingStatus);
            } catch (error) {
                setHasError(true);
            } finally {
                setIsLoading(false);
            }
        }

        initializeListingView();
    }, [listingIdentifier]);

    // Polling Logic: Runs only when isPurchaseProcessing is true
    useEffect(() => {
        let pollInterval: NodeJS.Timeout | null = null;

        if (isPurchaseProcessing && !isOwnedByUser && listingIdentifier) {
            pollInterval = setInterval(async () => {
                try {
                    const [ownershipStatus, processingStatus] = await Promise.all([
                        UserAPI.isOwner(listingIdentifier),
                        UserAPI.isPurchaseProcessingOwner(listingIdentifier)
                    ]);

                    if (ownershipStatus) {
                        setIsOwnedByUser(true);
                        setIsPurchaseProcessing(false);
                    } else if (!processingStatus) {
                        // Process finished but ownership not granted (possible failure)
                        setIsPurchaseProcessing(false);
                    }
                } catch (error) {
                    console.error("Polling failed", error);
                }
            }, 3000); // Poll every 3 seconds
        }

        return () => {
            if (pollInterval) clearInterval(pollInterval);
        };
    }, [isPurchaseProcessing, isOwnedByUser, listingIdentifier]);

    const handleToggleLibrary = async () => {
        if (!listing) return;
        setIsLibraryActionLoading(true);
        try {
            if (isSavedInLibrary) {
                await UserAPI.library.delete(listing.id);
                setIsSavedInLibrary(false);
            } else {
                await UserAPI.library.add(listing.id);
                setIsSavedInLibrary(true);
            }
        } finally {
            setIsLibraryActionLoading(false);
        }
    };

    const handlePurchaseListing = async () => {
        if (!listing || isOwnedByUser || isPurchaseProcessing) return;

        setIsPurchaseActionLoading(true);
        try {
            await UserAPI.purchaseListing(listing.id);
            // Kick off the polling by updating the state
            setIsPurchaseProcessing(true);
        } catch (error) {
            console.error("Purchase failed", error);
        } finally {
            setIsPurchaseActionLoading(false);
        }
    };

    const formatCurrency = (amount: number): string => {
        return new Intl.NumberFormat("de-DE", {
            style: "currency",
            currency: "EUR",
        }).format(amount / 100);
    };

    if (isLoading) return <div className="flex items-center justify-center h-full text-slate-500 italic">Loading listing details...</div>;
    if (hasError || !listing) return <div className="flex items-center justify-center h-full text-slate-500">Listing not found</div>;

    console.log(listing)
    return (
        <div className="w-full h-full flex flex-col gap-8 px-6 py-8 max-w-7xl mx-auto">
            <section className="flex flex-col gap-6 border-b border-slate-200 dark:border-slate-800 pb-6">
                <div className="flex flex-col md:flex-row justify-between gap-6">
                    <TitleAndDescription
                        title={listing.title ?? "[No title]"}
                        description={listing.description ?? "[No description]"}
                        titleClassname="text-4xl md:text-5xl font-extrabold tracking-tight"
                        descriptionClassname="max-w-2xl text-slate-600 dark:text-slate-400"
                    />

                    <div className="flex flex-col items-start md:items-end gap-4 min-w-[240px]">
                        <div className="text-3xl font-bold text-slate-900 dark:text-white">
                            {formatCurrency(listing.price)}
                        </div>

                        <div className="flex items-center gap-2 w-full md:w-auto">
                            <button
                                onClick={handleToggleLibrary}
                                disabled={isLibraryActionLoading}
                                className={`flex items-center justify-center gap-2 px-4 py-2.5 rounded-lg border border-slate-300 dark:border-slate-700 hover:bg-slate-100 dark:hover:bg-slate-800 transition-all ${isSavedInLibrary ? "bg-slate-100 dark:bg-slate-800 text-blue-600" : ""} ${isLibraryActionLoading ? "opacity-50 cursor-not-allowed" : ""}`}
                            >
                                {isSavedInLibrary ? <BookmarkCheck className="w-4 h-4" /> : <BookmarkPlus className="w-4 h-4" />}
                                <span className="text-sm font-medium">{isSavedInLibrary ? "Saved" : "Save"}</span>
                            </button>

                            <button
                                onClick={handlePurchaseListing}
                                disabled={isPurchaseActionLoading || isOwnedByUser || isPurchaseProcessing}
                                className={`flex-1 md:flex-none flex items-center justify-center gap-2 px-6 py-2.5 rounded-lg font-semibold transition-all
                                    ${isOwnedByUser ? "bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-400 border border-green-200 dark:border-green-800" : 
                                      isPurchaseProcessing ? "bg-amber-100 dark:bg-amber-900/30 text-amber-700 dark:text-amber-400 border border-amber-200 dark:border-amber-800" :
                                      "bg-blue-600 hover:bg-blue-700 text-white shadow-sm active:scale-95"}
                                    ${(isPurchaseActionLoading || isOwnedByUser || isPurchaseProcessing) ? "cursor-default" : ""}
                                `}
                            >
                                {isOwnedByUser ? (
                                    <><CheckCircle className="w-4 h-4" /><span>Owned</span></>
                                ) : isPurchaseProcessing ? (
                                    <><Loader2 className="w-4 h-4 animate-spin" /><span>Processing</span></>
                                ) : (
                                    <><ShoppingCart className="w-4 h-4" /><span>{isPurchaseActionLoading ? "Authorizing..." : "Purchase"}</span></>
                                )}
                            </button>
                        </div>
                    </div>
                </div>
            </section>
            
            {isPurchaseProcessing && !isOwnedByUser && (
                <div className="flex items-center gap-3 bg-amber-50 dark:bg-amber-900/10 border border-amber-200 dark:border-amber-800 rounded-xl p-4 text-sm text-amber-800 dark:text-amber-200 animate-pulse">
                    <Loader2 className="w-4 h-4 animate-spin" />
                    <p>Your transaction is being finalized on the ledger. This page will update automatically once complete.</p>
                </div>
            )}

            <section className="flex flex-col flex-1 min-h-0 overflow-auto gap-6">
                {listing.screenshotIds && listing.screenshotIds.length > 0 && (
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                        {listing.screenshotIds.map((sId) => (
                            <div key={sId} className="aspect-video rounded-xl overflow-hidden border border-slate-200 dark:border-slate-800">
                                <ListingScreenshot 
                                    listingId={listing.id} 
                                    screenshotId={sId} 
                                    className="w-full h-full" 
                                />
                            </div>
                        ))}
                    </div>
                )}
                
                <div className="rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-50 dark:bg-slate-900 p-8 text-slate-700 dark:text-slate-300">
                   <h3 className="text-lg font-semibold mb-4 text-slate-900 dark:text-white">Listing Information</h3>
                    <p className="text-sm leading-relaxed mb-4">
                        This section contains deployment instructions, versioning history, and technical screenshots for <strong>{listing.title}</strong>.
                    </p>
                    {/* Placeholder for future listing components */}
                    <div className="h-32 border-2 border-dashed border-slate-200 dark:border-slate-800 rounded-lg flex items-center justify-center text-slate-400 italic">
                        Technical Documentation & Assets
                    </div>
                </div>
            </section>
        </div>
    );
}