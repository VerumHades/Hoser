import { useEffect, useState, useCallback } from "react";
import debounce from "lodash.debounce"
import LoadingIcon from "../prefabs/LoadingIcon";
import {Search} from "lucide-react"
import type { HasClassname } from "../common";

interface QueryProps<T> extends HasClassname{
    children?: React.ReactNode,
    bodyBuilder: (item?: T) => React.ReactElement,
    queryBuilder: (query: string) => Record<string,string>,
    endpoint: string
    debounceDelay?: number
}

const DEBOUNCE_DELAY = 300;
export default function Query<T>({ bodyBuilder, endpoint, queryBuilder, debounceDelay, children, className}: QueryProps<T>) {
    const [query, setQuery] = useState("");
    const [results, setResults] = useState<T | undefined>(undefined);
    const [loading, setLoading] = useState(false);

    const fetchResultsDebounced = useCallback(
        debounce(async (query: string, controller: AbortController) => {
            setLoading(true);
            try {
                const params = new URLSearchParams(queryBuilder(query));
                const res = await fetch(`${endpoint}?${params}`, {
                    signal: controller.signal,
                    credentials: "include"
                });
                if (!res.ok) throw new Error("Failed to fetch results");

                const data: T = await res.json();
                setResults(data);
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
        <div className={"flex flex-col " + className}>
            <div className="flex flex-col sm:flex-row gap-3 items-center w-full mb-3 
                dark:bg-slate-800 px-5 shadow-md">
                <input
                    type="text"
                    placeholder="Query..."
                    value={query}
                    onChange={(e) => setQuery(e.target.value)}
                    className="flex-1 px-6 py-3 focus:outline-none sm:w-auto sm:flex-1 w-full 
                                text-gray-900 dark:text-gray-100 transition-all dark:focus:bg-slate-600 focus:bg-slate-200"
                />
                <Search className="text-gray-900 dark:text-gray-100" />
                {children}
            </div>

            {loading ? (
                <LoadingIcon></LoadingIcon>
            ) : bodyBuilder(results)}
        </div>
    );
}
