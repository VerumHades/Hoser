import { type Listing } from "../../backend/repositories/listing";
import { ListingScreenshot } from "../../ListingScreenshot";
import { ChevronRight, Cpu, HardDrive } from "lucide-react";

interface PublicListingRowProps {
    listing: Listing;
    onSelect?: () => void;
    className?: string;
}

/**
 * A horizontal row display for a listing.
 * Optimized for high-density lists where scannability and quick info are key.
 */
export function PublicListingRow({ listing, onSelect, className }: PublicListingRowProps) {
    return (
        <div
            onClick={onSelect}
            className={`
                group relative flex w-full cursor-pointer items-center gap-6 
                rounded-2xl border border-slate-200 bg-white p-3 mb-4
                transition-all duration-200 hover:border-blue-400 hover:shadow-lg
                dark:border-slate-800 dark:bg-slate-900 dark:hover:border-blue-500/50
                ${className ?? ""}
            `}
        >
            <RowThumbnail listingId={listing.id} screenshotId={listing.screenshotIds?.[0]} />

            <div className="flex flex-1 flex-col justify-center min-w-0 py-1">
                <RowHeader title={listing.title || ""} author={listing.author} />
                <RowDescription description={listing.description || ""} />
                <RowMetadata listing={listing} />
            </div>

            <RowPriceAndAction priceInCents={listing.price} />
        </div>
    );
}

/**
 * Small preview image for the row layout.
 */
function RowThumbnail({ listingId, screenshotId }: { listingId: string; screenshotId?: string }) {
    return (
        <div className="h-24 w-40 flex-shrink-0 overflow-hidden rounded-xl bg-slate-100 dark:bg-slate-800">
            {screenshotId ? (
                <ListingScreenshot
                    listingId={listingId}
                    screenshotId={screenshotId}
                    className="h-full w-full object-cover transition-transform group-hover:scale-105"
                />
            ) : (
                <div className="flex h-full items-center justify-center text-[10px] italic text-slate-400">
                    No preview
                </div>
            )}
        </div>
    );
}

/**
 * Primary text info: Title and Author.
 */
function RowHeader({ title, author }: { title: string; author: string }) {
    return (
        <div className="flex items-baseline gap-2">
            <h3 className="truncate text-lg font-bold text-slate-900 dark:text-slate-100">
                {title}
            </h3>
            <span className="text-xs text-slate-400">by {author}</span>
        </div>
    );
}

/**
 * Truncated description for the row.
 */
function RowDescription({ description }: { description: string }) {
    return (
        <p className="mt-1 line-clamp-1 text-sm text-slate-600 dark:text-slate-400">
            {description}
        </p>
    );
}

/**
 * Displays technical badges or metadata tags in a row.
 */
function RowMetadata({ listing }: { listing: Listing }) {
    return (
        <div className="mt-3 flex gap-3">
            {listing.recommended_hardware && (
                <div className="flex items-center gap-1.5 rounded-md bg-slate-50 px-2 py-1 text-[10px] font-bold uppercase tracking-wider text-slate-500 dark:bg-slate-800/50">
                    <Cpu className="h-3 w-3" /> Optimized
                </div>
            )}
            <div className="flex items-center gap-1.5 rounded-md bg-blue-50 px-2 py-1 text-[10px] font-bold uppercase tracking-wider text-blue-600 dark:bg-blue-900/20">
                <HardDrive className="h-3 w-3" /> V{listing.id.slice(0, 3)} {/* Placeholder for Version */}
            </div>
        </div>
    );
}

/**
 * Right-side price display and navigation hint.
 */
function RowPriceAndAction({ priceInCents }: { priceInCents: number }) {
    const formattedPrice = new Intl.NumberFormat("de-DE", {
        style: "currency",
        currency: "EUR",
    }).format(priceInCents / 100);

    return (
        <div className="flex flex-col items-end gap-1 px-4">
            <span className="text-xl font-black text-slate-900 dark:text-white">
                {formattedPrice}
            </span>
            <div className="flex items-center text-xs font-semibold text-blue-600 opacity-0 transition-opacity group-hover:opacity-100">
                View Details <ChevronRight className="h-4 w-4" />
            </div>
        </div>
    );
}