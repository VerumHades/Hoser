import { Cpu, MemoryStick, HardDrive } from "lucide-react";
import HardwareSlider from "../HardwareSlider";
import type { HardwareUpdate } from "../../backend";

interface HardwareSettingsProps {
    hardware: HardwareUpdate;
    onChange: (update: HardwareUpdate) => void;
}

export default function HardwareSettings({
    hardware,
    onChange,
}: HardwareSettingsProps) {
    return (
        <div className="p-4">
            <HardwareSlider
                label="CPU Cores"
                icon={<Cpu className="w-5 h-5" />}
                min={1}
                max={64}
                value={hardware.cpu ?? 1}
                onChange={(v) => onChange({ ...hardware, cpu: v })}
                unit="cores"
            />
            <HardwareSlider
                label="RAM"
                icon={<MemoryStick className="w-5 h-5" />}
                min={1}
                max={1024}
                value={hardware.ram ?? 8}
                onChange={(v) => onChange({ ...hardware, ram: v })}
                unit="GB"
            />
            <HardwareSlider
                label="Disk"
                icon={<HardDrive className="w-5 h-5" />}
                min={128}
                max={8192}
                value={hardware.disk ?? 256}
                onChange={(v) => onChange({ ...hardware, disk: v })}
                unit="GB"
            />
        </div>
    );
};