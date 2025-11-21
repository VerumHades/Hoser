import type { PricesRequest } from "../../backend";
import CurrencyInput from "../CurrencyInput";

// --- Prices Menu Component ---
interface PricesMenuProps {
    prices: PricesRequest;
    onChange: (prices: PricesRequest) => void;
}

export default function PricesMenu({ prices, onChange }: PricesMenuProps) {
    return (
        <div className="flex flex-col gap-4">
            <CurrencyInput
                value={prices.singlePurchase}
                onChange={(v) => onChange({ ...prices, singlePurchase: v })}
            />
            <CurrencyInput
                value={prices.monthlySubscription}
                onChange={(v) => onChange({ ...prices, monthlySubscription: v })}
            />
        </div>
    );
};
