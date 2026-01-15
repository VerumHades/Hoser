import { Cpu, MemoryStick, HardDrive } from "lucide-react";
import type { HardwareSpecification } from "../../backend/types";

interface HardwareSpecsDisplayProps {
    /** The hardware specification data to display */
    specification: HardwareSpecification;
}

/**
 * A read-only display window for hardware requirements.
 * Designed to mirror the visual style of the ListingSpecs component for UI consistency.
 */
export default function HardwareSpecsDisplay({ specification }: HardwareSpecsDisplayProps) {
    /**
     * Formats byte values into a human-readable string (e.g., 16 GB).
     */
    const formatBytes = (bytes: number): string => {
        if (bytes === 0) return "0 B";
        const units = ["B", "KB", "MB", "GB", "TB"];
        const index = Math.floor(Math.log(bytes) / Math.log(1024));
        const value = bytes / Math.pow(1024, index);
        
        return `${parseFloat(value.toFixed(1))} ${units[index]}`;
    };

    const hardwareItems = [
        {
            label: "CPU",
            value: `${specification.cpu} Cores`,
            icon: <Cpu className="w-4 h-4" />,
        },
        {
            label: "Memory",
            value: formatBytes(specification.ramBytes),
            icon: <MemoryStick className="w-4 h-4" />,
        },
        {
            label: "Storage",
            value: formatBytes(specification.diskBytes),
            icon: <HardDrive className="w-4 h-4" />,
        },
    ];

    return (
        <div className="p-6 rounded-3xl bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800">
            <h3 className="font-bold text-lg mb-4 text-slate-900 dark:text-white">Hardware Specification</h3>
            
            <ul className="space-y-4">
                {hardwareItems.map((item) => (
                    <li key={item.label} className="flex justify-between items-center">
                        <div className="flex items-center gap-2 text-slate-600 dark:text-slate-400">
                            {item.icon}
                            <span className="text-sm">{item.label}</span>
                        </div>
                        <span className="text-slate-900 dark:text-white font-mono text-xs font-semibold">
                            {item.value}
                        </span>
                    </li>
                ))}
            </ul>

            <div className="mt-6 pt-6 border-t border-slate-200 dark:border-slate-800">
                <p className="text-xs leading-relaxed text-slate-500 italic">
                    These specifications represent the minimum resource allocation required for optimal performance.
                </p>
            </div>
        </div>
    );
}