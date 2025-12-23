import React, { useReducer, useState, useCallback, useEffect } from "react";
import { AnimatePresence, motion } from "framer-motion";
import EditableText from "../../components/input/EditableText";
import HardwareSettings from "./hardware/HardwareSettings";
import VisibilitySettings from "./VisibilitySettings";
import PriceSettings from "./prices/PriceSettings";
import DeleteListingPrompt from "./DeleteListingPrompt";
import { API, type Money, type DeveloperListing, type HardwareSpecification } from "../../backend";

interface DeveloperListingDisplayProps {
    listing: DeveloperListing;
    onShouldClose?: () => void;
}

type ListingAction =
    | { type: "set"; listing: DeveloperListing }
    | { type: "setTitle"; title: string }
    | { type: "setDescription"; description: string }
    | { type: "setAccessMode"; mode: number }
    | { type: "setHardware"; hardware: HardwareSpecification }
    | { type: "setPrice"; currency: Money }
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

// Memoized Section wrapper
const Section = React.memo(function Section({ title, description, children }: { title: string; description?: string; children: React.ReactNode }) {
    return (
        <div className="mb-8 border-b border-gray-300 dark:border-gray-700 pb-6">
            <h2 className="text-lg font-semibold mb-1 text-gray-900 dark:text-white">{title}</h2>
            {description && <p className="text-sm text-gray-500 dark:text-gray-400 mb-3">{description}</p>}
            {children}
        </div>
    );
});

export default function DeveloperListingDisplay({ listing: sourceListing, onShouldClose }: DeveloperListingDisplayProps) {
    const [backupListing, setBackupListing] = useState<DeveloperListing>(sourceListing);
    const [listing, dispatch] = useReducer(listingReducer, sourceListing);
    const [hasUnsavedChanges, setHasUnsavedChanges] = useState(false);

    const markChanged = useCallback(() => setHasUnsavedChanges(true), []);
    const resetChanges = useCallback(() => {
        dispatch({ type: "reset", backup: backupListing });
        setHasUnsavedChanges(false);
    }, [backupListing]);

    const saveChanges = useCallback(async () => {
        const { json, ok } = await API.developer.listing.update(listing);
        if (!ok) { resetChanges(); return; }
        dispatch({ type: "reset", backup: json as DeveloperListing });
        setBackupListing(json as DeveloperListing);
        setHasUnsavedChanges(false);
    }, [listing, resetChanges]);

    const deleteListing = useCallback(async () => {
        if (!(await API.developer.listing.delete(listing.id)).ok) onShouldClose?.();
    }, [listing.id, onShouldClose]);

    // Memoized callbacks for children
    const handleTitleChange = useCallback((title: string) => {
        dispatch({ type: "setTitle", title });
        markChanged();
    }, [markChanged]);

    const handleDescriptionChange = useCallback((description: string) => {
        dispatch({ type: "setDescription", description });
        markChanged();
    }, [markChanged]);

    const handleAccessModeChange = useCallback((mode: number) => {
        dispatch({ type: "setAccessMode", mode });
        markChanged();
    }, [markChanged]);

    const handlePriceChange = useCallback((price: Money) => {
        dispatch({ type: "setPrice", currency: price });
        markChanged();
    }, [markChanged]);

    const handleHardwareChange = useCallback((hardware: HardwareSpecification) => {
        dispatch({ type: "setHardware", hardware });
        markChanged();
    }, [markChanged]);

    return (
        <div className="relative flex flex-col w-full h-full items-center overflow-y-auto">
            <AnimatePresence mode="wait">
                {hasUnsavedChanges && (
                    <motion.div
                        initial={{ opacity: 0, y: 50 }}
                        animate={{ opacity: 1, y: 0 }}
                        exit={{ opacity: 0, y: -50 }}
                        transition={{ duration: 0.3 }}
                        className="fixed bottom-0 z-10 bg-indigo-100 border-l-4 border-indigo-500 shadow-md p-4 rounded-md mb-4 flex justify-between items-center mx-4"
                    >
                        <span className="text-indigo-900 font-semibold mr-10">You have unsaved changes</span>
                        <div className="flex gap-2">
                            <button className="px-3 py-1 bg-indigo-500 text-white rounded hover:bg-indigo-600" onClick={saveChanges}>
                                Save
                            </button>
                            <button className="px-3 py-1 bg-slate-300 text-slate-900 rounded hover:bg-slate-400" onClick={resetChanges}>
                                Cancel
                            </button>
                        </div>
                    </motion.div>
                )}
            </AnimatePresence>

            <div className="w-full p-6 max-w-4xl">
                <Section title="General" description="Basic information about your listing.">
                    <EditableText
                        text={listing.title ?? "No Title"}
                        label="Title: "
                        onChange={handleTitleChange}
                    />
                    <EditableText
                        text={listing.description ?? "No Description"}
                        label="Description: "
                        onChange={handleDescriptionChange}
                    />
                </Section>

                <Section title="Visibility" description="Control who can access this listing.">
                    <VisibilitySettings
                        accessMode={listing.accessMode ?? 0}
                        onChange={handleAccessModeChange}
                    />
                </Section>

                <Section title="Pricing" description="Set the price and currency for this listing.">
                    <PriceSettings
                        price={listing.price}
                        onChange={handlePriceChange}
                    />
                </Section>

                <Section title="Hardware" description="Configure hardware requirements.">
                    <HardwareSettings
                        initialSpec={listing.hardware ?? {}}
                        onChange={handleHardwareChange}
                    />
                </Section>

                <Section title="Danger Zone" description="Be careful! These actions are irreversible.">
                    <DeleteListingPrompt onDelete={deleteListing} />
                </Section>
            </div>
        </div>
    );
}
