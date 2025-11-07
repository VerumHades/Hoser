
import type { ReactNode } from "react";

interface PromptProperties {
    children?: ReactNode,
    submit?: ReactNode,
    cancel?: ReactNode
}

export default function Prompt({ children,  submit, cancel }: PromptProperties) {
    return (
        <div className="absolute inset-0 bg-black bg-opacity-30 flex flex-col justify-center items-center z-50">
            <div className="flex flex-col bg-white dark:bg-gray-800 rounded-lg shadow-lg max-w-sm w-full p-6 sm:max-w-md md:max-w-lg">
                <div className="flex-1">
                    {children}
                </div>
                <div className="flex flex-row justify-between">
                    {cancel}
                    {submit}
                </div>
            </div>
        </div>
    );
};
