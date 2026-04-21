import { useMemo } from "react";
import SimpleMDE from "react-simplemde-editor";
import "easymde/dist/easymde.min.css";

interface MarkdownEditorProps {
    /** The markdown string to edit */
    value: string;
    /** Callback fired when content changes */
    onChange: (value: string) => void;
}

/**
 * A professional Markdown editor.
 * Styling is handled via the .markdown-editor-wrapper to avoid EasyMDE token errors.
 */
export default function MarkdownEditor({ value, onChange }: MarkdownEditorProps) {
    
    const editorOptions = useMemo(() => ({
        spellChecker: false,
        status: false,
        placeholder: "Enter detailed documentation...",
        minHeight: "400px",
        maxHeight: "600px",
        sideBySideFullscreen: false, 
        forceSync: true,
        // Use a single token here to avoid the InvalidCharacterError
        previewClass: "custom-markdown-preview",
    }), []);

    return (
        <div className="markdown-editor-wrapper w-full border border-slate-200 rounded-xl overflow-hidden bg-white">
            <SimpleMDE 
                value={value} 
                onChange={onChange} 
                options={editorOptions}
            />
        </div>
    );
}