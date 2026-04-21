import React, { useState, type ChangeEvent, type KeyboardEvent } from "react";
import { Pencil, PenOff } from "lucide-react";
import { colors, typography} from "../theme";

interface EditableTextProps {
    text: string;
    label?: string;
    onChange: (newText: string) => void;
    block?: boolean; // block-level full width
    className?: string; // extra custom styles
}

/**
 * EditableText is a reusable component for inline or block editing of text.
 * Fully theme-driven with support for light/dark mode, responsive width, and block/inline layout.
 */
export const EditableText: React.FC<EditableTextProps> = ({ text, onChange, label, block = false, className = "" }) => {
    const [editing, setEditing] = useState(false);

    const handleChange = (e: ChangeEvent<HTMLInputElement>) => onChange(e.target.value);
    const stopEditing = () => setEditing(false);
    const handleKeyDown = (e: KeyboardEvent<HTMLInputElement>) => { if (e.key === "Enter") stopEditing(); };

    const containerClasses = block
        ? `flex flex-col w-full gap-2 ${className}`
        : `flex flex-row items-center gap-2 ${className}`;

    const labelClasses = `${typography.label} ${block ? "mb-1" : "mr-2"} ${colors.text.light} dark:${colors.text.dark}`;
    const inputClasses = `${typography.input} border-b ${colors.border.light} dark:${colors.border.dark} focus:outline-none focus:ring-1 focus:ring-indigo-500 w-full`;
    const textClasses = `truncate ${colors.text.light} dark:${colors.text.dark} ${block ? "w-full" : ""}`;

    return (
        <div className={containerClasses}>
            {label && <label className={labelClasses}>{label}</label>}

            <div
                className={`flex ${block ? "flex-col w-full" : "flex-row items-center"} justify-between gap-2`}
                onDoubleClick={() => setEditing(true)}
            >
                {editing ? (
                    <div className="flex w-full items-center gap-2">
                        <input
                            value={text}
                            autoFocus
                            onChange={handleChange}
                            onBlur={stopEditing}
                            onKeyDown={handleKeyDown}
                            className={inputClasses}
                        />
                        <PenOff className="cursor-pointer" onClick={stopEditing} />
                    </div>
                ) : (
                    <div className="flex w-full items-center justify-between gap-2">
                        <span className={textClasses}>{text}</span>
                        <Pencil className="cursor-pointer" onClick={() => setEditing(true)} />
                    </div>
                )}
            </div>
        </div>
    );
};

export default EditableText;
