import { Cpu, MemoryStick, HardDrive } from "lucide-react";
import HardwareSlider from "../../../components/input/HardwareSlider";
import type { HardwareSpecification } from "../../../backend";
import { useState } from "react";

interface HardwareSettingsProps {
    /** Optional initial hardware specification */
    initialSpec?: HardwareSpecification;
    /** Fired whenever the hardware spec changes */
    onChange: (spec: HardwareSpecification) => void;
}

export default function HardwareSettings({
    initialSpec,
    onChange,
}: HardwareSettingsProps) {
    const [spec, setSpec] = useState<HardwareSpecification>({
        cpu: initialSpec?.cpu ?? 4,
        ramBytes: initialSpec?.ramBytes ?? 8192,
        diskBytes: initialSpec?.diskBytes ?? 100,
    });

    const byteUnits = [
        { label: "B", multiplier: 1 },
        { label: "KB", multiplier: 1024 },
        { label: "MB", multiplier: 1024 ** 2 },
        { label: "GB", multiplier: 1024 ** 3 },
        { label: "TB", multiplier: 1024 ** 4 },
    ];

    const handleChange = (updatedSpec: Partial<HardwareSpecification>) => {
        const newSpec = { ...spec, ...updatedSpec };
        setSpec(newSpec);
        onChange(newSpec);
    };

    return (
        <div className="p-4 space-y-4">
            {/* CPU cores */}
            <HardwareSlider
                label="CPU Cores"
                icon={<Cpu className="w-5 h-5" />}
                min={1}
                max={64}
                value={spec.cpu || 4}
                onChange={(v) => handleChange({ cpu: v })}
            />

            {/* RAM */}
            <HardwareSlider
                label="RAM"
                icon={<MemoryStick className="w-5 h-5" />}
                min={1024 ** 2}          // 1 MB
                max={16 * 1024 ** 3}   // 1 TB (for demonstration)
                value={spec.ramBytes}
                onChange={(v) => handleChange({ ramBytes: Math.floor(v) })}
                units={byteUnits}
            />

            {/* Disk */}
            <HardwareSlider
                label="Disk"
                icon={<HardDrive className="w-5 h-5" />}
                min={1024 ** 2}            // 10 GB
                max={16 * 1024 ** 3}     // 10 TB
                value={spec.diskBytes}
                onChange={(v) => handleChange({ diskBytes: Math.floor(v) })}
                units={byteUnits}
            />
        </div>
    );
}
