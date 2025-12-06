import { useState } from "react";
import Prompt from "./Prompt";

interface SelectBoxOption {
    label: string,
    description?: string
}

interface SelectBoxProps {
    options: Record<string, SelectBoxOption>;
    defaultValue?: string;
    onSelected?: (value: string, oldValue: string) => void;
}

export default function SelectBox({
    options,
    defaultValue,
    onSelected = () => { },
}: SelectBoxProps) {
    const [selected, setSelected] = useState<string>(defaultValue || Object.keys(options)[0] || "");
    const [open, setOpen] = useState<boolean>(false)

    const select = (name: string) => {
        const oldValue = selected
        setSelected(name)
        onSelected(name, oldValue);
    }

    const makeSelectHandler = (key: string) => {
        return () => {
            select(key)
            setOpen(false)
        }
    }

    return (
        <>
            <div
                className="border border-gray-300 rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                onClick={() => setOpen(true)}
            >
                {options[selected].label}
            </div>
            {open ?
                <Prompt cancel={
                    <button
                        onClick={() => { setOpen(false); }}
                        className="bg-gray-300 hover:bg-gray-400 text-gray-800 font-semibold py-2 px-4 rounded
                        dark:bg-gray-700 dark:text-gray-200 dark:hover:bg-gray-600">
                        Cancel
                    </button>
                }>
                    <div className="flex flex-col">
                        {Object.keys(options).map((key) => {
                            const option = options[key];
                            return <div
                                key={key}
                                onClick={makeSelectHandler(key)}
                                className={`p-4 mb-2 rounded-lg border shadow-sm transition-colors cursor-pointer
                                    ${key == selected
                                        ? "bg-blue-600 border-blue-600 text-white dark:bg-blue-500 dark:border-blue-500"
                                        : "bg-white border-gray-200 text-gray-900 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-100 hover:bg-gray-50 dark:hover:bg-gray-700"
                                    }`}
                            >
                                <label
                                    className="block font-medium"
                                >
                                    {option.label}
                                </label>
                                <p className="mt-1 text-sm">
                                    {option.description}
                                </p>
                            </div>

                        })
                        }
                    </div>
                </Prompt> : <></>}
        </>
    );
}