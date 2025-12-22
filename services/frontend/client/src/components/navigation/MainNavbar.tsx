import React, { createContext, useContext, useState } from "react";
import { NavLink } from "react-router-dom";
import { Menu, X } from "lucide-react";
import Logo from "../prefabs/Logo";
import { AnimatePresence, motion } from "framer-motion";
import { UserSessionDisplay } from "../restriction/UserSession";

interface NavigationLinkProps {
    to: string;
    onClick?: () => void;
    className?: string;
    children?: React.ReactNode;
}

function NavigationLink({
    children,
    to,
    onClick,
    className
}: NavigationLinkProps) {
    return (
        <NavLink
            to={to}
            end
            onClick={onClick}
            className={({ isActive }) =>
                `
                px-3 py-2 rounded-lg text-sm font-medium transition
                ${
                    isActive
                        ? "text-slate-900 dark:text-slate-100 bg-slate-100 dark:bg-slate-800"
                        : "text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-slate-200 hover:bg-slate-100 dark:hover:bg-slate-800"
                }
                ${className ?? ""}
                `
            }
        >
            {children}
        </NavLink>
    );
}

interface NavbarProps {
    children?: React.ReactNode;
    className?: string;
}

interface NavbarContextType {
    open: () => void;
    close: () => void;
}

const NavbarContext = createContext<NavbarContextType>({
    open: () => {},
    close: () => {}
});

export function useNavbar() {
    return useContext(NavbarContext);
}

function MainNavbar({ children }: NavbarProps) {
    const [isOpen, setIsOpen] = useState(false);

    const open = () => setIsOpen(true);
    const close = () => setIsOpen(false);

    const links = (
        <>
            <NavigationLink to="/search" onClick={close}>
                Explore
            </NavigationLink>
        </>
    );

    return (
        <NavbarContext.Provider value={{ open, close }}>
            <div className="fixed top-0 inset-x-0 z-50">
                <nav className="bg-white dark:bg-slate-900 border-b border-slate-200 dark:border-slate-800">
                    <div className="max-w-7xl mx-auto h-16 px-5 flex items-center justify-between">
                        <div className="flex items-center gap-6">
                            <Logo />
                            <div className="hidden sm:flex items-center gap-1">
                                {links}
                            </div>
                        </div>

                        <div className="hidden sm:flex items-center">
                            <UserSessionDisplay />
                        </div>

                        <button
                            onClick={() => setIsOpen(!isOpen)}
                            className="
                                sm:hidden
                                p-2
                                rounded-lg
                                text-slate-700
                                dark:text-slate-300
                                hover:bg-slate-100
                                dark:hover:bg-slate-800
                                transition
                            "
                        >
                            {isOpen ? <X size={22} /> : <Menu size={22} />}
                        </button>
                    </div>
                </nav>

                <AnimatePresence>
                    {isOpen && (
                        <motion.div
                            initial={{ opacity: 0, y: -8 }}
                            animate={{ opacity: 1, y: 0 }}
                            exit={{ opacity: 0, y: -8 }}
                            transition={{ duration: 0.2, ease: "easeOut" }}
                            className="
                                sm:hidden
                                bg-white
                                dark:bg-slate-900
                                border-b
                                border-slate-200
                                dark:border-slate-800
                            "
                        >
                            <div className="max-w-7xl mx-auto px-5 py-4 flex flex-col gap-4">
                                <div className="flex flex-col gap-1">
                                    {links}
                                </div>

                                <div className="pt-4 border-t border-slate-200 dark:border-slate-800">
                                    <UserSessionDisplay />
                                </div>

                                {children && (
                                    <div className="pt-2">
                                        {children}
                                    </div>
                                )}
                            </div>
                        </motion.div>
                    )}
                </AnimatePresence>
            </div>
        </NavbarContext.Provider>
    );
}

export default MainNavbar;
