import React, { useReducer, useState, useCallback, useEffect } from "react";
import { AnimatePresence, motion } from "framer-motion";
import HardwareSettings from "./hardware/HardwareSettings";
import VisibilitySettings from "./VisibilitySettings";
import PriceSettings from "./prices/PriceSettings";
import DeleteListingPrompt from "./DeleteListingPrompt";

import { DeveloperListingAPI, type DeveloperListing, type ListingGithubSetup } from "../../backend/repositories/developer_listing";
import type { HardwareSpecification, ListingAccessMode } from "../../backend/types";
import toast, { Toaster } from "react-hot-toast";
import { Button } from "../../templates/components/Button";
import Section from "../../templates/components/Section";
import EditableText from "../../templates/components/EditableText";
import CustomSelect from "../../templates/components/SelectBox";
import { ListingAPI } from "../../backend/repositories/listing";

interface DeveloperListingDisplayProps {
    listing: DeveloperListing;
    onShouldClose?: () => void;
}

type ListingAction =
    | { type: "set"; listing: DeveloperListing }
    | { type: "setTitle"; title: string }
    | { type: "setDescription"; description: string }
    | { type: "setAccessMode"; mode: ListingAccessMode }
    | { type: "setHardware"; hardware: HardwareSpecification }
    | { type: "setPrice"; currency: number }
    | { type: "reset"; backup: DeveloperListing }
    | { type: "addScreenshot"; id: string }
    | { type: "removeScreenshot"; id: string };

function listingReducer(state: DeveloperListing, action: ListingAction): DeveloperListing {
    switch (action.type) {
        case "setTitle": return { ...state, title: action.title };
        case "setDescription": return { ...state, description: action.description };
        case "setAccessMode": return { ...state, accessMode: action.mode };
        case "setHardware": return { ...state, hardware: { ...state.hardware, ...action.hardware } };
        case "setPrice": return { ...state, price: action.currency };
        case "reset": return action.backup;
        case "set": return action.listing;
        case "addScreenshot": 
            return { ...state, screenshotIds: [...(state.screenshotIds || []), action.id] };
        case "removeScreenshot":
            return { ...state, screenshotIds: (state.screenshotIds || []).filter(id => id !== action.id) };
        default: return state;
    }
}

