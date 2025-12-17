import { Cpu, MemoryStick, HardDrive } from "lucide-react";
import HardwareSlider from "../../../components/input/HardwareSlider";
import type { HardwareUpdate } from "../../../backend";

interface HardwareSettingsProps {
    hardware: HardwareUpdate;
    onChange: (update: HardwareUpdate) => void;
}

export default function HardwareSettings({
    hardware,
    onChange,
}: HardwareSettingsProps) {
    const GB = 1024 ** 3;

    const byteUnits = [
        { label: "B", multiplier: 1 },
        { label: "KB", multiplier: 1024 },
        { label: "MB", multiplier: 1024 ** 2 },
        { label: "GB", multiplier: 1024 ** 3 },
        { label: "TB", multiplier: 1024 ** 4 },
    ];

    return (
        <div className="p-4">
            {/* CPU cores */}
            <HardwareSlider
                label="CPU Cores"
                icon={<Cpu className="w-5 h-5" />}
                min={1}
                max={64}
                value={hardware.cpu ?? 1}
                onChange={(v) => onChange({ ...hardware, cpu: v })}
            />

            {/* RAM */}
            <HardwareSlider
                label="RAM"
                icon={<MemoryStick className="w-5 h-5" />}
                min={1024}               // 1 GB in bytes
                max={1024 * GB}            // 1024 GB in bytes
                value={hardware.ramBytes ?? 8 * GB}
                onChange={(v) => onChange({ ...hardware, ramBytes: Math.floor(v) })}
                units={byteUnits}
            />

            {/* Disk */}
            <HardwareSlider
                label="Disk"
                icon={<HardDrive className="w-5 h-5" />}
                min={1024}             // 128 GB in bytes
                max={1024 * GB}            // 1024 GB in bytes
                value={hardware.diskBytes ?? 256 * GB}
                onChange={(v) => onChange({ ...hardware, diskBytes: Math.floor(v) })}
                units={byteUnits}
            />
        </div>
    );
}