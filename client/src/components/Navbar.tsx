// src/components/Navbar.tsx
import React from "react";
import { Menu, X } from "lucide-react"; // Lucide icons
import { motion } from "framer-motion";
import Logo from "./prefabs/Logo";
import { UserSessionDisplay } from "./UserSession";


interface NavbarProps {
    children?: React.ReactNode,
    className?: string,
    open: boolean,
    onToggleOpen: () => void
}

function Navbar({ children, className, onToggleOpen, open }: NavbarProps) {
    return (
        <div className={className}>
            <div className="flex justify-between items-center h-16 px-5 relative">
                {/* Logo + links */}
                <Logo></Logo>

                {/* Account button desktop */}
                <div className="hidden md:flex"><UserSessionDisplay></UserSessionDisplay></div>

                {/* Mobile menu toggle */}
                <div className="md:hidden flex items-center">
                    <button
                        onClick={onToggleOpen}
                        className="p-2 focus:outline-none rounded hover:bg-gray-100 dark:hover:bg-gray-800 transition text-slate-800 dark:text-gray-200"
                    >
                        {open ? <X size={24} /> : <Menu size={24} />}
                    </button>
                </div>
            </div>

            {open && <motion.nav
                initial={{ y: -10, opacity: 0 }}
                animate={{ y: 0, opacity: 1 }}
                exit={{ y: -10, opacity: 0 }}
                transition={{ type: "tween", duration: 0.3 }}
                className="fixed inset-0 top-16 z-50 flex flex-1 flex-col justify-between bg-white dark:bg-gray-900 pointer-events-auto md:hidden"
            >
                <div className="flex-1">
                    {children}
                </div>
                <UserSessionDisplay></UserSessionDisplay>
            </motion.nav>}
        </div>
    );
};

export default Navbar;
