import { ExternalLink, Play } from "lucide-react"
import { useNavigate } from "react-router-dom"
import type { Listing } from "../../backend"
import { gotoListing } from "../../explore/PublicListingView"
import { UserLibraryItem } from "./UserLibraryItem"

type UserLibraryCardProps = {
    listing: Listing
    onCreateInstance: (listingId: string) => void
}

export function UserLibraryCard({
    listing,
    onCreateInstance
}: UserLibraryCardProps) {
    const navigate = useNavigate()

    return (
        <div className="rounded-xl border p-4 hover:shadow transition flex flex-col gap-3">
            <div
                className="cursor-pointer"
                onClick={() => gotoListing(navigate, listing)}
            >
                <UserLibraryItem listing={listing} />
            </div>

            <div className="flex justify-end gap-2">
                <button
                    className="p-2 rounded hover:bg-gray-200"
                    onClick={() => gotoListing(navigate, listing)}
                >
                    <ExternalLink size={18} />
                </button>

                <button
                    className="p-2 rounded hover:bg-gray-200"
                    onClick={() => onCreateInstance(listing.id)}
                >
                    <Play size={18} />
                </button>
            </div>
        </div>
    )
}
