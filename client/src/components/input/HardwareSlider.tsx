import React, { useState, useEffect } from "react";

interface HardwareUnit {
    label: string;       // e.g., "B", "KB", "MB", "GB"
    multiplier: number;  // multiplier to convert to base value
}

interface HardwareSliderProps {
    label: string;
    icon: React.ReactNode;
    min: number;       // min in base units
    max: number;       // max in base units
    step?: number;     // step in base units
    value?: number;    // value in base units
    onChange: (value: number) => void;
    units?: HardwareUnit[];  // optional array of units
}

export default function HardwareSlider({
    label,
    icon,
    min,
    max,
    step = 1,
    value,
    onChange,
    units,
}: HardwareSliderProps) {
    const [selectedUnit, setSelectedUnit] = useState<HardwareUnit | undefined>(units?.[units.length - 1]);

    // Automatically select the most appropriate unit when value changes
    useEffect(() => {
        if (!units || value === undefined) return;

        // find largest unit where value/unit >= 1
        const unit = [...units].reverse().find((u) => value >= u.multiplier) || units[0];
        setSelectedUnit(unit);
    }, [value, units]);

    const minInUnit = selectedUnit ? min / selectedUnit.multiplier : min;
    const maxInUnit = selectedUnit ? max / selectedUnit.multiplier : max;
    const stepInUnit = selectedUnit ? step / selectedUnit.multiplier : step;

    const displayValue = selectedUnit ? (value ?? min) / selectedUnit.multiplier : value ?? min;

    const clamp = (val: number) => Math.min(Math.max(val, minInUnit), maxInUnit);

    const handleValueChange = (val: number) => {
        const clamped = clamp(val);
        if (selectedUnit) onChange(clamped * selectedUnit.multiplier);
        else onChange(clamped);
    };

    const handleUnitChange = (unitLabel: string) => {
        const unit = units?.find((u) => u.label === unitLabel);
        if (unit) setSelectedUnit(unit);
    };

    return (
        <div className="flex flex-col mb-4">
            <div className="flex items-center justify-between mb-1">
                <div className="flex items-center gap-2 text-slate-700 dark:text-slate-300">
                    {icon}
                    <span className="font-medium">{label}</span>
                </div>
            </div>
            <div className="flex items-center gap-2">
                <input
                    type="range"
                    min={minInUnit}
                    max={maxInUnit}
                    step={stepInUnit}
                    value={displayValue}
                    onChange={(e) => handleValueChange(Number(e.target.value))}
                    className="flex-1 h-2 bg-slate-300 rounded-lg appearance-none cursor-pointer dark:bg-slate-700 accent-indigo-500"
                />
                <input
                    type="number"
                    min={minInUnit}
                    max={maxInUnit}
                    step={stepInUnit}
                    value={displayValue}
                    onChange={(e) => handleValueChange(Number(e.target.value))}
                    className="w-20 p-1 border rounded-md dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                />
                {units && (
                    <select
                        value={selectedUnit?.label}
                        onChange={(e) => handleUnitChange(e.target.value)}
                        className="p-1 border rounded-md dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                    >
                        {units.map((u) => (
                            <option key={u.label} value={u.label}>
                                {u.label}
                            </option>
                        ))}
                    </select>
                )}
            </div>
        </div>
    );
}
