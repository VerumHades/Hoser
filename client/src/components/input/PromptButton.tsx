import { useState, type ReactNode } from "react";
import Prompt from "./Prompt";

interface PromptButtonProps {
    children?: ReactNode,
    promptContents: ReactNode,
    className?: string,
    isSubmitLocked?: () => boolean,
    submitClassName?: string,
    onCancel?: () => void,
    onSubmit?: () => void
}

export default function PromptButton({ children, promptContents, className, onSubmit, onCancel, isSubmitLocked, submitClassName }: PromptButtonProps) {
    const [open, setOpen] = useState<boolean>(false)

    const submit = () => {
        onSubmit?.()
        setOpen(false);
    }
    const cancel = () => {
        onCancel?.()
        setOpen(false)
    }
    return (
        <>
            {open ?
                <Prompt
                    submit={
                        <button onClick={submit} disabled={isSubmitLocked?.()} className={submitClassName}>{children}</button>
                    }
                    cancel={
                        <button
                            onClick={cancel}
                            className="bg-gray-300 hover:bg-gray-400 text-gray-800 font-semibold py-2 px-4 rounded
             dark:bg-gray-700 dark:text-gray-200 dark:hover:bg-gray-600">
                            Cancel
                        </button>
                    }
                >
                    {promptContents}
                </Prompt> : <></>
            }
            <button className={className} onClick={() => setOpen(true)}>
                {children}
            </button>
        </>
    );
};
