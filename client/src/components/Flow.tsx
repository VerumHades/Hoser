import { AnimatePresence, motion } from "framer-motion";
import { ChevronLeft } from "lucide-react";
import React from "react";
import { createContext, useContext, useState } from "react";

interface FlowContextType {
    page: number,
    next: () => void,
    prev: () => void,
    setPage: (index: number) => void
}

const FlowContext = createContext<FlowContextType>({ page: 0, next: () => { }, prev: () => { }, setPage: () => {} });

interface FlowProps {
    children?: React.ReactNode
}

export function Flow({ children }: FlowProps) {
    const total = React.Children.count(children);
    const childArray = React.Children.toArray(children);

    const [page, setFlow] = useState(0);

    const next = () => setFlow((p) => (p + 1) % total);
    const prev = () => setFlow((p) => (p - 1 + total) % total);
    const setPage: (index: number) => void = (index:number) => setFlow(Math.max(Math.min(index, total - 1),0)) 

    return (
        <FlowContext.Provider value={{ page, next, prev, setPage }}>
            <AnimatePresence mode="wait">
                {childArray.map((child, i) => {
                    if (i !== page) return null;

                    return (
                        <motion.div
                            key={i} // important for AnimatePresence to detect changes
                            initial={{ opacity: 0, x: 50 }}
                            animate={{ opacity: 1, x: 0 }}
                            exit={{ opacity: 0, x: -50 }}
                            transition={{ duration: 0.3 }}
                            className="w-full h-full"
                        >
                            {child}
                        </motion.div>
                    );
                })}
            </AnimatePresence>
        </FlowContext.Provider>
    );
}

interface FlowSwitchProps {
    children?: React.ReactNode
    direction: string
    className?: string
}

export function FlowSwitch({ direction, children, className }: FlowSwitchProps) {
    const { next, prev } = useContext(FlowContext);
    const action = direction === "next" ? next : prev;

    return (
        <div onClick={action} className={className}>
            {children}
        </div>
    );
}

interface FlowTopBackWrapperProps {
    children?: React.ReactNode
}

export function FlowTopBackWrapper({ children }: FlowTopBackWrapperProps) {
    return <div className="space-y-4 flex flex-col w-full h-full">
        <FlowSwitch
            className="my-3 py-2 flex text-slate-800 dark:text-slate-100 flex-row items-center hover:bg-slate-200 dark:hover:bg-gray-700 transition-all"
            direction="back"
        >
            <ChevronLeft size={32} />
            <label>Back</label>
        </FlowSwitch>
        <div className="flex-1">
            {children}
        </div>
    </div>
}

export function useFlow() {
    const ctx = useContext(FlowContext);
    if (!ctx) throw new Error("useFlow must be used within a Flow");
    return ctx;
}