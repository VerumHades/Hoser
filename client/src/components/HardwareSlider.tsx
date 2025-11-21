// --- Hardware Slider Component ---
interface HardwareSliderProps {
    label: string;
    icon: React.ReactNode;
    min: number;
    max: number;
    value?: number;
    onChange: (value: number) => void;
    unit?: string;
}

export default function HardwareSlider({
    label,
    icon,
    min,
    max,
    value,
    onChange,
    unit,
}: HardwareSliderProps) {
    return (
        <div className="flex flex-col mb-4">
            <div className="flex items-center justify-between mb-1">
                <div className="flex items-center gap-2 text-slate-700 dark:text-slate-300">
                    {icon}
                    <span className="font-medium">{label}</span>
                </div>
                <span className="text-slate-900 dark:text-slate-100">
                    {value} {unit}
                </span>
            </div>
            <input
                type="range"
                min={min}
                max={max}
                value={value}
                onChange={(e) => onChange(Number(e.target.value))}
                className="w-full h-2 bg-slate-300 rounded-lg appearance-none cursor-pointer dark:bg-slate-700 accent-indigo-500"
            />
        </div>
    );
};
