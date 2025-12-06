// src/components/AccountPage.tsx
import React, { useReducer, useState } from "react";
import EditableText from "../../components/input/EditableText";
import { API, ListingAccessModes, type CurrencyRequest, type HardwareUpdate, type Listing, type PricingEntry } from "../../backend";
import { ChevronLeft, Delete } from "lucide-react";
import PromptButton from "../../components/input/PromptButton";

import HardwareSettings from "./hardware/HardwareSettings";

import { FlowSwitch } from "../../components/navigation/Flow";
import { AnimatePresence, motion } from "framer-motion";
import SelectBox from "../../components/input/SelectBox";
import PricesMenu from "./prices/PricesMenu";

interface ListingDisplayProps {
    listing: Listing,
    onShouldClose?: () => void
}

export type ListingAction =
    | { type: "setTitle"; title: string }
    | { type: "setDescription"; description: string }
    | { type: "setAccessMode"; mode: number }
    | { type: "setHardware"; hardware: HardwareUpdate }
    | { type: "addPrice"; price: PricingEntry }
    | { type: "updatePrice"; priceId: string; currency: CurrencyRequest | null }
    | { type: "deletePrice"; priceId: string }
    | { type: "reset"; backup: Listing };

    
function listingReducer(state: Listing, action: ListingAction): Listing {
    switch (action.type) {
        case "setTitle":
            return { ...state, title: action.title };
        case "setDescription":
            return { ...state, description: action.description };
        case "setAccessMode":
            return { ...state, accessMode: action.mode };
        case "setHardware":
            return { ...state, hardware: { ...state.hardware, ...action.hardware } };
        case "addPrice":
            return { ...state, prices: [...state.prices, action.price] };
        case "updatePrice":
            return {
                ...state,
                prices: state.prices.map((p) =>
                    p.id === action.priceId ? { ...p, currency: action.currency } : p
                ),
            };
        case "deletePrice":
            return { ...state, prices: state.prices.filter((p) => p.id !== action.priceId) };
        case "reset":
            return action.backup;
        default:
            return state;
    }
}
export default function DeveloperListingDisplay({ listing: sourceListing, onShouldClose }: ListingDisplayProps) {
    const [backupListing, setBackupListing] = useState<Listing>(sourceListing);
    const [listing, dispatch] = useReducer(listingReducer, sourceListing);

    const [currentText, setText] = useState<string>("");
    const deleteKeyword = "Delete";
    const isDeleteLocked = currentText !== deleteKeyword;
    const [hasUnsavedChanges, setHasUnsavedChanges] = useState<boolean>(false);

    const handleChange = (event: React.KeyboardEvent<HTMLInputElement>) => {
        setText((event.target as HTMLInputElement).value);
    };

    const deleteThisListing = async () => {
        if (!(await API.developer.listing.delete(listing.id)).ok) onShouldClose?.();
    };

    const reset = () => {
        dispatch({ type: "reset", backup: backupListing });
        setHasUnsavedChanges(false);
    };

    const changeListing = () => setHasUnsavedChanges(true);

    return (
        <div className="relative flex flex-col w-full h-full min-h-0 overflow-hidden">
            <FlowSwitch className="my-3 py-2 flex flex-row items-center hover:bg-slate-200 dark:hover:bg-gray-700 transition-all" direction="back">
                <ChevronLeft size={32} />
                <label>Back</label>
            </FlowSwitch>

            <AnimatePresence mode="wait">
                {hasUnsavedChanges && <motion.div
                    initial={{ opacity: 0, y: 50 }}
                    animate={{ opacity: 1, y: 0 }}
                    exit={{ opacity: 0, y: -50 }}
                    transition={{ duration: 0.3 }}
                    className="absolute bottom-0 left-0 right-0 bg-indigo-100 dark:bg-slate-700 border-l-4 border-indigo-500 shadow-md p-4 rounded-md mb-4 flex justify-between items-center mx-4"
                >
                    <span className="text-indigo-900 dark:text-white font-semibold">You have unsaved changes</span>
                    <div className="flex gap-2">
                        <button
                            className="px-3 py-1 bg-indigo-500 text-white rounded hover:bg-indigo-600"
                            onClick={async () => {
                                if (listing.prices) {
                                    const removed = backupListing.prices.filter(item => !listing.prices.includes(item));
                                    
                                    for(const price of removed){
                                        await API.developer.listing.prices.delete(listing.id, price.id);
                                    }
                                }

                                const{ json, ok } = await API.developer.listing.update(listing);
                                if (!ok) {
                                    reset();
                                    return;
                                }

                                dispatch({ type: "reset", backup: json as Listing });
                                setBackupListing(json as Listing);
                                setHasUnsavedChanges(false);
                            }}
                        >
                            Save
                        </button>
                        <button
                            className="px-3 py-1 bg-slate-300 dark:bg-slate-700 text-slate-900 dark:text-white rounded hover:bg-slate-400"
                            onClick={reset}
                        >
                            Cancel
                        </button>
                    </div>
                </motion.div>}
            </AnimatePresence>

            <div className="flex-1 overflow-y-auto">
                <div className="m-3">
                    <div className="flex flex-row justify-between mb-10">
                        <SelectBox
                            options={ListingAccessModes}
                            onSelected={(mode) => {
                                dispatch({ type: "setAccessMode", mode: parseInt(mode) });
                                changeListing();
                            }}
                            defaultValue={"" + listing.accessMode}
                        />
                        <PromptButton
                            className="bg-red-500 hover:bg-red-600 text-white font-semibold py-2 px-4 rounded transition-colors dark:bg-red-600 dark:hover:bg-red-700 flex flex-row gap-2"
                            submitClassName={`flex flex-row gap-2 font-semibold py-2 px-4 rounded transition-colors ${isDeleteLocked
                                ? "bg-red-300 text-gray-200 cursor-not-allowed dark:bg-red-700 dark:text-gray-400"
                                : "bg-red-500 hover:bg-red-600 text-white dark:bg-red-600 dark:hover:bg-red-700"
                                }`}
                            promptContents={
                                <>
                                    <p>
                                        This operation will delete this listing. If you are sure you want to do this type "Delete" as confirmation.
                                    </p>
                                    <input
                                        onKeyUp={handleChange}
                                        type="text"
                                        className="w-full p-3 my-10 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-red-500 focus:border-red-500 dark:bg-gray-700 dark:border-gray-600 dark:text-gray-100 dark:placeholder-gray-400 dark:focus:ring-red-400 dark:focus:border-red-400"
                                    ></input>
                                </>
                            }
                            isSubmitLocked={() => isDeleteLocked}
                            onCancel={() => setText("")}
                            onSubmit={deleteThisListing}
                        >
                            Delete <Delete />
                        </PromptButton>
                    </div>

                    <EditableText
                        text={listing.title ?? "No Title"}
                        label="Title: "
                        onChange={(title: string) => {
                            dispatch({ type: "setTitle", title });
                            changeListing();
                        }}
                    />

                    <EditableText
                        text={listing.description ?? "No Description"}
                        label="Description: "
                        onChange={(description: string) => {
                            dispatch({ type: "setDescription", description });
                            changeListing();
                        }}
                    />

                    {/* Hardware sliders */}
                    <HardwareSettings
                        hardware={listing.hardware ?? {}}
                        onChange={(newHardware) => {
                            dispatch({ type: "setHardware", hardware: newHardware });
                            changeListing();
                        }}
                    />

                    {/* Prices menu */}
                    <PricesMenu
                        listing={listing}
                        onAction={(action) => {dispatch(action); changeListing()}}
                    />
                </div>
            </div>
        </div>
    );
}
