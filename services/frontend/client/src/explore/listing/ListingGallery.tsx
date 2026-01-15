import { useState, useEffect } from "react";
import { Maximize2, Loader2 } from "lucide-react";
import Lightbox from "yet-another-react-lightbox";
import Zoom from "yet-another-react-lightbox/plugins/zoom";
import { ScreenshotCache } from "../../backend/utils/screenshot_cache";
import { ListingScreenshot } from "../../ListingScreenshot";
interface GalleryProps {
    listingId: string;
    screenshotIds: string[];
}

export function ListingGallery({ listingId, screenshotIds }: GalleryProps) {
    const [index, setIndex] = useState(0);
    const [isOpen, setIsOpen] = useState(false);
    const [resolvedSlides, setResolvedSlides] = useState<{ src: string }[]>([]);
    const [isResolving, setIsResolving] = useState(false);

    // Resolve all URLs for the Lightbox whenever screenshotIds change
    useEffect(() => {
        async function resolveAll() {
            if (!screenshotIds.length) return;
            setIsResolving(true);
            try {
                const urls = await Promise.all(
                    screenshotIds.map(id => ScreenshotCache.getUrl(listingId, id))
                );
                setResolvedSlides(urls.map(src => ({ src })));
            } catch (err) {
                console.error("Gallery failed to resolve slides", err);
            } finally {
                setIsResolving(false);
            }
        }
        resolveAll();
    }, [listingId, screenshotIds]);

    if (!screenshotIds?.length) {
        return (
            <div className="aspect-video rounded-3xl border-2 border-dashed border-slate-200 dark:border-slate-800 flex items-center justify-center text-slate-400 italic">
                No screenshots available
            </div>
        );
    }

    return (
        <div className="flex flex-col gap-4">
            {/* Main Display */}
            <div className="relative group aspect-video rounded-3xl overflow-hidden border border-slate-200 dark:border-slate-800 bg-slate-100 dark:bg-slate-900">
                <ListingScreenshot 
                    listingId={listingId} 
                    screenshotId={screenshotIds[index]} 
                    className="w-full h-full object-cover transition-transform duration-700 group-hover:scale-105" 
                />
                
                <div className="absolute inset-0 bg-black/0 group-hover:bg-black/20 transition-all flex items-center justify-center">
                    <button 
                        disabled={isResolving}
                        onClick={() => setIsOpen(true)} 
                        className="p-4 bg-white/90 dark:bg-slate-900/90 rounded-full scale-90 opacity-0 group-hover:opacity-100 group-hover:scale-100 transition-all shadow-2xl disabled:cursor-not-allowed"
                    >
                        {isResolving ? (
                            <Loader2 className="w-6 h-6 animate-spin text-blue-500" />
                        ) : (
                            <Maximize2 className="w-6 h-6 text-slate-900 dark:text-white" />
                        )}
                    </button>
                </div>
            </div>

            {/* Thumbnail Strip */}
            <div className="flex gap-3 overflow-x-auto pb-2 no-scrollbar">
                {screenshotIds.map((sId, i) => (
                    <button 
                        key={sId} 
                        onClick={() => setIndex(i)}
                        className={`relative flex-shrink-0 w-32 aspect-video rounded-xl overflow-hidden border-2 transition-all ${
                            index === i ? "border-blue-500 ring-4 ring-blue-500/10" : "border-transparent opacity-60 hover:opacity-100"
                        }`}
                    >
                        <ListingScreenshot listingId={listingId} screenshotId={sId} className="w-full h-full object-cover" />
                    </button>
                ))}
            </div>

            {/* Lightbox - only opens if slides are ready */}
            {resolvedSlides.length > 0 && (
                <Lightbox 
                    open={isOpen} 
                    close={() => setIsOpen(false)} 
                    index={index} 
                    slides={resolvedSlides} 
                    plugins={[Zoom]} 
                />
            )}
        </div>
    );
}