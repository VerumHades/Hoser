import { DollarSign } from "lucide-react";
import type { CurrencyRequest } from "../backend";

interface CurrencyInputProps {
    value?: CurrencyRequest;
    onChange: (value: CurrencyRequest) => void;
}

const currencyPresets: CurrencyRequest[] = [
    { value: 0, name: "US Dollar", short: "USD" },
    { value: 0, name: "Euro", short: "EUR" },
    { value: 0, name: "British Pound", short: "GBP" },
];

export default function CurrencyInput({
    value,
    onChange,
}: CurrencyInputProps) {
    return (
        <div className="flex flex-col gap-2 mb-4 p-4 bg-slate-50 dark:bg-slate-800 rounded-lg shadow-md">
            <label className="flex items-center gap-2 text-slate-700 dark:text-slate-300">
                <DollarSign className="w-5 h-5" />
                Currency
            </label>
            <select
                className="p-2 rounded-md border border-slate-300 dark:border-slate-600 bg-white dark:bg-slate-700 text-slate-900 dark:text-slate-100"
                value={value?.short}
                onChange={(e) => {
                    const preset = currencyPresets.find((c) => c.short === e.target.value);
                    if (preset) onChange({ ...preset, value: value?.value ?? 0 });
                }}
            >
                {currencyPresets.map((c) => (
                    <option key={c.short} value={c.short}>
                        {c.name} ({c.short})
                    </option>
                ))}
            </select>
            <input
                type="number"
                className="p-2 rounded-md border border-slate-300 dark:border-slate-600 bg-white dark:bg-slate-700 text-slate-900 dark:text-slate-100"
                value={value?.value ?? 0}
                onChange={(e) =>
                    onChange({ ...value!, value: parseFloat(e.target.value) })
                }
            />
        </div>
    );
};