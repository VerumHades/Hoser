import { useState, useEffect } from "react";
import { DollarSign } from "lucide-react";
import type { Money } from "../../../backend";

interface CurrencyInputProps {
    value?: Money;
    text?: string,
    onChange: (value: Money | null) => void;
}

const currencyPresets: Money[] = [
    { amount: 0, name: "US Dollar", code: "USD" },
    { amount: 0, name: "Euro", code: "EUR" },
    { amount: 0, name: "British Pound", code: "GBP" },
];

export default function CurrencyInput({ value, onChange, text }: CurrencyInputProps) {
    const [enabled, setEnabled] = useState<boolean>(value !== undefined);

    // If user disables the control → send undefined upward
    useEffect(() => {
        if (!enabled) onChange(null);
    }, [enabled]);

    return (
        <div className="flex flex-col gap-2 mb-4 p-">
            <label className="flex items-center justify-between text-slate-700 dark:text-slate-300 mb-1">
                <span className="flex items-center gap-2">
                    <DollarSign className="w-5 h-5" />
                    {text ?? "Currency"}
                </span>

                {/* Enable/Disable checkbox */}
                <label className="flex items-center gap-1 cursor-pointer">
                    <input
                        type="checkbox"
                        checked={enabled}
                        onChange={(e) => setEnabled(e.target.checked)}
                    />
                    <span className="text-sm">Enable</span>
                </label>
            </label>

            <select
                disabled={!enabled}
                className="p-2 rounded-md border border-slate-300 dark:border-slate-600 
                           bg-white dark:bg-slate-700 text-slate-900 dark:text-slate-100
                           disabled:bg-slate-200 dark:disabled:bg-slate-600 disabled:text-slate-500"
                value={value?.code ?? ""}
                onChange={(e) => {
                    const preset = currencyPresets.find((c) => c.code === e.target.value);
                    if (preset) onChange({ ...preset, amount: value?.amount ?? 0 });
                }}
            >
                <option value="" disabled>Select currency…</option>
                {currencyPresets.map((c) => (
                    <option key={c.code} value={c.code}>
                        {c.name} ({c.code})
                    </option>
                ))}
            </select>

            <input
                type="number"
                disabled={!enabled}
                className="p-2 rounded-md border border-slate-300 dark:border-slate-600 
                           bg-white dark:bg-slate-700 text-slate-900 dark:text-slate-100
                           disabled:bg-slate-200 dark:disabled:bg-slate-600 disabled:text-slate-500"
                value={value?.amount ?? ""}
                onChange={(e) =>
                    onChange({ ...value!, amount: parseFloat(e.target.value) })
                }
            />
        </div>
    );
}
