// src/components/AccountPage.tsx
import { Pencil, PenOff } from "lucide-react";
import { useState, type ChangeEvent } from "react";

interface EditableTextProps {
    text: string;
    onChange: (newText: string) => void;
}

export default function EditableText({ text, onChange }: EditableTextProps) {
    const [editing, setEditing] = useState<boolean>(false);

    const handleChange = (event: ChangeEvent<HTMLInputElement>) => {
        onChange(event.target.value);
    };

    const stopEditing = () => {
        setEditing(false);
    };

    const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
        if (e.key === "Enter" && editing) {
            stopEditing();
        }
    };

    return (
        <div className="flex flex-row justify-between items-baseline" onDoubleClick={() => setEditing(true)}>
            {editing ? (
                <>
                    <input
                        value={text}
                        autoFocus
                        onChange={handleChange}
                        onBlur={stopEditing}
                        onKeyDown={handleKeyDown}
                        className="border-b border-gray-300 dark:border-gray-600 focus:outline-none focus:ring-1 focus:ring-indigo-500 dark:bg-gray-800 dark:text-white px-1"
                    />
                    <PenOff className="cursor-pointer" onClick={stopEditing} />
                </>
            ) : (
                <>
                    <span>{text}</span>
                    <Pencil className="cursor-pointer" onClick={() => setEditing(true)} />
                </>
            )}
        </div>
    );
};
