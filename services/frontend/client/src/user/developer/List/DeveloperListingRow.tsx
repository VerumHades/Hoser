
import type { DeveloperListing } from "../../../backend/repositories/developer_listing"
import TableRow from "../../../components/table/TableRow"

import { DeveloperListingItem } from "./DeveloperListingItem"


type DeveloperListingRowProps = {
    listing: DeveloperListing
    onSelect: (listing: DeveloperListing) => void
}

export function DeveloperListingRow({
    listing,
    onSelect
}: DeveloperListingRowProps) {
    return (
        <TableRow onClick={() => onSelect(listing)}>
            <DeveloperListingItem listing={listing} />
        </TableRow>
    )
}
