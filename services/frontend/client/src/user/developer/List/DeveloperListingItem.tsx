import { Earth, EarthLock } from "lucide-react"
import type { DeveloperListing } from "../../../backend/repositories/developer_listing"


type DeveloperListingItemProps = {
    listing: DeveloperListing
}

export function DeveloperListingItem({
    listing
}: DeveloperListingItemProps) {
    return (
        [
            <h3 className="font-semibold text-lg">
                {listing.title}
            </h3>,

            <p className="text-sm text-gray-600">
                {listing.description}
            </p>,

            <div className="flex items-center gap-2 text-sm">
                {listing.accessMode === 1 ? (
                    <>
                        <Earth className="w-4 h-4" />
                        <span>Public</span>
                    </>
                ) : (
                    <>
                        <EarthLock className="w-4 h-4" />
                        <span>Private</span>
                    </>
                )}
            </div>
        ]    
    )
}
