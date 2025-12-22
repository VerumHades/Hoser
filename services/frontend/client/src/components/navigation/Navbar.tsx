import React from "react";
import { Menu, X } from "lucide-react";
import { motion, AnimatePresence } from "framer-motion";
import Logo from "../prefabs/Logo";
import { UserSessionDisplay } from "../restriction/UserSession";

interface NavbarProps {
    children?: React.ReactNode;
    className?: string;
    open: boolean;
    onToggleOpen: () => void;
}

function Navbar({ children, className, onToggleOpen, open }: NavbarProps) {
    return (
        <header
            className={`
                ${className ?? ""}
                sticky top-0 z-50
                bg-white dark:bg-gray-900
                border-b border-slate-200 dark:border-gray-800
            `}
        >
            <div className="max-w-full flex justify-between items-center h-16 px-5">
                <Logo />

                <div className="hidden md:flex">
                    <UserSessionDisplay />
                </div>

                <button
                    onClick={onToggleOpen}
                    className="
                        md:hidden
                        p-2
                        rounded-lg
                        text-slate-700 dark:text-gray-300
                        hover:bg-slate-100 dark:hover:bg-gray-800
                        transition
                    "
                >
                    {open ? <X size={22} /> : <Menu size={22} />}
                </button>
            </div>

            <AnimatePresence>
                {open && (
                    <motion.nav
                        initial={{ opacity: 0, y: -6 }}
                        animate={{ opacity: 1, y: 0 }}
                        exit={{ opacity: 0, y: -6 }}
                        transition={{ duration: 0.2, ease: "easeOut" }}
                        className="
                            md:hidden
                            border-t border-slate-200 dark:border-gray-800
                            bg-white dark:bg-gray-900
                        "
                    >
                        <div className="px-4 py-4 flex flex-col gap-4">
                            <div className="flex-1 overflow-y-auto">
                                {children}
                            </div>

                            <div className="pt-4 border-t border-slate-200 dark:border-gray-800">
                                <UserSessionDisplay />
                            </div>
                        </div>
                    </motion.nav>
                )}
            </AnimatePresence>
        </header>
    );
}

export default Navbar;
