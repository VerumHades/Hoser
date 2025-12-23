import type { Listing } from "../../../backend"

type UserLibraryItemProps = {
    listing: Listing
}

export function UserLibraryItem({
    listing
}: UserLibraryItemProps) {
    return (
        <>
            <h3 className="font-semibold text-lg">
                {listing.title}
            </h3>

            <p className="text-sm text-gray-600">
                {listing.description}
            </p>
        </>
    )
}
