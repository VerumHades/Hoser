import React from "react";
import SelectBox from "../../components/input/SelectBox";
import { ListingAccessModes } from "../../backend";

interface VisibilitySettingsProps {
    accessMode: number;
    onChange: (mode: number) => void;
}

export default function VisibilitySettings({ accessMode, onChange }: VisibilitySettingsProps) {
    const options: Record<string, { label: string; description: string }> = Object.fromEntries(
        Object.entries(ListingAccessModes).map(([key, value]) => [key, { label: value.label, description: value.description }])
    );

    return (
        <div className="mb-4">
            <label className="block font-semibold mb-1">Visibility</label>
            <SelectBox
                options={options}
                defaultValue={accessMode.toString()}
                onSelected={(value) => onChange(Number(value))}
            />
        </div>
    );
}
