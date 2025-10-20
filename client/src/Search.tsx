import { useEffect, useState, useCallback } from "react";
import backend_constants from "./backend_constants";
import debounce from "lodash.debounce"

interface SearchItem {
    ID: string | number;
    Title: string;
    Description: string;
    icon?: string;
    image?: string;
    author: string;
}

const DEBOUNCE_DELAY = 300;

export default function SearchWithFilters() {
    const [query, setQuery] = useState("");
    const [results, setResults] = useState<SearchItem[]>([]);
    const [loading, setLoading] = useState(false);

    const fetchResultsDebounced = useCallback(
        debounce(async (q: string, controller: AbortController) => {
            setLoading(true);
            try {
                const params = new URLSearchParams({ q });
                const res = await fetch(`${backend_constants.address}/rentals/public?${params}`, {
                    signal: controller.signal,
                });
                if (!res.ok) throw new Error("Failed to fetch results");

                const data: SearchItem[] = await res.json();
                setResults(data || []);
            } catch (err) {
                if (err instanceof DOMException && err.name === "AbortError") return;
                console.error("Fetch error:", err);
            } finally {
                setLoading(false);
            }
        }, DEBOUNCE_DELAY),
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
        <div className="max-w-7xl mx-auto p-6 space-y-2 w-screen h-screen">
            <div className="flex flex-col sm:flex-row gap-3 items-center w-full">
                <input
                    type="text"
                    placeholder="Search..."
                    value={query}
                    onChange={(e) => setQuery(e.target.value)}
                    className="flex-1 px-6 py-3 focus:outline-none bg-gray-300 dark:bg-gray-800 sm:w-auto sm:flex-1 w-full"
                />
            </div>

            {loading ? (
                <p className="text-center text-gray-500">Loading...</p>
            ) : results.length === 0 ? (
                <p className="text-center text-gray-400">No results found.</p>
            ) : (
                <div className="flex flex-col overflow-y-scroll">
                    {results.map((item) => (
                        <div
                            key={item.ID}
                            className="flex items-start gap-4 p-4 rounded-md shadow-lg bg-gray-900 mt-3"
                        >
                            <img
                                src={item.icon || item.image || "/placeholder.png"}
                                alt=""
                                className="w-12 h-12 object-cover rounded-xl"
                            />
                            <div>
                                <h3 className="font-semibold text-lg">{item.Title}</h3>
                                <p className="text-sm text-gray-600">{item.Description}</p>
                                <p className="text-xs text-gray-400 mt-1">By {item.author}</p>
                            </div>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
}
