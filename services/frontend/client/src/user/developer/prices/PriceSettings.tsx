import React, { useState, useEffect } from "react";
import { Euro } from "lucide-react";

interface PriceSettingsProps {
    /** Price in cents (e.g., 100 = 1.00€) */
    priceCents: number;
    onChange: (cents: number) => void;
}

export default function PriceSettings({ priceCents, onChange }: PriceSettingsProps) {
    // We use a string for the internal state to allow users to type decimals comfortably
    const [displayValue, setDisplayValue] = useState((priceCents / 100).toString());

    // Sync internal state if the parent price changes (e.g., on "Reset" or Undo)
    useEffect(() => {
        const decimalValue = (priceCents / 100).toString();
        if (parseFloat(displayValue) !== priceCents / 100) {
            setDisplayValue(decimalValue);
        }
    }, [priceCents]);

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const val = e.target.value;
        // Allow numeric input and single decimal point/comma
        if (/^\d*[.,]?\d{0,2}$/.test(val) || val === "") {
            setDisplayValue(val.replace(",", "."));
        }
    };

    const handleBlur = () => {
        const parsed = parseFloat(displayValue);
        if (isNaN(parsed)) {
            // Revert to 0 if input is invalid
            setDisplayValue("0");
            onChange(0);
        } else {
            // Convert to cents (integer) and round to prevent floating point math errors
            const cents = Math.round(parsed * 100);
            onChange(cents);
            // Re-format display value to fixed 2 decimals (e.g., "1" becomes "1.00")
            setDisplayValue((cents / 100).toFixed(2));
        }
    };

    return (
        <div className="flex flex-col gap-1.5 max-w-[240px]">
            <label className="text-sm font-semibold text-slate-700 dark:text-slate-300">
                Price (EUR)
            </label>
            
            <div className="relative group">
                <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                    <Euro className="h-4 w-4 text-slate-400 group-focus-within:text-blue-500 transition-colors" />
                </div>
                
                <input
                    type="text" // Using text instead of number for better control over decimals/formatting
                    inputMode="decimal"
                    value={displayValue}
                    onChange={handleChange}
                    onBlur={handleBlur}
                    placeholder="0.00"
                    className="w-full pl-9 pr-4 py-2 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl font-medium focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 outline-none transition-all"
                />
            </div>
        </div>
    );
}