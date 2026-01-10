import { useEffect, useState, type JSX } from "react";
import { useParams, type NavigateFunction } from "react-router-dom";
import { TitleAndDescription } from "../components/prefabs/TitleAndDescription";
import { BookmarkPlus, BookmarkCheck, ShoppingCart, CheckCircle } from "lucide-react";
import { ListingAPI, type Listing } from "../backend/repositories/listing";
import { UserAPI } from "../backend/repositories/user";

/**
 * Navigates the user to the detailed view of a specific listing.
 */
export function gotoListing(
    navigate: NavigateFunction,
    listing: Listing | undefined
) {
    if (!listing) return;
    navigate("/listing/" + listing.id);
}

/**
 * PublicListingView displays the details of a listing and allows users
 * to save it to their library or purchase it.
 */
export function PublicListingView(): JSX.Element {
    const { id: listingIdentifier } = useParams<{ id: string }>();

    const [listing, setListing] = useState<Listing | null>(null);
    const [isLoading, setIsLoading] = useState<boolean>(true);
    const [hasError, setHasError] = useState<boolean>(false);
    
    const [isSavedInLibrary, setIsSavedInLibrary] = useState<boolean>(false);
    const [isOwnedByUser, setIsOwnedByUser] = useState<boolean>(false);
    
    const [isLibraryActionLoading, setIsLibraryActionLoading] = useState<boolean>(false);
    const [isPurchaseActionLoading, setIsPurchaseActionLoading] = useState<boolean>(false);

    useEffect(() => {
        /**
         * Loads the listing details and checks the user's relationship with it.
         */
        async function initializeListingView() {
            if (!listingIdentifier) {
                setHasError(true);
                setIsLoading(false);
                return;
            }

            try {
                const [listingResult, libraryStatus, ownershipStatus] = await Promise.all([
                    ListingAPI.get(listingIdentifier),
                    UserAPI.library.check(listingIdentifier),
                    UserAPI.isOwner(listingIdentifier)
                ]);

                setListing(listingResult);
                setIsSavedInLibrary(libraryStatus);
                setIsOwnedByUser(ownershipStatus);
            }
            catch (error) {
                setHasError(true);
            }
            finally {
                setIsLoading(false);
            }
        }

        initializeListingView();
    }, [listingIdentifier]);

    /**
     * Toggles the presence of the listing in the user's personal library.
     */
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

    /**
     * Initiates the purchase flow for the current listing.
     */
    const handlePurchaseListing = async () => {
        if (!listing || isOwnedByUser) return;

        setIsPurchaseActionLoading(true);
        try {
            await UserAPI.purchaseListing(listing.id);
            setIsOwnedByUser(true);
        } catch (error) {
            console.error("Purchase failed", error);
        } finally {
            setIsPurchaseActionLoading(false);
        }
    };

    /**
     * Formats the numeric price into a localized EUR currency string.
     */
    const formatCurrency = (amount: number): string => {
        return new Intl.NumberFormat("de-DE", {
            style: "currency",
            currency: "EUR",
        }).format(amount);
    };

    if (isLoading) {
        return (
            <div className="flex items-center justify-center h-full text-slate-500">
                Loading listing…
            </div>
        );
    }

    if (hasError || !listing) {
        return (
            <div className="flex items-center justify-center h-full text-slate-500">
                Listing not found
            </div>
        );
    }

    return (
        <div className="w-full h-full flex flex-col gap-8 px-6 py-8 max-w-7xl mx-auto">
            <section className="flex flex-col gap-6 border-b border-slate-200 dark:border-slate-800 pb-6">
                <div className="flex flex-col md:flex-row justify-between gap-6">
                    <TitleAndDescription
                        title={listing.title ?? "[No title]"}
                        description={listing.description ?? "[No description]"}
                        titleClassname="text-4xl md:text-5xl"
                        descriptionClassname="max-w-3xl"
                    />

                    <div className="flex flex-col items-start md:items-end gap-4">
                        <div className="text-3xl font-bold text-slate-900 dark:text-white">
                            {formatCurrency(listing.price)}
                        </div>

                        <div className="flex items-center gap-2">
                            <button
                                onClick={handleToggleLibrary}
                                disabled={isLibraryActionLoading}
                                className={`
                                    flex items-center gap-2
                                    px-4 py-2.5
                                    rounded-lg
                                    border
                                    border-slate-300
                                    dark:border-slate-700
                                    hover:bg-slate-100
                                    dark:hover:bg-slate-800
                                    transition
                                    ${isSavedInLibrary ? "bg-slate-100 dark:bg-slate-800" : ""}
                                    ${isLibraryActionLoading ? "opacity-50 cursor-not-allowed" : ""}
                                `}
                            >
                                {isSavedInLibrary ? (
                                    <BookmarkCheck className="w-4 h-4 text-blue-500" />
                                ) : (
                                    <BookmarkPlus className="w-4 h-4" />
                                )}
                                <span className="text-sm font-medium">
                                    {isSavedInLibrary ? "Saved to Library" : "Save to Library"}
                                </span>
                            </button>

                            <button
                                onClick={handlePurchaseListing}
                                disabled={isPurchaseActionLoading || isOwnedByUser}
                                className={`
                                    flex items-center gap-2
                                    px-6 py-2.5
                                    rounded-lg
                                    font-semibold
                                    transition
                                    ${isOwnedByUser 
                                        ? "bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-400 border border-green-200 dark:border-green-800" 
                                        : "bg-blue-600 hover:bg-blue-700 text-white shadow-sm"
                                    }
                                    ${(isPurchaseActionLoading || isOwnedByUser) ? "cursor-default" : "active:scale-95"}
                                `}
                            >
                                {isOwnedByUser ? (
                                    <>
                                        <CheckCircle className="w-4 h-4" />
                                        <span>Purchased</span>
                                    </>
                                ) : (
                                    <>
                                        <ShoppingCart className="w-4 h-4" />
                                        <span>{isPurchaseActionLoading ? "Processing..." : "Purchase Now"}</span>
                                    </>
                                )}
                            </button>
                        </div>
                    </div>
                </div>
            </section>

            <section className="flex flex-col flex-1 min-h-0 overflow-auto gap-6">
                <div className="rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-50 dark:bg-slate-900 p-6 text-slate-700 dark:text-slate-300">
                    <p className="text-sm leading-relaxed">
                        Additional listing content goes here. Deployment instructions, screenshots, supported providers, versioning, etc.
                    </p>
                </div>
            </section>
        </div>
    );
}