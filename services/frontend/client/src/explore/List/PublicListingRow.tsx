import React from "react"
import type { Listing } from "../../backend/repositories/listing"
interface PublicListingRowProps {
    listing: Listing
    onSelect?: () => void
    className?: string
}

export function PublicListingRow({
    listing,
    onSelect,
    className,
}: PublicListingRowProps) {
    return (
        <div
            onClick={onSelect}
            className={`
                cursor-pointer
                rounded-xl
                border
                border-slate-200
                bg-white
                p-4
                mb-5
                transition
                hover:border-slate-400
                hover:shadow-md
                dark:border-slate-800
                dark:bg-slate-900
                dark:hover:border-slate-600
                ${className ?? ""}
            `}
        >
            <div
                className="
                    flex
                    min-w-0
                    flex-col
                    justify-between
                    gap-2
                "
            >
                <h3
                    className="
                        truncate
                        text-lg
                        font-semibold
                        text-slate-900
                        dark:text-slate-200
                    "
                >
                    {listing.title}
                </h3>

                <p
                    className="
                        text-xs
                        text-slate-600
                        dark:text-slate-500
                    "
                >
                    By {listing.author}
                </p>
            </div>

            <p
                className="
                    mt-2
                    line-clamp-2
                    text-sm
                    text-slate-700
                    dark:text-slate-400
                "
            >
                {listing.description}
            </p>

            {listing.price && (
                <div
                    className="
                        mt-2
                        flex
                        items-center
                        justify-end
                    "
                >
                </div>
            )}
        </div>
    )
}
