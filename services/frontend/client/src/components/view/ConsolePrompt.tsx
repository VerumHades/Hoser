import { useState } from "react"

type ConsolePromptProps = {
    errorText: string
}

/**
 * Full-screen overlay that mimics a console window.
 * Supports basic ANSI colors and resets.
 */
export function ConsolePrompt({ errorText }: ConsolePromptProps) {
    const [isOpen, setIsOpen] = useState(false)

    const ansiColorMap: Record<string, string> = {
        "30": "text-black",
        "31": "text-red-500 dark:text-red-400",
        "32": "text-green-500 dark:text-green-400",
        "33": "text-yellow-500 dark:text-yellow-400",
        "34": "text-blue-500 dark:text-blue-400",
        "35": "text-purple-500 dark:text-purple-400",
        "36": "text-cyan-500 dark:text-cyan-400",
        "37": "text-white dark:text-slate-200",
        "90": "text-gray-400 dark:text-gray-500",
        "91": "text-red-400 dark:text-red-300",
        "92": "text-green-400 dark:text-green-300",
        "93": "text-yellow-400 dark:text-yellow-300",
        "94": "text-blue-400 dark:text-blue-300",
        "95": "text-purple-400 dark:text-purple-300",
        "96": "text-cyan-400 dark:text-cyan-300",
        "97": "text-white dark:text-white",
    }

    function parseAnsi(text: string) {
        const parts: React.ReactNode[] = []
        const ansiRegex = /\x1b\[(\d+(;\d+)*)m/g
        let lastIndex = 0
        let match: RegExpExecArray | null
        let currentClass = ""

        while ((match = ansiRegex.exec(text)) !== null) {
            if (match.index > lastIndex) {
                parts.push(
                    <span className={currentClass} key={lastIndex}>
                        {text.slice(lastIndex, match.index)}
                    </span>
                )
            }

            const codes = match[1].split(";")
            if (codes.includes("0")) currentClass = ""
            const colorCode = codes.find(code => ansiColorMap[code])
            if (colorCode) currentClass = ansiColorMap[colorCode]

            lastIndex = match.index + match[0].length
        }

        if (lastIndex < text.length) {
            parts.push(
                <span className={currentClass} key={lastIndex}>
                    {text.slice(lastIndex)}
                </span>
            )
        }

        return parts
    }

    return (
        <>
            <button
                type="button"
                onClick={() => setIsOpen(true)}
                className="text-xs font-medium text-red-700 dark:text-red-300 hover:underline"
            >
                Open Console
            </button>

            {isOpen && (
                <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 dark:bg-black/80 p-4 text-left">
                    <div className="relative w-full max-w-5xl max-h-[90vh] flex flex-col bg-gray-900 dark:bg-gray-800 rounded-lg shadow-lg border border-gray-700">
                        {/* Title bar */}
                        <div className="flex justify-between items-center bg-gray-800 dark:bg-gray-900 px-4 py-2 rounded-t-lg border-b border-gray-700">
                            <span className="text-sm text-gray-200 font-mono">Console Output</span>
                            <button
                                onClick={() => setIsOpen(false)}
                                className="text-sm font-bold text-red-500 hover:text-red-400"
                            >
                                ✕
                            </button>
                        </div>

                        {/* Console body */}
                        <pre className="flex-1 overflow-auto p-4 text-xs font-mono text-gray-200 dark:text-gray-100 bg-gray-900 dark:bg-gray-800 whitespace-pre-wrap break-words">
                            {parseAnsi(errorText)}
                        </pre>
                    </div>
                </div>
            )}
        </>
    )
}