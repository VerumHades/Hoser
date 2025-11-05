
import type { ReactNode } from "react";

interface PromptProperties {
    children?: ReactNode,
    text?: ReactNode,
    cancelClassName?: string,
    submitClassName?: string,
    isSubmitLocked?: () => boolean,
    onCancel?: () => void,
    onSubmit?: () => void
}

export default function Prompt({ children, text, onCancel, onSubmit, cancelClassName, submitClassName, isSubmitLocked }: PromptProperties) {
    return (
        <div className="absolute inset-0 bg-black bg-opacity-30 flex flex-col justify-center items-center z-50">
            <div className="flex flex-col bg-white dark:bg-gray-800 rounded-lg shadow-lg max-w-sm w-full p-6 sm:max-w-md md:max-w-lg">
                <div className="flex-1">
                    {children}
                </div>
                <div className="flex flex-row justify-between">
                    <button onClick={onCancel} className={cancelClassName ? cancelClassName : "bg-gray-300 hover:bg-gray-400 text-gray-800 font-semibold py-2 px-4 rounded dark:bg-gray-700 dark:text-gray-200 dark:hover:bg-gray-600"}>Cancel</button>
                    <button onClick={onSubmit} disabled={isSubmitLocked?.()} className={submitClassName}>{text ? text : "Submit"}</button>
                </div>
            </div>
        </div>
    );
};
