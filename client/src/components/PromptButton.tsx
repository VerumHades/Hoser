import { useState, type ReactNode } from "react";
import Prompt from "./Prompt";

interface PromptButtonProps {
    children?: ReactNode,
    contents: ReactNode,
    className?: string,
    isSubmitLocked?: () => boolean,
    submitClassName?: string,
    onCancel?: () => void,
    onSubmit?: () => void
}

export default function PromptButton({ children, contents, className, onSubmit, onCancel, isSubmitLocked, submitClassName }: PromptButtonProps) {
    const [open, setOpen] = useState<boolean>(false)

    const submit = () => {
        onSubmit?.()
        setOpen(false); 
    }
    const cancel = () => {
        onCancel?.()
        setOpen(false)
    }

    if (open) {
        return <Prompt text={contents} isSubmitLocked={isSubmitLocked} onSubmit={submit} onCancel={cancel} submitClassName={submitClassName}>{children}</Prompt>
    }

    return (
        <button className={className} onClick={() => setOpen(true)}>
            {contents}
        </button>
    );
};
