import { useState, useRef, useEffect } from "react";
import { ChevronDown } from "lucide-react";
import { colors, borders, typography } from "../theme";

interface SelectOption<T = string> {
    label: string;
    value: T;
}

interface CustomSelectProps<T = string> {
    options: SelectOption<T>[];
    value: T | null;
    onChange: (value: T) => void;
    placeholder?: string;
    block?: boolean;
    className?: string;
}

/**
 * Custom dropdown component fully theme-driven, responsive, block-capable, dark/light ready.
 */
export function CustomSelect<T extends string | number>({
    options,
    value,
    onChange,
    placeholder = "Select...",
    block = false,
    className = "",
}: CustomSelectProps<T>) {
    const [open, setOpen] = useState(false);
    const containerRef = useRef<HTMLDivElement>(null);

    // Close dropdown on outside click
    useEffect(() => {
        const handleClickOutside = (e: MouseEvent) => {
            if (containerRef.current && !containerRef.current.contains(e.target as Node)) {
                setOpen(false);
            }
        };
        document.addEventListener("mousedown", handleClickOutside);
        return () => document.removeEventListener("mousedown", handleClickOutside);
    }, []);

    const handleSelect = (val: T) => {
        onChange(val);
        setOpen(false);
    };

    const containerClasses = block ? `w-full relative ${className}` : `relative inline-block ${className}`;
    const buttonClasses = `
        flex justify-between items-center cursor-pointer
        ${typography.input} ${borders.rounded} border ${colors.border.light} dark:${colors.border.dark}
        ${colors.text.light} dark:${colors.text.dark} px-3 py-1 w-full
    `;
    const listClasses = `
        absolute mt-1 w-full bg-white dark:bg-gray-900 border ${colors.border.light} dark:${colors.border.dark} 
        rounded shadow-lg z-50 max-h-60 overflow-auto
    `;
    const optionClasses = `
        px-3 py-2 cursor-pointer hover:bg-indigo-100 dark:hover:bg-indigo-800
        ${colors.text.light} dark:${colors.text.dark} 
    `;

    const selectedLabel = options.find((o) => o.value === value)?.label ?? placeholder;

    return (
        <div ref={containerRef} className={containerClasses}>
            <button
                type="button"
                className={buttonClasses}
                onClick={() => setOpen((prev) => !prev)}
            >
                <span>{selectedLabel}</span>
                <ChevronDown className="ml-2 w-4 h-4" />
            </button>

            {open && (
                <ul className={listClasses}>
                    {options.map((opt) => (
                        <li
                            key={String(opt.value)}
                            className={optionClasses}
                            onClick={() => handleSelect(opt.value)}
                        >
                            {opt.label}
                        </li>
                    ))}
                </ul>
            )}
        </div>
    );
}

export default CustomSelect;
