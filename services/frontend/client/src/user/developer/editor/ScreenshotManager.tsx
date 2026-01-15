import { useEffect, useState } from "react";
import { DeveloperListingAPI, type DeveloperListing } from "../../../backend/repositories/developer_listing";
import toast from "react-hot-toast";
import Section from "../../../templates/components/Section";
import { ScreenshotCache } from "../../../backend/utils/screenshot_cache";

interface Props {
    listingId: string;
    screenshotIds: string[];
    onUpdate: (listing: DeveloperListing) => void;
    onRemove: (id: string) => void;
}

export function ScreenshotManager({ listingId, screenshotIds, onUpdate, onRemove }: Props) {
    const [isUploading, setIsUploading] = useState(false);
    const [urls, setUrls] = useState<Record<string, string>>({});

    useEffect(() => {
        screenshotIds.forEach(async (id) => {
            if (!urls[id]) {
                const url = await ScreenshotCache.getUrl(listingId, id);
                setUrls(prev => ({ ...prev, [id]: url }));
            }
        });
    }, [screenshotIds, listingId]);

    const handleUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (!file) return;
        setIsUploading(true);
        const tid = toast.loading("Uploading...");
        try {
            const { uploadUrl } = await DeveloperListingAPI.screenshots.createUploadUrl(listingId);
            await fetch(uploadUrl, { method: "PUT", body: file, headers: { "Content-Type": file.type } });
            const updated = await DeveloperListingAPI.get(listingId);
            onUpdate(updated);
            toast.success("Uploaded", { id: tid });
        } catch (err) {
            toast.error("Failed", { id: tid });
        } finally { setIsUploading(false); }
    };

    return (
        <Section title="Screenshots" description="Visuals for your listing (Max 5).">
            <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
                {screenshotIds.map(id => (
                    <div key={id} className="relative group aspect-video bg-gray-200 dark:bg-gray-700 rounded-lg overflow-hidden border">
                        {urls[id] && <img src={urls[id]} className="w-full h-full object-cover" />}
                        <button 
                            onClick={() => onRemove(id)}
                            className="absolute top-2 right-2 bg-red-600 p-1.5 rounded-full opacity-0 group-hover:opacity-100 transition-opacity"
                        >
                            <svg className="w-4 h-4 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path d="M6 18L18 6M6 6l12 12" /></svg>
                        </button>
                    </div>
                ))}
                {screenshotIds.length < 5 && (
                    <label className="flex flex-col items-center justify-center aspect-video border-2 border-dashed rounded-lg cursor-pointer hover:bg-gray-50">
                        {isUploading ? <div className="animate-spin h-6 w-6 border-b-2 border-indigo-500" /> : <p className="text-xs">Add Screenshot</p>}
                        <input type="file" className="hidden" onChange={handleUpload} disabled={isUploading} accept="image/*" />
                    </label>
                )}
            </div>
        </Section>
    );
}