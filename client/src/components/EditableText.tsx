// src/components/AccountPage.tsx
import { Pencil, PenOff } from "lucide-react"
import { useState, type ChangeEvent } from "react";

interface EditableTextProps {
    onChange: (newText: string) => void,
    text: string
}

export default function EditableText({ onChange, text }: EditableTextProps) {
    const [editing, setEditing] = useState<boolean>(false);
    const [currentText, setText] = useState<string>(text);

    const handleChange = (event: ChangeEvent<HTMLInputElement>) => {
        setText(event.target.value);
    };

    const stopEditing = () => {
        setEditing(false)
        onChange(currentText)
    }

    const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement> ) => {
        if (e.key === "Enter" && editing) {
            stopEditing()
        }
    };

    return (
        <div className="flex flex-row justify-between items-baseline" onDoubleClick={() => setEditing(true)}>
            {
                editing ? (
                    <>
                    <input value={currentText} autoFocus onChange={handleChange} onBlur={stopEditing} onKeyDown={handleKeyDown}></input><PenOff onClick={() => setEditing(false)}/>
                    </>
                ) : (
                    <>
                    {currentText}<Pencil onClick={() => setEditing(true)} />
                    </>
                )
            }
        </div>
    );
};
