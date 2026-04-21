import { useState } from "react";
import PromptButton from "../../components/input/PromptButton";
import { Delete } from "lucide-react";

interface DeleteListingPromptProps {
    onDelete: () => void;
}

export default function DeleteListingPrompt({ onDelete }: DeleteListingPromptProps) {
    const [confirmationText, setConfirmationText] = useState("");
    const deleteKeyword = "Delete";
    const isLocked = confirmationText !== deleteKeyword;

    return (
        <PromptButton
            className="bg-red-500 hover:bg-red-600 text-white font-semibold py-2 px-4 rounded flex gap-2"
            submitClassName={`flex gap-2 font-semibold py-2 px-4 rounded transition-colors ${
                isLocked
                    ? "bg-red-300 text-gray-200 cursor-not-allowed"
                    : "bg-red-500 hover:bg-red-600 text-white"
            }`}
            promptContents={
                <>
                    <p>
                        Deleting this listing is irreversible. All data and users' access will be lost. 
                        Type "Delete" to confirm.
                    </p>
                    <input
                        onKeyUp={(e) => setConfirmationText((e.target as HTMLInputElement).value)}
                        type="text"
                        className="w-full p-2 border rounded"
                    />
                </>
            }
            isSubmitLocked={() => isLocked}
            onCancel={() => setConfirmationText("")}
            onSubmit={onDelete}
        >
            Delete <Delete />
        </PromptButton>
    );
}
