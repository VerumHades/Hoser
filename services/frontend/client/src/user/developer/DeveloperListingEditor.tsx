import React, { useReducer, useState, useCallback, useEffect } from "react";
import { AnimatePresence, motion } from "framer-motion";
import HardwareSettings from "./hardware/HardwareSettings";
import VisibilitySettings from "./VisibilitySettings";
import PriceSettings from "./prices/PriceSettings";
import DeleteListingPrompt from "./DeleteListingPrompt";

import { DeveloperListingAPI, type DeveloperListing, type ListingGithubSetup } from "../../backend/repositories/developer_listing";
import type { HardwareSpecification, ListingAccessMode } from "../../backend/types";
import toast from "react-hot-toast";
import { Button } from "../../templates/components/Button";
import Section from "../../templates/components/Section";
import EditableText from "../../templates/components/EditableText";
import CustomSelect from "../../templates/components/SelectBox";

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
    | { type: "reset"; backup: DeveloperListing };

function listingReducer(state: DeveloperListing, action: ListingAction): DeveloperListing {
    switch (action.type) {
        case "setTitle": return { ...state, title: action.title };
        case "setDescription": return { ...state, description: action.description };
        case "setAccessMode": return { ...state, accessMode: action.mode };
        case "setHardware": return { ...state, hardware: { ...state.hardware, ...action.hardware } };
        case "setPrice": return { ...state, price: action.currency };
        case "reset": return action.backup;
        case "set": return action.listing;
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

    return (
        <div className="relative flex flex-col w-full h-full items-center overflow-y-auto">
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
