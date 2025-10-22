import { useEffect, useState, useCallback } from "react";
import debounce from "lodash.debounce"

interface SearchProps<T> {
    children?: React.ReactNode,
    itemBuilder: (item: T) => React.ReactElement,
    queryBuilder: (query: string) => {},
    endpoint: string
    debounceDelay?: number
}

const DEBOUNCE_DELAY = 300;

export default function Search<T>({ itemBuilder, endpoint, queryBuilder, debounceDelay, children }: SearchProps<T>) {
    const [query, setQuery] = useState("");
    const [results, setResults] = useState<T[]>([]);
    const [loading, setLoading] = useState(false);

    const fetchResultsDebounced = useCallback(
        debounce(async (query: string, controller: AbortController) => {
            setLoading(true);
            try {
                const params = new URLSearchParams(queryBuilder(query));
                const res = await fetch(`${endpoint}?${params}`, {
                    signal: controller.signal,
                });
                if (!res.ok) throw new Error("Failed to fetch results");

                const data: T[] = await res.json();
                setResults(data || []);
            } catch (err) {
                if (err instanceof DOMException && err.name === "AbortError") return;
                console.error("Fetch error:", err);
            } finally {
                setLoading(false);
            }
        }, debounceDelay ?? DEBOUNCE_DELAY),
        []
    );

    useEffect(() => {
        const controller = new AbortController();
        fetchResultsDebounced(query, controller);

        return () => {
            controller.abort();
            fetchResultsDebounced.cancel(); // cancel any pending debounce
        };
    }, [query, fetchResultsDebounced]);

    return (
        <div className="max-w-7xl flex flex-col mx-auto p-6 space-y-2 w-screen h-screen">
            <div className="flex flex-col sm:flex-row gap-3 items-center w-full">
                <input
                    type="text"
                    placeholder="Search..."
                    value={query}
                    onChange={(e) => setQuery(e.target.value)}
                    className="flex-1 px-6 py-3 focus:outline-none bg-gray-300 dark:bg-gray-800 sm:w-auto sm:flex-1 w-full"
                />
                {children}
            </div>

            {loading ? (
                <p className="text-center text-gray-500">Loading...</p>
            ) : results.length === 0 ? (
                <p className="text-center text-gray-400">No results found.</p>
            ) : (
                <div className="flex flex-col overflow-y-auto flex-1">
                    {results.map((item) => itemBuilder(item))}
                </div>
            )}
        </div>
    );
}
