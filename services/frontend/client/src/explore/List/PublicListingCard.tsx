import React from "react";
import { ListingScreenshot } from "../../ListingScreenshot";
import { Bookmark, ShoppingCart, Cpu } from "lucide-react";
import type { Listing } from "../../backend/repositories/listing";

interface PublicListingCardProps {
    listing: Listing;
    onSelect?: () => void;
}

/**
 * A highly visual card representing a marketplace listing.
 * Designed to entice users through screenshot previews and clear pricing.
 */
export function PublicListingCard({ listing, onSelect }: PublicListingCardProps) {
    return (
        <div
            onClick={onSelect}
            className="group flex flex-col overflow-hidden rounded-2xl border border-slate-200 bg-white transition-all duration-300 hover:-translate-y-1 hover:border-blue-400 hover:shadow-xl dark:border-slate-800 dark:bg-slate-900 cursor-pointer"
        >
            <CardThumbnail listingId={listing.id} screenshotId={listing.screenshotIds?.[0]} />
            
            <div className="flex flex-col gap-3 p-5">
                <CardHeader title={listing.title} author={listing.author} />
                <CardDescription description={listing.description} />
                
                <div className="mt-2 flex items-center justify-between border-t border-slate-100 pt-4 dark:border-slate-800">
                    <CardPrice priceInCents={listing.price} />
                    <CardActionIcons hasHardwareSpecs={!!listing.recommended_hardware} />
                </div>
            </div>
        </div>
    );
}

/**
 * Renders the primary visual preview for the listing.
 */
function CardThumbnail({ listingId, screenshotId }: { listingId: string; screenshotId?: string }) {
    return (
        <div className="relative aspect-video w-full overflow-hidden bg-slate-100 dark:bg-slate-800">
            {screenshotId ? (
                <ListingScreenshot
                    listingId={listingId}
                    screenshotId={screenshotId}
                    className="h-full w-full object-cover transition-transform duration-500 group-hover:scale-110"
                />
            ) : (
                <div className="flex h-full items-center justify-center text-slate-400">
                    <span className="text-xs italic">No preview available</span>
                </div>
            )}
            <div className="absolute inset-0 bg-gradient-to-t from-black/20 to-transparent opacity-0 transition-opacity group-hover:opacity-100" />
        </div>
    );
}

/**
 * Renders the title and author information.
 */
function CardHeader({ title, author }: { title: string; author: string }) {
    return (
        <div className="flex flex-col">
            <h3 className="line-clamp-1 text-lg font-bold text-slate-900 transition-colors group-hover:text-blue-600 dark:text-slate-100">
                {title}
            </h3>
            <span className="text-xs font-medium text-slate-500">
                by <span className="text-blue-500/80">{author}</span>
            </span>
        </div>
    );
}

/**
 * Renders the short description text.
 */
function CardDescription({ description }: { description: string }) {
    return (
        <p className="line-clamp-2 text-sm leading-relaxed text-slate-600 dark:text-slate-400">
            {description}
        </p>
    );
}

/**
 * Formats and renders the listing price.
 */
function CardPrice({ priceInCents }: { priceInCents: number }) {
    const formattedPrice = new Intl.NumberFormat("de-DE", {
        style: "currency",
        currency: "EUR",
    }).format(priceInCents / 100);

    return (
        <span className="text-xl font-black text-slate-900 dark:text-white">
            {formattedPrice}
        </span>
    );
}

/**
 * Renders visual cues for features like hardware requirements or save status.
 */
function CardActionIcons({ hasHardwareSpecs }: { hasHardwareSpecs: boolean }) {
    return (
        <div className="flex gap-2 text-slate-400">
            {hasHardwareSpecs && (
                <Cpu className="h-4 w-4 transition-colors hover:text-slate-600" title="Specs included" />
            )}
            <Bookmark className="h-4 w-4 transition-colors hover:text-blue-500" />
            <ShoppingCart className="h-4 w-4 transition-colors hover:text-blue-500" />
        </div>
    );
}