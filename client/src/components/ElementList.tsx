// src/components/DashboardApp.tsx
import { ChevronLeft } from "lucide-react";
import Search from "../components/Search";
import { AnimatePresence, motion } from "framer-motion";

interface ElementListProps<T> {
    element: T | undefined,
    children?: React.ReactNode,

    itemBuilder: (item: T) => React.ReactElement,
    itemViewBuilder: (item: T) => React.ReactElement,
    onResetElement: () => void,

    endpoint: string,
    queryBuilder: (query: string) => {},
}

export default function ElementList<T>({ queryBuilder, itemBuilder, itemViewBuilder, element, endpoint, onResetElement }: ElementListProps<T>) {

    const builder = (element: T) => {
        return itemBuilder(element)
    }

    return <AnimatePresence mode="wait">
        {element ? (
            <motion.div
                key="detail-view"
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -20 }}
                transition={{ duration: 0.3 }}
                className="space-y-4 flex flex-col w-full h-full text-gray-900 dark:text-gray-100"
            >
                <div
                    className="my-3 py-2 flex flex-row items-center hover:bg-slate-200 dark:hover:bg-gray-700 rounded-md transition-all"
                    onClick={() => onResetElement()}
                >
                    <ChevronLeft size={32} />
                    <label>Back</label>
                </div>
                <div className="flex-1">
                    {itemViewBuilder(element)}
                </div>
            </motion.div>
        ) : (
            <motion.div
                key="search-view"
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -20 }}
                transition={{ duration: 0.3 }}
                className="space-y-4 w-full h-full flex flex-col items-center"
            >
                <Search itemBuilder={builder} queryBuilder={queryBuilder} endpoint={endpoint} />
            </motion.div>
        )}
    </AnimatePresence>
}
