import { AnimatePresence, motion } from "framer-motion";
import HardwareSettings from "./hardware/HardwareSettings";
import PriceSettings from "./prices/PriceSettings";
import DeleteListingPrompt from "./DeleteListingPrompt";

import { DeveloperListingAPI, type DeveloperListing } from "../../backend/repositories/developer_listing";
import { Button } from "../../templates/components/Button";
import Section from "../../templates/components/Section";
import EditableText from "../../templates/components/EditableText";
import CustomSelect from "../../templates/components/SelectBox";
import { ScreenshotManager } from "./editor/ScreenshotManager";
import { GithubSetupForm } from "./editor/GithubSetupForm";
import { useListingEditor } from "./editor/useListingEditor";
import { Toaster } from "react-hot-toast";
import MarkdownEditor from "../../components/input/MarkdownEditor";
import type { ListingAccessMode } from "../../backend/types";

interface DeveloperListingDisplayProps {
    listing: DeveloperListing;
    onShouldClose?: () => void;
}

/**
 * Primary editor component for managing Developer Listing details.
 * Implements a two-column layout for configuration and hardware specifications.
 */
export default function DeveloperListingEditor({ listing: sourceListing, onShouldClose }: DeveloperListingDisplayProps) {
    const { 
        listing, 
        dispatch, 
        hasUnsavedChanges, 
        markChanged, 
        resetChanges, 
        saveChanges 
    } = useListingEditor(sourceListing);

    const deleteListing = async () => {
        await DeveloperListingAPI.delete(listing.id);
        onShouldClose?.();
    };

    const handleDocumentationChange = (markdown: string) => {
        dispatch({ type: "setDocumentationMarkdown", markdown });
        markChanged();
    };

    return (
        <div className="relative flex flex-col w-full h-full items-center overflow-y-auto p-6">
            <Toaster position="top-right" />
            
            <AnimatePresence>
                {hasUnsavedChanges && (
                    <UnsavedChangesBanner onSave={saveChanges} onReset={resetChanges} />
                )}
            </AnimatePresence>

            <div className="w-full max-w-6xl space-y-6">
                <Section title="General">
                    <EditableText 
                        label="Title" 
                        text={listing.title || ""} 
                        onChange={title => { dispatch({ type: "setTitle", title }); markChanged(); }} 
                        block 
                    />
                    <EditableText 
                        label="Description" 
                        text={listing.description || ""} 
                        onChange={description => { dispatch({ type: "setDescription", description }); markChanged(); }} 
                        block 
                    />
                </Section>
                <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 items-start">
                    <div className="space-y-6">
                        <Section title="Visibility">
                            <CustomSelect 
                                value={listing.accessMode as number} 
                                onChange={mode => { dispatch({ type: "setAccessMode", mode: mode as ListingAccessMode }); markChanged(); }} 
                                options={[{ label: "Public", value: 1 }, { label: "Private", value: 0 }]} 
                                block 
                            />
                        </Section>
                    
                        <PriceSettings 
                            priceCents={listing.price} 
                            onChange={currency => { dispatch({ type: "setPrice", currency }); markChanged(); }} 
                        />
                    </div>

                    <div className="h-full">
                        <Section title="Hardware Specification">
                            <HardwareSettings 
                                initialSpec={listing.recommended_hardware} 
                                onChange={hardware => { dispatch({ type: "setHardware", hardware }); markChanged(); }} 
                            />
                        </Section>
                        
                    </div>
                </div>

                <Section title="Documentation">
                    <MarkdownEditor
                        value={listing.documentation_markdown} onChange={handleDocumentationChange}>
                    </MarkdownEditor>
                </Section>

                <ScreenshotManager 
                    listingId={listing.id} 
                    screenshotIds={listing.screenshotIds || []} 
                    onUpdate={updatedListing => dispatch({ type: "set", listing: updatedListing })}
                    onRemove={id => { dispatch({ type: "removeScreenshot", id }); markChanged(); }}
                />

                <GithubSetupForm listingId={listing.id} />

                <Section title="Danger Zone">
                    <DeleteListingPrompt onDelete={deleteListing} />
                </Section>
            </div>
        </div>
    );
}

/**
 * Notification banner displayed when the local state diverges from the persisted server state.
 */
function UnsavedChangesBanner({ onSave, onReset }: { onSave: () => void, onReset: () => void }) {
    return (
        <motion.div 
            initial={{ y: 50, opacity: 0 }} 
            animate={{ y: 0, opacity: 1 }} 
            exit={{ y: 50, opacity: 0 }} 
            className="fixed bottom-6 z-10 bg-indigo-100 border-l-4 border-indigo-500 p-4 rounded-md flex justify-between items-center w-[90%] max-w-3xl shadow-2xl"
        >
            <span className="font-semibold text-indigo-900">You have unsaved changes</span>
            <div className="flex gap-3">
                <Button variant="primary" onClick={onSave}>Save Changes</Button>
                <Button variant="secondary" onClick={onReset}>Discard</Button>
            </div>
        </motion.div>
    );
}