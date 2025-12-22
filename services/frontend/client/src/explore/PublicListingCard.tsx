import React from "react";
import { PriceTag } from "../user/developer/prices/PriceTag";
import { type Listing } from "../backend";

interface PublicListingCardProps {
    listing: Listing;
    onSelect?: () => void;
}

export function PublicListingCard({ listing, onSelect }: PublicListingCardProps) {
    return (
        <div
            onClick={onSelect}
            className="rounded-xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 p-4 hover:border-slate-400 dark:hover:border-slate-600 hover:shadow-md transition cursor-pointer flex flex-col gap-2"
        >
            <h3 className="text-lg font-semibold text-slate-900 dark:text-slate-200 truncate">
                {listing.title}
            </h3>
            <p className="text-xs text-slate-600 dark:text-slate-500">
                By {listing.author}
            </p>
            <p className="text-sm text-slate-700 dark:text-slate-400 line-clamp-3">
                {listing.description}
            </p>
            {listing.price && (
                <div className="flex items-center justify-end mt-2">
                    <PriceTag prices={[listing.price]} rate="one time" />
                </div>
            )}
        </div>
    );
}
