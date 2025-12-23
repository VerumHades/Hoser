import { useEffect, useState } from "react"
import LoadingIcon from "../prefabs/LoadingIcon"
import { Search } from "lucide-react"
import type { HasClassname } from "../common"

interface QueryProps<ResultType> extends HasClassname {
    endpoint: string
    queryBuilder: (queryText: string) => Record<string, string>
    render: (results: ResultType | undefined) => React.ReactElement
    debounceDelayMilliseconds?: number
    toolbarRightSlot?: React.ReactNode
    searchPlaceholder?: string
}

const DEFAULT_DEBOUNCE_DELAY_MILLISECONDS = 300

/**
 * Query provides debounced, abort-safe search querying with a controlled render boundary.
 */
export default function Query<ResultType>({
    endpoint,
    queryBuilder,
    render,
    debounceDelayMilliseconds = DEFAULT_DEBOUNCE_DELAY_MILLISECONDS,
    toolbarRightSlot,
    searchPlaceholder = "Search...",
    className
}: QueryProps<ResultType>) {
    const [queryText, setQueryText] = useState("")
    const [results, setResults] = useState<ResultType | undefined>(undefined)
    const [isLoading, setIsLoading] = useState(false)

    useEffect(() => {
        const abortController = new AbortController()
        const timeoutHandle = window.setTimeout(async () => {
            setIsLoading(true)

            try {
                const searchParameters = new URLSearchParams(
                    queryBuilder(queryText)
                )

                const response = await fetch(
                    `${endpoint}?${searchParameters.toString()}`,
                    {
                        credentials: "include",
                        signal: abortController.signal
                    }
                )

                if (!response.ok) {
                    throw new Error("Query request failed")
                }

                const parsedResults: ResultType = await response.json()
                setResults(parsedResults)
            } catch (error) {
                if (
                    error instanceof DOMException &&
                    error.name === "AbortError"
                ) {
                    return
                }
                console.error("Query error:", error)
            } finally {
                setIsLoading(false)
            }
        }, debounceDelayMilliseconds)

        return () => {
            abortController.abort()
            window.clearTimeout(timeoutHandle)
        }
    }, [queryText])

    return (
        <div className={`flex flex-col gap-4 ${className ?? ""}`}>
            <div
                className="
                    flex items-center gap-3
                    w-full
                    px-4 py-3
                    rounded-xl
                    border
                    border-slate-200
                    dark:border-slate-700
                    bg-white
                    dark:bg-slate-900
                    shadow-sm
                    focus-within:ring-2
                    focus-within:ring-slate-400
                    dark:focus-within:ring-slate-600
                    transition
                "
            >
                <Search className="w-5 h-5 text-slate-500 dark:text-slate-400" />

                <input
                    type="text"
                    value={queryText}
                    onChange={(event) =>
                        setQueryText(event.target.value)
                    }
                    placeholder={searchPlaceholder}
                    className="
                        flex-1
                        bg-transparent
                        outline-none
                        text-sm
                        text-slate-900
                        dark:text-slate-100
                        placeholder:text-slate-400
                        dark:placeholder:text-slate-500
                    "
                />

                {toolbarRightSlot}
            </div>

            <div className="flex min-h-[3rem]">
                {isLoading ? (
                    <div className="absolute inset-0 flex items-center justify-center">
                        <LoadingIcon />
                    </div>
                ) : (
                    render(results)
                )}
            </div>
        </div>
    )
}
