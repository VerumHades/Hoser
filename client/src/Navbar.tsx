// src/components/Navbar.tsx
import React, { useEffect, useState } from "react";
import { NavLink, useLocation, useNavigate } from "react-router-dom";
import { getUser, type User } from "./user";
import { Menu, X, User as UserIcon } from "lucide-react"; // Lucide icons

const Navbar: React.FC = () => {
    const [isOpen, setIsOpen] = useState(false);
    const [user, setUser] = useState<User | undefined>(undefined);
    const navigate = useNavigate();
    const location = useLocation();

    const links = [
        { name: "Explore", to: "/explore" },
        { name: "Rent", to: "/rent" },
        { name: "Develop", to: "/developer/dashboard" },
    ];

    const refreshUser = () => {
        async function fetchUser() {
            const u = await getUser();
            setUser(u);
        }
        fetchUser();
    };

    useEffect(refreshUser, []);

    useEffect(() => {
        if (location.pathname === "/user_logged_out") {
            setUser(undefined);
            navigate("/");
        } else if (location.pathname === "/user_logged_in") {
            refreshUser();
            navigate("/");
        }
    }, [location, navigate]);

    const builtLinks = links.map((x, i) => (
        <NavLink
            key={i}
            to={x.to}
            end
            className="text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-800 hover:text-gray-900 dark:hover:text-white px-4 py-2 rounded transition-colors"
            onClick={() => setIsOpen(false)}
        >
            {x.name}
        </NavLink>
    ));

    const accountButton = (
        <NavLink
            to={user ? "/account" : "/login"}
            onClick={() => setIsOpen(false)}
            className="flex items-center gap-2 p-2 m-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors"
        >
            <UserIcon size={20} className="text-gray-700 dark:text-gray-300" />
            <span className="font-medium text-gray-900 dark:text-gray-100">
                {user ? user.Username : "Login"}
            </span>
        </NavLink>
    );

    return (
        <div className="fixed inset-x-0 top-0 flex flex-col pointer-events-none z-50">
            <nav className="bg-white dark:bg-gray-900 shadow-md w-full pointer-events-auto">
                <div className="flex justify-between items-center h-16 px-5 relative">
                    {/* Logo + links */}
                    <div className="flex items-center space-x-4">
                        <img src="/ghost_logo.svg" className="h-8 w-auto" alt="Logo" />
                        <span className="text-2xl font-bold text-gray-900 dark:text-white">
                            gHost
                        </span>
                        <div className="hidden sm:flex space-x-2">{builtLinks}</div>
                    </div>

                    {/* Account button desktop */}
                    <div className="hidden sm:flex">{accountButton}</div>

                    {/* Mobile menu toggle */}
                    <div className="sm:hidden flex items-center">
                        <button
                            onClick={() => setIsOpen(!isOpen)}
                            className="p-2 focus:outline-none rounded hover:bg-gray-100 dark:hover:bg-gray-800 transition"
                        >
                            {isOpen ? <X size={24} /> : <Menu size={24} />}
                        </button>
                    </div>
                </div>
            </nav>

            {/* Mobile dropdown */}
            <div
                className={`sm:hidden transition-all overflow-hidden ${isOpen ? "max-h-screen" : "max-h-0"
                    } flex flex-col bg-white dark:bg-gray-900 pointer-events-auto`}
            >
                <div className="flex flex-col px-2 py-2">{builtLinks}</div>
                {accountButton}
            </div>
        </div>
    );
};

export default Navbar;
