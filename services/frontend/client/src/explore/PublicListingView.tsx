
import { useEffect, useState, type JSX } from "react";
import { useParams, useNavigate, type NavigateFunction } from "react-router-dom";

import "yet-another-react-lightbox/styles.css";
import { useUserSession } from "../components/restriction/UserSession";
import { ListingAPI, type Listing } from "../backend/repositories/listing";
import { UserAPI } from "../backend/repositories/user";
import { ListingHeader } from "./listing/ListingHeader";
import { ListingGallery } from "./listing/ListingGallery";
import { ListingSpecs } from "./listing/ListingSpecs";
import HardwareSpecsDisplay from "../components/view/HardwareSpecDisplay";
import MarkdownViewer from "../components/view/MarkdownViewer";
import MainNavbar from "../components/navigation/MainNavbar";


export function gotoListing(navigate: NavigateFunction, listing: Listing | undefined) {
    if (!listing) return;
    navigate("/listing/" + listing.id);
}

export function PublicListingView(): JSX.Element {
    const { id: listingIdentifier } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const { user } = useUserSession();
    const isAuthenticated = !!user;

    const [listing, setListing] = useState<Listing | null>(null);
    const [isLoading, setIsLoading] = useState(true);
    const [hasError, setHasError] = useState(false);
    
    const [isSavedInLibrary, setIsSavedInLibrary] = useState(false);
    const [isOwnedByUser, setIsOwnedByUser] = useState(false);
    const [isPurchaseProcessing, setIsPurchaseProcessing] = useState(false);
    
    const [isLibLoading, setIsLibLoading] = useState(false);
    const [isPurLoading, setIsPurLoading] = useState(false);

    useEffect(() => {
        async function initialize() {
            if (!listingIdentifier) { setHasError(true); setIsLoading(false); return; }
            try {
                const res = await ListingAPI.get(listingIdentifier);
                setListing(res);
                if (isAuthenticated) {
                    const [lib, own, proc] = await Promise.all([
                        UserAPI.library.check(listingIdentifier),
                        UserAPI.isOwner(listingIdentifier),
                        UserAPI.isPurchaseProcessingOwner(listingIdentifier)
                    ]);
                    setIsSavedInLibrary(lib); setIsOwnedByUser(own); setIsPurchaseProcessing(proc);
                }
            } catch { setHasError(true); } 
            finally { setIsLoading(false); }
        }
        initialize();
    }, [listingIdentifier, isAuthenticated]);

    // Handlers
    const handleLogin = () => navigate("/login", { state: { from: `/listing/${listingIdentifier}` } });
    
    const handleSave = async () => {
        if (!isAuthenticated) return handleLogin();
        setIsLibLoading(true);
        try {
            isSavedInLibrary ? await UserAPI.library.delete(listing!.id) : await UserAPI.library.add(listing!.id);
            setIsSavedInLibrary(!isSavedInLibrary);
        } finally { setIsLibLoading(false); }
    };

    /**
     * Polling configuration for purchase verification.
     */
    const MAX_PURCHASE_POLL_ATTEMPTS = 5;
    const PURCHASE_POLL_INTERVAL_MS = 3000;

    /**
     * Executes the purchase flow and initiates status polling.
     */
    const handlePurchase = async () => {
        if (!isAuthenticated) {
            return handleLogin();
        }

        setIsPurLoading(true);

        try {
            await UserAPI.purchaseListing(listing!.id);
            setIsPurchaseProcessing(true);
            startPurchaseStatusPolling();
        } finally {
            setIsPurLoading(false);
        }
    };

    /**
     * Initiates a recursive poll to verify if the purchase has been finalized.
     */
    async function startPurchaseStatusPolling(attemptCount: number = 0) {
        if (attemptCount >= MAX_PURCHASE_POLL_ATTEMPTS) {
            return;
        }

        const isStillProcessing = await UserAPI.isPurchaseProcessingOwner(listing!.id);

        if (isStillProcessing) {
            return scheduleNextPurchasePoll(attemptCount);
        }

        await finalizePurchaseState();
    }

    /**
     * Schedules the next poll attempt after a set delay.
     */
    function scheduleNextPurchasePoll(currentAttempt: number) {
        setTimeout(() => {
            startPurchaseStatusPolling(currentAttempt + 1);
        }, PURCHASE_POLL_INTERVAL_MS);
    }

    /**
     * Updates the component state once the purchase is no longer processing.
     */
    async function finalizePurchaseState() {
        const isNowOwned = await UserAPI.isOwner(listing!.id);

        setIsPurchaseProcessing(false);
        setIsOwnedByUser(isNowOwned);
    }

    if (isLoading) return <div className="p-10 text-slate-500 italic">Loading...</div>;
    if (hasError || !listing) return <div className="p-10">Listing not found</div>;
    
    return (
        <div>
            <MainNavbar></MainNavbar>

            <div className="w-full h-full flex flex-col gap-8 px-6 py-8 max-w-7xl mx-auto">
                <ListingHeader 
                    listing={listing}
                    isAuthenticated={isAuthenticated}
                    isSaved={isSavedInLibrary}
                    isOwned={isOwnedByUser}
                    isProcessing={isPurchaseProcessing}
                    isActionLoading={{ library: isLibLoading, purchase: isPurLoading }}
                    onSave={handleSave}
                    onPurchase={handlePurchase}
                    onLogin={handleLogin}
                />

                <div className="grid grid-cols-1 lg:grid-cols-12 gap-8">
                    <div className="lg:col-span-8">
                        <ListingGallery listingId={listing.id} screenshotIds={listing.screenshotIds ?? []} />
                    </div>
                    <div className="lg:col-span-4">
                        {/*<ListingSpecs /> <div className="h-10"></div>*/}
                        
                        {listing.recommended_hardware && <HardwareSpecsDisplay specification={listing.recommended_hardware}></HardwareSpecsDisplay> }
                    </div>
                </div>
                <div className="w-full bg-white p-8 rounded-3xl border border-slate-200">
                    <MarkdownViewer content={listing.documentation_markdown} />
                </div>
            </div>

        </div>
    );
}