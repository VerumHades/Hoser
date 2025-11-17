// src/components/Navbar.tsx
import React, { useState } from "react";
import { NavLink, useNavigate } from "react-router-dom";
import { Menu, X } from "lucide-react"; // Lucide icons
import Logo from "./components/prefabs/Logo";
import { UserSessionDisplay, useUserSession } from "./components/UserSession";
import { AnimatePresence, motion } from "framer-motion";

interface NavigationLinkProps {
    to: string,
    onClick?: () => void,
    className?: string,
    children?: React.ReactNode
}
function NavigationLink({ children, to, onClick, className }: NavigationLinkProps) {
    return <NavLink
        to={to}
        end
        className={className}
        onClick={onClick}
    >
        {children}
    </NavLink>
}

interface DevLinkProps {
    onClick?: () => void,
    className?: string,
    children?: React.ReactNode
}

function DeveloperLink({ className, children, onClick }: DevLinkProps) {
    const navigate = useNavigate()
    const { user } = useUserSession()
    const go = () => {
        onClick?.()
        if (user?.isDeveloper) {
            navigate("/account", { state: { dashpage: "Listings" } })
            return
        }

        navigate("/developer/join")
    }
    return <div onClick={go} className={className}>{children}</div>
}

interface NavbarProps {
    children?: React.ReactNode,
    className?: string
}

function Navbar({ children, className }: NavbarProps) {
    const [open, setIsOpen] = useState(false);

    const close = () => {
        setIsOpen(false)
    }

    const links = <>
        <NavigationLink to="/explore/public" onClick={close} className="navigation-link">
            Public Servers
        </NavigationLink>
        <NavigationLink to="/explore/listings" onClick={close} className="navigation-link">
            Listings
        </NavigationLink>
        <DeveloperLink onClick={close} className="navigation-link">Develop</DeveloperLink>
    </>

    return (
        <div className={"md:relative fixed inset-0 top-0 flex flex-col pointer-events-none z-50"}>
            <nav className="bg-white dark:bg-gray-900 shadow-md w-full pointer-events-auto">
                <div className="flex justify-between items-center h-16 px-5 relative">
                    {/* Logo + links */}
                    <div className="flex flex-row">
                        <Logo></Logo>
                        <div className="hidden sm:flex space-x-2">
                            {links}
                        </div>
                    </div>

                    {/* Account button desktop */}
                    <div className="hidden sm:flex"><UserSessionDisplay></UserSessionDisplay></div>

                    {/* Mobile menu toggle */}
                    <div className="sm:hidden flex items-center">
                        <button
                            onClick={() => setIsOpen(!open)}
                            className="p-2 focus:outline-none rounded hover:bg-gray-100 dark:hover:bg-gray-800 transition text-slate-800 dark:text-gray-200"
                        >
                            {open ? <X size={24} /> : <Menu size={24} />}
                        </button>
                    </div>
                </div>
            </nav>

            {/* Mobile dropdown */}
            <AnimatePresence>
                {open && <motion.nav
                    initial={{ y: -10, opacity: 0 }}
                    animate={{ y: 0, opacity: 1 }}
                    exit={{ y: -10, opacity: 0 }}
                    transition={{ type: "tween", duration: 0.3 }}
                    className="flex flex-1 flex-col justify-between bg-white dark:bg-gray-900 pointer-events-auto md:hidden"
                >
                    <div className="flex flex-col px-2 py-2 pointer-events-auto">{links}</div>
                    <div className="flex-1">
                        {children}
                    </div>
                    <UserSessionDisplay></UserSessionDisplay>
                </motion.nav>}
            </AnimatePresence> 
        </div>
    );
};

export default Navbar;
