import { useEffect, useState, useCallback } from "react";
import debounce from "lodash.debounce";
import LoadingIcon from "../prefabs/LoadingIcon";
import { Search } from "lucide-react";
import type { HasClassname } from "../common";

interface QueryProps<T> extends HasClassname {
    children?: React.ReactNode;
    bodyBuilder: (item?: T) => React.ReactElement;
    queryBuilder: (query: string) => Record<string, string>;
    endpoint: string;
    debounceDelay?: number;
}

const DEFAULT_DEBOUNCE_DELAY = 300;

export default function Query<T>({
    bodyBuilder,
    endpoint,
    queryBuilder,
    debounceDelay,
    children,
    className
}: QueryProps<T>) {
    const [query, setQuery] = useState("");
    const [results, setResults] = useState<T | undefined>(undefined);
    const [loading, setLoading] = useState(false);

    const fetchResultsDebounced = useCallback(
        debounce(async (queryValue: string, controller: AbortController) => {
            setLoading(true);

            try {
                const params = new URLSearchParams(queryBuilder(queryValue));
                const response = await fetch(`${endpoint}?${params}`, {
                    signal: controller.signal,
                    credentials: "include"
                });

                if (!response.ok) {
                    throw new Error("Failed to fetch results");
                }

                const data: T = await response.json();
                setResults(data);
            } catch (error) {
                if (error instanceof DOMException && error.name === "AbortError") {
                    return;
                }
                console.error("Query fetch error:", error);
            } finally {
                setLoading(false);
            }
        }, debounceDelay ?? DEFAULT_DEBOUNCE_DELAY),
        [endpoint, queryBuilder, debounceDelay]
    );

    useEffect(() => {
        const controller = new AbortController();
        fetchResultsDebounced(query, controller);

        return () => {
            controller.abort();
            fetchResultsDebounced.cancel();
        };
    }, [query, fetchResultsDebounced]);

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
                    placeholder="Search setups, stacks, infrastructure..."
                    value={query}
                    onChange={(event) => setQuery(event.target.value)}
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

                {children}
            </div>

            <div className="relative min-h-[3rem]">
                {loading ? (
                    <div className="absolute inset-0 flex items-center justify-center">
                        <LoadingIcon />
                    </div>
                ) : (
                    bodyBuilder(results)
                )}
            </div>
        </div>
    );
}
