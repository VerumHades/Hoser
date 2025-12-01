import type { CurrencyRequest, Listing, PricingEntry } from "../../backend";
import CurrencyInput from "../pricing/CurrencyInput";

// --- Prices Menu Component ---
interface PricesMenuProps {
    listing: Listing;
    onChange: (id: string, pricing: CurrencyRequest | null) => void;
}

export default function PricesMenu({ listing, onChange }: PricesMenuProps) {
    return (
        <div className="flex flex-col gap-4">
            {
                 listing.prices?.map?.((x: PricingEntry) => 
                    x.currency == undefined ? <></> :
                    <CurrencyInput
                        value={x.currency}
                        text=""
                        onChange={(v) => {
                            onChange(x.id, v)
                        }}
                    />
                )
            }
        </div>
    );
};
