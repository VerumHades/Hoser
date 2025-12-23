import React, { useState } from "react";
import Prompt from "./Prompt";

interface SelectBoxOption {
    label: string;
    description?: string;
}

interface SelectBoxProps {
    options: Record<string, SelectBoxOption>;
    defaultValue?: string;
    onSelected?: (value: string, oldValue: string) => void;
    onEmpty?: () => React.ReactNode;
}

export default function SelectBox({
    options,
    defaultValue,
    onSelected = () => {},
    onEmpty,
}: SelectBoxProps) {
    if (!options || Object.keys(options).length === 0) {
        return onEmpty?.();
    }

    const [selected, setSelected] = useState<string>(
        defaultValue || Object.keys(options)[0] || ""
    );
    const [open, setOpen] = useState<boolean>(false);

    const selectOption = (key: string) => {
        const previousValue = selected;
        setSelected(key);
        onSelected(key, previousValue);
    };

    const makeSelectHandler = (key: string) => () => {
        selectOption(key);
        setOpen(false);
    };

    return (
        <>
            <div
                className="border border-gray-300 dark:border-gray-600 rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:focus:ring-blue-400 dark:focus:border-blue-400 cursor-pointer"
                onClick={() => setOpen(true)}
            >
                {options[selected].label}
            </div>

            {open && (
                <Prompt
                    cancel={
                        <button
                            onClick={() => setOpen(false)}
                            className="bg-gray-200 hover:bg-gray-300 text-gray-900 font-semibold py-2 px-4 rounded dark:bg-gray-700 dark:text-gray-100 dark:hover:bg-gray-600"
                        >
                            Cancel
                        </button>
                    }
                >
                    <div className="flex flex-col">
                        {Object.keys(options).map((key) => {
                            const option = options[key];
                            const isSelected = key === selected;

                            const baseClasses =
                                "p-4 mb-2 rounded-lg border shadow-sm transition-colors cursor-pointer";
                            const selectedClasses =
                                "bg-blue-600 border-blue-600 text-white dark:bg-blue-500 dark:border-blue-500";
                            const unselectedClasses =
                                "bg-white border-gray-200 text-gray-900 hover:bg-gray-50 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-100 dark:hover:bg-gray-700";

                            return (
                                <div
                                    key={key}
                                    onClick={makeSelectHandler(key)}
                                    className={`${baseClasses} ${
                                        isSelected
                                            ? selectedClasses
                                            : unselectedClasses
                                    }`}
                                >
                                    <label className="block font-medium">
                                        {option.label}
                                    </label>
                                    {option.description && (
                                        <p className="mt-1 text-sm text-gray-700 dark:text-gray-300">
                                            {option.description}
                                        </p>
                                    )}
                                </div>
                            );
                        })}
                    </div>
                </Prompt>
            )}
        </>
    );
}
