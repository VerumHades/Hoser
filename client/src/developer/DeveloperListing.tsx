// src/components/AccountPage.tsx
import React, { use, useRef, useState } from "react";
import EditableText from "../components/EditableText";
import { API, ListingAccessModes, type HardwareUpdate, type Listing, type PricesRequest } from "../backend";
import { ChevronLeft, Delete } from "lucide-react";
import PromptButton from "../components/PromptButton";
import SelectBox from "../components/SelectBox";
import HardwareSettings from "../components/prefabs/HardwareSettings";
import PricesMenu from "../components/prefabs/PricesMenu";
import { FlowSwitch } from "../components/Flow";

interface ListingDisplayProps {
    listing: Listing,
    onShouldClose?: () => void
}


export default function DeveloperListingDisplay({ listing: sourceListing, onShouldClose }: ListingDisplayProps) {
    const [listing, setListing] = useState<Listing>(sourceListing);

    const [currentText, setText] = useState<string>("");
    const deleteKeyword = "Delete"
    const isDeleteLocked = currentText != deleteKeyword

    const [hardware, setHardware] = useState<HardwareUpdate | undefined>(sourceListing.hardware);
    const [prices, setPrices] = useState<PricesRequest | undefined>(sourceListing.prices);

    const [title, setTitle] = useState<string | undefined>(listing.title)
    const [description, setDescription] = useState<string | undefined>(listing.description)

    const [hasUnsavedChanges, setHasUnsavedChanges] = useState<boolean>(false)

    const changeListing = () => {
        setHasUnsavedChanges(true)
    }
    const handleChange = (event: React.KeyboardEvent<HTMLInputElement>) => {
        setText((event.target as HTMLInputElement).value);
    };

    const deleteThisListing = async () => {
        if (!(await API.developer.deleteListing(listing.id)).ok)
            onShouldClose?.()
    }

    const onSelectedAccessMode = (mode: string) => {
        API.developer.editListing({ id: listing.id, accessMode: parseInt(mode) as (0 | 1) })
    }

    return (
        <div className="flex flex-col w-full h-full min-h-0 overflow-hidden">
            <FlowSwitch
                className="my-3 py-2 flex flex-row items-center hover:bg-slate-200 dark:hover:bg-gray-700 transition-all"
                direction="back"
            >
                <ChevronLeft size={32} />
                <label>Back</label>
            </FlowSwitch>
            {hasUnsavedChanges && (
                <div className="bg-indigo-100 dark:bg-slate-700 border-l-4 border-indigo-500 shadow-md p-4 rounded-md mb-4 flex justify-between items-center mx-4">
                    <span className="text-indigo-900 dark:text-white font-semibold">
                        You have unsaved changes
                    </span>
                    <div className="flex gap-2">
                        <button
                            className="px-3 py-1 bg-indigo-500 text-white rounded hover:bg-indigo-600"
                            onClick={async () => {
                                const changes = { title, description, hardware, prices };
                                let updated_listing = (await API.developer.editListing({ id: listing.id, ...changes })).json as Listing;
                                console.log(updated_listing)
                                setListing(  updated_listing );
                                setHasUnsavedChanges(false);
                            }}
                        >
                            Save
                        </button>
                        <button
                            className="px-3 py-1 bg-slate-300 dark:bg-slate-700 text-slate-900 dark:text-white rounded hover:bg-slate-400"
                            onClick={() => {
                                setTitle(listing.title ?? "");
                                setDescription(listing.description ?? "");
                                setHardware(listing.hardware ?? {});
                                setPrices(listing.prices ?? {});
                                setHasUnsavedChanges(false);
                            }}
                        >
                            Cancel
                        </button>
                    </div>
                </div>
            )}
            <div className="flex-1 overflow-y-auto">
                <div className="m-3">
                    <div className="flex flex-row justify-between mb-10">
                        <SelectBox options={ListingAccessModes} onSelected={onSelectedAccessMode} defaultValue={"" + listing.accessMode}></SelectBox>
                        <PromptButton
                            className="bg-red-500 hover:bg-red-600 text-white font-semibold py-2 px-4 rounded transition-colors dark:bg-red-600 dark:hover:bg-red-700 flex flex-row gap-2"
                            submitClassName={`
            flex flex-row gap-2 font-semibold py-2 px-4 rounded transition-colors
            ${isDeleteLocked
                                    ? "bg-red-300 text-gray-200 cursor-not-allowed dark:bg-red-700 dark:text-gray-400"
                                    : "bg-red-500 hover:bg-red-600 text-white dark:bg-red-600 dark:hover:bg-red-700"
                                }
          `}
                            promptContents={
                                <>
                                    <p>This operation will delete this listing. If you are sure you want to do this type "Delete" as confirmation.</p>
                                    <input onKeyUp={handleChange} type="text"
                                        className="w-full p-3 my-10 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-red-500 focus:border-red-500
          dark:bg-gray-700 dark:border-gray-600 dark:text-gray-100 dark:placeholder-gray-400 dark:focus:ring-red-400 dark:focus:border-red-400"></input>
                                </>
                            }
                            isSubmitLocked={() => isDeleteLocked}
                            onCancel={() => setText("")}
                            onSubmit={() => deleteThisListing()}
                        >
                            Delete <Delete></Delete>
                        </PromptButton>
                    </div>
                    <h1 className="text-2xl font-bold mb-4 text-gray-900 dark:text-gray-100">
                        <EditableText
                            text={title ?? "No Title"}
                            onChange={(name: string) => { setTitle(name); changeListing() }}>

                        </EditableText>
                    </h1>
                    <p className="mb-2 text-gray-700 dark:text-gray-300">
                        <EditableText
                            text={description ?? "No Description"}
                            onChange={(name: string) => { setDescription(name); changeListing() }}>

                        </EditableText>
                    </p>

                    {/* Hardware sliders */}
                    <HardwareSettings hardware={hardware ?? {}} onChange={(newHardware) => {
                        setHardware(newHardware);
                        changeListing()
                    }} />

                    {/* Prices menu */}
                    <PricesMenu prices={prices ?? {}} onChange={(newPrices) => {
                        console.log(newPrices)
                        setPrices(newPrices);
                        changeListing()
                    }} />
                </div>

            </div>
        </div>

    );
};
