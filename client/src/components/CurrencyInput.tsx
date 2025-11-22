import { useState, useEffect } from "react";
import { DollarSign } from "lucide-react";
import type { CurrencyRequest } from "../backend";

interface CurrencyInputProps {
    value?: CurrencyRequest;
    onChange: (value: CurrencyRequest | undefined) => void;
}

const currencyPresets: CurrencyRequest[] = [
    { value: 0, name: "US Dollar", short: "USD" },
    { value: 0, name: "Euro", short: "EUR" },
    { value: 0, name: "British Pound", short: "GBP" },
];

export default function CurrencyInput({ value, onChange }: CurrencyInputProps) {
    const [enabled, setEnabled] = useState<boolean>(value !== undefined);

    // If user disables the control → send undefined upward
    useEffect(() => {
        if (!enabled) onChange(undefined);
    }, [enabled]);

    return (
        <div className="flex flex-col gap-2 mb-4 p-4 bg-slate-50 dark:bg-slate-800 rounded-lg shadow-md">
            <label className="flex items-center justify-between text-slate-700 dark:text-slate-300 mb-1">
                <span className="flex items-center gap-2">
                    <DollarSign className="w-5 h-5" />
                    Currency
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
                value={value?.short ?? ""}
                onChange={(e) => {
                    const preset = currencyPresets.find((c) => c.short === e.target.value);
                    if (preset) onChange({ ...preset, value: value?.value ?? 0 });
                }}
            >
                <option value="" disabled>Select currency…</option>
                {currencyPresets.map((c) => (
                    <option key={c.short} value={c.short}>
                        {c.name} ({c.short})
                    </option>
                ))}
            </select>

            <input
                type="number"
                disabled={!enabled}
                className="p-2 rounded-md border border-slate-300 dark:border-slate-600 
                           bg-white dark:bg-slate-700 text-slate-900 dark:text-slate-100
                           disabled:bg-slate-200 dark:disabled:bg-slate-600 disabled:text-slate-500"
                value={value?.value ?? ""}
                onChange={(e) =>
                    onChange({ ...value!, value: parseFloat(e.target.value) })
                }
            />
        </div>
    );
}
