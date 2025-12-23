import { ExternalLink, Play } from "lucide-react"
import { useNavigate } from "react-router-dom"
import type { Listing } from "../../../backend"
import { gotoListing } from "../../../explore/PublicListingView"
import { UserLibraryItem } from "./UserLibraryItem"

type UserLibraryCardProps = {
    listing: Listing
    onCreateInstance: (listingId: string) => void
}

export function UserLibraryCard({ listing, onCreateInstance }: UserLibraryCardProps) {
    const navigate = useNavigate()

    return (
        <div
            className="
                flex
                flex-col
                gap-3
                rounded-xl
                border
                border-slate-200
                bg-white
                p-4
                shadow-sm
                transition
                duration-150
                hover:shadow-md
                dark:border-slate-800
                dark:bg-slate-900
                dark:hover:border-slate-600
            "
        >
            <div
                className="cursor-pointer focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 rounded"
                onClick={() => gotoListing(navigate, listing)}
                tabIndex={0} // makes card clickable via keyboard
            >
                <UserLibraryItem listing={listing} />
            </div>

            <div className="flex justify-end gap-2">
                <button
                    onClick={() => gotoListing(navigate, listing)}
                    className="
                        p-2
                        rounded-md
                        hover:bg-gray-100
                        dark:hover:bg-slate-700
                        transition
                        duration-150
                        focus:outline-none
                        focus:ring-2
                        focus:ring-blue-500
                        focus:ring-offset-2
                    "
                    aria-label="View Listing"
                >
                    <ExternalLink size={18} />
                </button>

                <button
                    onClick={() => onCreateInstance(listing.id)}
                    className="
                        p-2
                        rounded-md
                        hover:bg-green-100
                        dark:hover:bg-green-900
                        transition
                        duration-150
                        focus:outline-none
                        focus:ring-2
                        focus:ring-green-500
                        focus:ring-offset-2
                    "
                    aria-label="Create Instance"
                >
                    <Play size={18} />
                </button>
            </div>
        </div>
    )
}
