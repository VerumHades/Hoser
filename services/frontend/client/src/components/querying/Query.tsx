import { useEffect, useState, useCallback } from "react"
import LoadingIcon from "../prefabs/LoadingIcon"
import { Search } from "lucide-react"
import type { HasClassname } from "../common"

interface PaginatedResponse<ItemType> {
    items: ItemType[]
    cursor?: string
}

interface QueryProps<ItemType> extends HasClassname {
    endpoint: string
    queryBuilder: (queryText: string) => Record<string, string>
    render: (
        results: ItemType[],
        goToNext?: () => void,
        goToPrev?: () => void,
        currentPage?: number,
        totalPages?: number
    ) => React.ReactElement
    debounceDelayMilliseconds?: number
    toolbarRightSlot?: React.ReactNode
    searchPlaceholder?: string
}

const DEFAULT_DEBOUNCE_DELAY_MILLISECONDS = 300

export default function Query<ItemType>({
    endpoint,
    queryBuilder,
    render,
    debounceDelayMilliseconds = DEFAULT_DEBOUNCE_DELAY_MILLISECONDS,
    toolbarRightSlot,
    searchPlaceholder = "Search...",
    className
}: QueryProps<ItemType>) {
    const [queryText, setQueryText] = useState("")
    const [results, setResults] = useState<ItemType[]>([])
    
    const [cursorHistory, setCursorHistory] = useState<(string | undefined)[]>([undefined])
    const [currentCursorIndex, setCurrentCursorIndex] = useState(0)
    const [nextCursor, setNextCursor] = useState<string | undefined>(undefined)
    
    const [isLoading, setIsLoading] = useState(false)

    const fetchBatch = useCallback(async (cursor?: string) => {
        setIsLoading(true)
        try {
            const params = new URLSearchParams(queryBuilder(queryText))
            if (cursor) params.set("cursor", cursor)

            const response = await fetch(`${endpoint}?${params.toString()}`, {
                credentials: "include"
            })
            if (!response.ok) throw new Error("Query request failed")

            const data: PaginatedResponse<ItemType> = await response.json()
            setResults(data.items)
            setNextCursor(data.cursor)

            // Update cursor history if we moved forward
            setCursorHistory(prev => {
                const newHistory = prev.slice(0, currentCursorIndex + 1)
                if (data.cursor) newHistory.push(data.cursor)
                return newHistory
            })
        } catch (error) {
            console.error("Query error:", error)
        } finally {
            setIsLoading(false)
        }
    }, [endpoint, queryText, queryBuilder, currentCursorIndex])

    // Debounced search
    useEffect(() => {
        const timeout = window.setTimeout(() => fetchBatch(undefined), debounceDelayMilliseconds)
        return () => window.clearTimeout(timeout)
    }, [queryText, debounceDelayMilliseconds])

    const goToNext = useCallback(() => {
        if (nextCursor && currentCursorIndex < cursorHistory.length - 1) {
            setCurrentCursorIndex(prev => prev + 1)
            fetchBatch(cursorHistory[currentCursorIndex + 1])
        } else if (nextCursor) {
            setCurrentCursorIndex(prev => prev + 1)
            fetchBatch(nextCursor)
        }
    }, [nextCursor, currentCursorIndex, cursorHistory, fetchBatch])

    const goToPrev = useCallback(() => {
        if (currentCursorIndex > 0) {
            setCurrentCursorIndex(prev => prev - 1)
            fetchBatch(cursorHistory[currentCursorIndex - 1])
        }
    }, [currentCursorIndex, cursorHistory, fetchBatch])

    const currentPage = currentCursorIndex + 1
    const totalPages = nextCursor ? cursorHistory.length + 1 : cursorHistory.length

    return (
        <div className={`flex flex-col gap-4 ${className ?? ""}`}>
            <div className="flex items-center gap-3 w-full px-4 py-3 rounded-xl border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-900 shadow-sm focus-within:ring-2 focus-within:ring-slate-400 dark:focus-within:ring-slate-600 transition">
                <Search className="w-5 h-5 text-slate-500 dark:text-slate-400" />
                <input
                    type="text"
                    value={queryText}
                    onChange={e => {
                        setQueryText(e.target.value)
                        setCursorHistory([undefined])
                        setCurrentCursorIndex(0)
                    }}
                    placeholder={searchPlaceholder}
                    className="flex-1 bg-transparent outline-none text-sm text-slate-900 dark:text-slate-100 placeholder:text-slate-400 dark:placeholder:text-slate-500"
                />
                {toolbarRightSlot}
            </div>

            <div className="flex flex-col min-h-[3rem] relative">
                {isLoading && results.length === 0 ? (
                    <div className="absolute inset-0 flex items-center justify-center">
                        <LoadingIcon />
                    </div>
                ) : (
                    render(results, goToNext, goToPrev, currentPage, totalPages)
                )}
            </div>
        </div>
    )
}
