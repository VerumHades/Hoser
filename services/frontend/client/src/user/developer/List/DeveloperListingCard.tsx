import type { DeveloperListing } from "../../../backend"
import { DeveloperListingItem } from "./DeveloperListingItem"

type DeveloperListingCardProps = {
    listing: DeveloperListing
    onSelect: (listing: DeveloperListing) => void
}

export function DeveloperListingCard({ listing, onSelect }: DeveloperListingCardProps) {
    return (
        <div
            onClick={() => onSelect(listing)}
            className={`
                cursor-pointer
                rounded-xl
                border
                border-slate-200
                bg-white
                p-4
                shadow-sm
                transition
                duration-150
                hover:shadow-md
                hover:border-slate-300
                dark:border-slate-800
                dark:bg-slate-900
                dark:hover:border-slate-600
                focus:outline-none
                focus:ring-2
                focus:ring-blue-500
                focus:ring-offset-2
            `}
            tabIndex={0} // makes card focusable for keyboard
        >
            <DeveloperListingItem listing={listing} />
        </div>
    )
}
