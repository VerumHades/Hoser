
import { Delete } from "lucide-react";
import type { ListingAction } from "../DeveloperListing";
import TableRow from "../../../components/table/TableRow";
import { BillingFrequencyName, type Listing, type PricingEntry } from "../../../backend";
import Table from "../../../components/table/Table";


// --- Prices Menu Component ---
interface PricesMenuProps {
    listing: Listing;
    onAction:  (action:  ListingAction) => void,
}

export default function PricesMenu({ listing, onAction}: PricesMenuProps) {
    return (
        <Table>
            {
                 listing.prices?.map?.((x: PricingEntry) => 
                    <TableRow>
                        <label>{x.currency.value}</label>
                        <label>{`${x.currency.name} ${x.currency.short}`}</label>
                        <label>{BillingFrequencyName[x.type]}</label>
                        <Delete onClick={() => onAction({ type: "deletePrice", priceId: x.id })}></Delete>
                    </TableRow>
                )
            }
        </Table>
    );
};
