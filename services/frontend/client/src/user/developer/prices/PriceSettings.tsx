import React, { useRef } from "react";
import SelectBox from "../../../components/input/SelectBox";
import type { Money } from "../../../backend";
interface PriceSettingsProps {
    price?: Money;
    onChange: (price: Money) => void;
}

export default function PriceSettings({ price, onChange }: PriceSettingsProps) {
    const [internalAmount, setInternalAmount] = React.useState(price?.amount ?? 0);

    React.useEffect(() => {
        setInternalAmount(price?.amount ?? 0); // sync if parent changes externally
    }, [price?.amount]);

    const handleBlur = () => {
        onChange({ amount: internalAmount, code: price?.code ?? "EUR" });
    };

    return (
        <div className="flex flex-col gap-2">
            <SelectBox
                options={{ EUR: { label: "Euro" }, USD: { label: "US Dollar" } }}
                value={price?.code ?? "EUR"}
                onSelected={(currency) => onChange({ amount: internalAmount, code: currency })}
            />
            <input
                type="number"
                value={internalAmount}
                onChange={(e) => setInternalAmount(Number(e.target.value))}
                onBlur={handleBlur}
            />
        </div>
    );
}
