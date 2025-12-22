import { FlowSwitch } from "../../components/navigation/Flow"
import type { DeveloperListing } from "../../backend"
import { DeveloperListingItem } from "./DeveloperListingItem"

type DeveloperListingCardProps = {
    listing: DeveloperListing
    onSelect: (listing: DeveloperListing) => void
}

export function DeveloperListingCard({
    listing,
    onSelect
}: DeveloperListingCardProps) {
    return (
        <div onClick={() => onSelect(listing)} className="rounded-xl border p-4 hover:shadow transition">
            <DeveloperListingItem listing={listing} />
        </div>
    )
}