export default function DeveloperListingEditor({ listing: sourceListing, onShouldClose }: DeveloperListingDisplayProps) {
    const [backupListing, setBackupListing] = useState<DeveloperListing>(sourceListing);
    const [listing, dispatch] = useReducer(listingReducer, sourceListing);
    const [hasUnsavedChanges, setHasUnsavedChanges] = useState(false);

    const markChanged = useCallback(() => setHasUnsavedChanges(true), []);
    const resetChanges = useCallback(() => {
        dispatch({ type: "reset", backup: backupListing });
        setHasUnsavedChanges(false);
    }, [backupListing]);

    const saveChanges = useCallback(async () => {
        try {
            const updated = await DeveloperListingAPI.update(listing);
            dispatch({ type: "reset", backup: updated });
            setBackupListing(updated);
            setHasUnsavedChanges(false);
        } catch (err) {
            toast.error(err + "");
            resetChanges();
        }
    }, [listing, resetChanges]);

    const deleteListing = useCallback(async () => {
        try {
            await DeveloperListingAPI.delete(listing.id);
            onShouldClose?.();
        } catch (err) {
            toast.error(err + "");
        }
    }, [listing.id, onShouldClose]);

    // ---------------- GitHub Setup ----------------
    const [setup, setSetup] = useState<ListingGithubSetup | null>(null);
    const [loadingSetup, setLoadingSetup] = useState(true);
    const [setupDraft, setSetupDraft] = useState<{ repoUrl: string; accessToken: string }>({ repoUrl: "", accessToken: "" });

    useEffect(() => {
        let canceled = false;
        async function fetchSetup() {
            setLoadingSetup(true);
            try {
                const setupData = await DeveloperListingAPI.githubSetup.get(listing.id);
                if (!canceled) {
                    setSetup(setupData);
                    setSetupDraft({ repoUrl: setupData?.repositoryURL ?? "", accessToken: setupData?.accessToken ?? "" });
                }
            } catch (err) {
                toast.error(err + "");
            }
            setLoadingSetup(false);
        }
        fetchSetup();
        return () => { canceled = true; };
    }, [listing.id]);

    const saveSetup = useCallback(async () => {
        setLoadingSetup(true);
        try {
            const updatedSetup = await DeveloperListingAPI.githubSetup.attachOrUpdate(
                listing.id,
                setupDraft.repoUrl,
                setupDraft.accessToken
            );
            setSetup(updatedSetup);
        } catch (err) {
            toast.error(err + "");
        }
        setLoadingSetup(false);
    }, [listing.id, setupDraft]);

    const removeSetup = useCallback(async () => {
        setLoadingSetup(true);
        try {
            await DeveloperListingAPI.githubSetup.remove(listing.id);
            setSetup(null);
            setSetupDraft({ repoUrl: "", accessToken: "" });
        } catch (err) {
            toast.error(err + "");
        }
        setLoadingSetup(false);
    }, [listing.id]);
    // ... existing state ...
    const [isUploading, setIsUploading] = useState(false);
    const [screenshotUrls, setScreenshotUrls] = useState<Record<string, string>>({});

    // ---------------- Screenshot Link Caching ----------------

    const fetchScreenshotUrls = useCallback(async (ids: string[]) => {
        const newUrls: Record<string, string> = {};
        
        await Promise.all(ids.map(async (id) => {
            try {
                // Use the new endpoint we created in the Go backend
                const response = await ListingAPI.getScreenshotReadURL(listing.id, id);
                newUrls[id] = response.readUrl;
            } catch (err) {
                console.error(`Failed to fetch URL for screenshot ${id}`, err);
            }
        }));

        setScreenshotUrls(prev => ({ ...prev, ...newUrls }));
    }, [listing.id]);

    // Fetch URLs whenever the listing's screenshot ID list changes
    useEffect(() => {
        if (listing.screenshotIds && listing.screenshotIds.length > 0) {
            const missingIds = listing.screenshotIds.filter(id => !screenshotUrls[id]);
            if (missingIds.length > 0) {
                fetchScreenshotUrls(missingIds);
            }
        }
    }, [listing.screenshotIds, fetchScreenshotUrls, screenshotUrls]);

    // ---------------- Modified Upload Logic ----------------

    const handleFileUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (!file) return;

        setIsUploading(true);
        const toastId = toast.loading("Uploading screenshot...");

        try {
            // 1. Get signed PUT URL
            const response = await DeveloperListingAPI.screenshots.createUploadUrl(listing.id);
            
            // 2. Upload to MinIO
            const uploadResponse = await fetch(response.uploadUrl, {
                method: "PUT",
                body: file,
                headers: { "Content-Type": file.type }
            });

            if (!uploadResponse.ok) throw new Error("Cloud storage upload failed");

            // 3. Refresh Listing Data
            const updated = await DeveloperListingAPI.get(listing.id);
            dispatch({ type: "set", listing: updated });
            
            // 4. Specifically fetch the URL for the new screenshot immediately
            if (updated.screenshotIds) {
                await fetchScreenshotUrls(updated.screenshotIds);
            }

            toast.success("Screenshot uploaded", { id: toastId });
        } catch (err) {
            toast.error("Upload failed: " + err, { id: toastId });
        } finally {
            setIsUploading(false);
            e.target.value = "";
        }
    };

    const handleDeleteScreenshot = async (screenshotId: string) => {
        try {
            await DeveloperListingAPI.screenshots.delete(listing.id, screenshotId);
            dispatch({ type: "removeScreenshot", id: screenshotId });
            toast.success("Screenshot removed");
        } catch (err) {
            toast.error("Delete failed: " + err);
        }
    };

    return (
        <div className="relative flex flex-col w-full h-full items-center overflow-y-auto">
            <Toaster position="top-right" />
            
            <AnimatePresence mode="wait">
                {hasUnsavedChanges && (
                    <motion.div
                        initial={{ opacity: 0, y: 50 }}
                        animate={{ opacity: 1, y: 0 }}
                        exit={{ opacity: 0, y: -50 }}
                        transition={{ duration: 0.3 }}
                        className="fixed bottom-0 z-10 bg-indigo-100 dark:bg-indigo-900 border-l-4 border-indigo-500 shadow-md p-4 rounded-md mb-4 flex justify-between items-center mx-4"
                    >
                        <span className="text-indigo-900 dark:text-indigo-100 font-semibold mr-10">You have unsaved changes</span>
                        <div className="flex gap-2">
                            <Button variant="primary" onClick={saveChanges}>Save</Button>
                            <Button variant="secondary" onClick={resetChanges}>Cancel</Button>
                        </div>
                    </motion.div>
                )}
            </AnimatePresence>

            <div className="w-full p-6 max-w-4xl">
                <Section title="General" description="Basic information about your listing.">
                    <EditableText
                        text={listing.title ?? "No Title"}
                        label="Title:"
                        block
                        onChange={(t) => { dispatch({ type: "setTitle", title: t }); markChanged(); }}
                    />
                    <EditableText
                        text={listing.description ?? "No Description"}
                        label="Description:"
                        block
                        onChange={(d) => { dispatch({ type: "setDescription", description: d }); markChanged(); }}
                    />
                </Section>

                <Section title="Visibility" description="Control who can access this listing.">
                    <CustomSelect<number>
                        value={listing.accessMode ?? 0}
                        onChange={(m) => { dispatch({ type: "setAccessMode", mode: m as ListingAccessMode }); markChanged(); }}
                        options={[
                            { label: "Public", value: 0 },
                            { label: "Private", value: 1 },
                        ]}
                        placeholder="Select access mode"
                        block
                    />
                </Section>

                <Section title="Pricing" description="Set the price and currency for this listing.">
                    <PriceSettings
                        price={listing.price}
                        onChange={(p) => { dispatch({ type: "setPrice", currency: p }); markChanged(); }}
                    />
                </Section>

                <Section title="Hardware" description="Configure hardware requirements.">
                    <HardwareSettings
                        initialSpec={listing.hardware ?? {}}
                        onChange={(h) => { dispatch({ type: "setHardware", hardware: h }); markChanged(); }}
                    />
                </Section>

                <Section title="Screenshots" description="Visuals for your listing (Max 5).">
                    <div className="grid grid-cols-2 md:grid-cols-3 gap-4 mb-4">
                        {listing.screenshotIds?.map((id) => (
                            <div key={id} className="relative group aspect-video bg-gray-200 dark:bg-gray-700 rounded-lg overflow-hidden border dark:border-gray-600">
                                {screenshotUrls[id] ? (
                                    <motion.img 
                                        initial={{ opacity: 0 }}
                                        animate={{ opacity: 1 }}
                                        src={screenshotUrls[id]} 
                                        alt="Screenshot" 
                                        className="object-cover w-full h-full"
                                    />
                                ) : (
                                    <div className="w-full h-full flex items-center justify-center">
                                        <div className="animate-pulse bg-gray-300 dark:bg-gray-600 w-full h-full" />
                                    </div>
                                )}
                                
                                <button 
                                    onClick={() => handleDeleteScreenshot(id)}
                                    className="absolute top-2 right-2 bg-red-600/80 hover:bg-red-600 text-white p-1.5 rounded-full opacity-0 group-hover:opacity-100 transition-opacity backdrop-blur-sm"
                                >
                                    <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                                    </svg>
                                </button>
                            </div>
                        ))}
                        
                         {(listing.screenshotIds?.length ?? 0) < 5 && (

                            <label className="flex flex-col items-center justify-center aspect-video border-2 border-dashed border-gray-300 dark:border-gray-600 rounded-lg cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors">

                                <div className="flex flex-col items-center justify-center pt-5 pb-6">

                                    {isUploading ? (

                                        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-500"></div>

                                    ) : (

                                        <>

                                            <svg className="w-8 h-8 mb-3 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 4v16m8-8H4"></path></svg>

                                            <p className="text-xs text-gray-500">Add Screenshot</p>

                                        </>

                                    )}

                                </div>

                                <input type="file" className="hidden" onChange={handleFileUpload} disabled={isUploading} accept="image/*" />

                            </label>

                        )}
                    </div>
                </Section>

                <Section title="GitHub Setup" description="Attach a GitHub repository for automated setup.">
                    <div className="flex flex-col gap-2">
                        <input
                            type="text"
                            placeholder="GitHub Repo URL"
                            className="px-3 py-1 border rounded w-full dark:bg-gray-800 dark:border-gray-600 dark:text-white"
                            value={setupDraft.repoUrl}
                            onChange={(e) => setSetupDraft({ ...setupDraft, repoUrl: e.target.value })}
                        />
                        <input
                            type="text"
                            placeholder="Access Token (optional)"
                            className="px-3 py-1 border rounded w-full dark:bg-gray-800 dark:border-gray-600 dark:text-white"
                            value={setupDraft.accessToken}
                            onChange={(e) => setSetupDraft({ ...setupDraft, accessToken: e.target.value })}
                        />
                        <div className="flex gap-2">
                            <Button disabled={loadingSetup} variant="success" onClick={saveSetup}>
                                {setup ? "Update Setup" : "Attach Setup"}
                            </Button>
                            {setup && (
                                <Button disabled={loadingSetup} variant="danger" onClick={removeSetup}>
                                    Remove Setup
                                </Button>
                            )}
                        </div>
                    </div>
                </Section>

                <Section title="Danger Zone" description="Be careful! These actions are irreversible.">
                    <DeleteListingPrompt onDelete={deleteListing} />
                </Section>
            </div>
        </div>
    );
}
