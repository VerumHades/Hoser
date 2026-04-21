
import SelectBox from "../../components/input/SelectBox";

interface VisibilitySettingsProps {
    accessMode: number;
    onChange: (mode: number) => void;
}

export default function VisibilitySettings({ accessMode, onChange }: VisibilitySettingsProps) {
    return (
        <div className="mb-4">
            <label className="block font-semibold mb-1">Visibility</label>
            <SelectBox
                options={{}}
                defaultValue={accessMode.toString()}
                onSelected={(value) => onChange(Number(value))}
            />
        </div>
    );
}
