import { ExternalLink, Play } from "lucide-react"
import { useNavigate } from "react-router-dom"
import TableRow from "../../../components/table/TableRow"

import { gotoListing } from "../../../explore/PublicListingView"
import { UserLibraryItem } from "./UserLibraryItem"
import type { Listing } from "../../../backend/repositories/listing"

type UserLibraryRowProps = {
    listing: Listing
    onCreateInstance: (listingId: string) => void
}

export function UserLibraryRow({
    listing,
    onCreateInstance
}: UserLibraryRowProps) {
    const navigate = useNavigate()

    return (
        <TableRow className="flex justify-between items-center px-4 py-2 hover:bg-gray-50 rounded">
            <div
                className="flex-1 cursor-pointer"
                onClick={() => gotoListing(navigate, listing)}
            >
                <UserLibraryItem listing={listing} />
            </div>

            <div className="flex space-x-2">
                <button
                    className="p-2 rounded hover:bg-gray-200"
                    onClick={() => gotoListing(navigate, listing)}
                >
                    <ExternalLink size={20} />
                </button>

                <button
                    className="p-2 rounded hover:bg-gray-200"
                    onClick={() => onCreateInstance(listing.id)}
                >
                    <Play size={20} />
                </button>
            </div>
        </TableRow>
    )
}
