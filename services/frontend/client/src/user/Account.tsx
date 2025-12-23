import React, { useState } from "react";
import { Routes, Route, useNavigate } from "react-router-dom";
import { BarChart3, Images, Library, User as UserIcon, Wallet, Menu, AppWindow } from "lucide-react";
import backend_constants from "../backend_constants";
import RequireLogin from "../components/restriction/RequireLogin";
import { useUserSession } from "../components/restriction/UserSession";

import UserLibrary from "./library/Library";
import DeveloperListings from "./developer/DeveloperListings";
import UserBilling from "./billing/UserAccounts";
import Logo from "../components/prefabs/Logo";
import CreateInstance from "./library/CreateInstance";
import UserInstances from "./instances/UserInstances";
import DeveloperListingLoader from "./developer/DeveloperListingPage";

const AccountInfo: React.FC = () => {
    const session = useUserSession();
    const navigate = useNavigate();

    const handleLogout = async () => {
        try {
            await fetch(`${backend_constants.address}/logout`, {
                method: "POST",
                credentials: "include",
            });
            session.refresh();
            navigate("/");
        } catch (err) {
            console.error("Logout failed", err);
        }
    };

    return (
        <div className="w-full mx-auto p-6">
            <p className="mb-2 text-gray-700 dark:text-gray-300">
                <strong>Username:</strong> {session.user?.username}
            </p>
            <button
                onClick={handleLogout}
                className="px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700 transition"
            >
                Logout
            </button>
        </div>
    );
};

const AccountPage: React.FC = () => {
    const session = useUserSession();
    const [isSidebarOpen, setIsSidebarOpen] = useState(false);

    return (
        <RequireLogin>
            <div className="flex h-screen bg-slate-50 dark:bg-gray-950">
                {/* Sidebar */}
                <aside
                    className={`fixed md:static z-20 top-0 left-0 bottom-0 w-64 border-r border-slate-200 dark:border-gray-800 p-4 bg-white dark:bg-gray-900 transform ${
                        isSidebarOpen ? "translate-x-0" : "-translate-x-full"
                    } transition-transform duration-300 md:translate-x-0`}
                >
                    <div className="mb-6 flex justify-between md:block">
                        <Logo />
                        <button
                            className="md:hidden p-2"
                            onClick={() => setIsSidebarOpen(false)}
                        >
                            ✕
                        </button>
                    </div>
                    <nav className="flex flex-col gap-2">
                        <NavLink to="/dashboard/info" icon={<UserIcon />} label="Account Info" onClick={() => setIsSidebarOpen(false)} />
                        <NavLink to="/dashboard/library" icon={<Library />} label="My Library" onClick={() => setIsSidebarOpen(false)} />
                        <NavLink to="/dashboard/instances" icon={<AppWindow />} label="My Instances" onClick={() => setIsSidebarOpen(false)} />
                        <NavLink to="/dashboard/billing" icon={<Wallet />} label="Billing" onClick={() => setIsSidebarOpen(false)} />
                        {session.user?.isDeveloper && (
                            <>
                                <NavLink to="/dashboard/listings" icon={<Images />} label="Listings" onClick={() => setIsSidebarOpen(false)} />
                                <NavLink to="/dashboard/earnings" icon={<BarChart3 />} label="Earnings" onClick={() => setIsSidebarOpen(false)} />
                            </>
                        )}
                    </nav>
                </aside>

                {/* Mobile menu toggle */}
                {
                    !isSidebarOpen && <button
                        className="fixed top-4 left-4 z-30 md:hidden p-2 rounded bg-white dark:bg-gray-900 shadow"
                        onClick={() => setIsSidebarOpen(true)}
                    >
                        <Menu className="w-5 h-5" />
                    </button>
                }


                {/* Main content */}
                <main className="flex flex-1 overflow-auto min-h-0">
                    <Routes>
                        <Route path="info" element={<AccountInfo />} />
                        <Route path="library" element={<UserLibrary />} />
                        <Route path="create-instance/:listingId" element={<CreateInstance />} />
                        <Route path="instances" element={<UserInstances />} />
                        <Route path="developer/listing/:id" element={<DeveloperListingLoader />} />
                        <Route path="billing/*" element={<UserBilling />} />
                        {session.user?.isDeveloper && (
                            <>
                                <Route path="listings" element={<DeveloperListings />} />
                                <Route path="earnings" element={<div>Earnings content goes here</div>} />
                            </>
                        )}
                        <Route path="*" element={<AccountInfo />} />
                    </Routes>
                </main>
            </div>
        </RequireLogin>
    );
};

interface NavLinkProps {
    to: string;
    icon?: React.ReactNode;
    label: string;
    onClick?: () => void;
}

const NavLink: React.FC<NavLinkProps> = ({ to, icon, label, onClick }) => {
    const navigate = useNavigate();
    const handleClick = () => {
        navigate(to);
        onClick?.();
    };

    return (
        <button
            onClick={handleClick}
            className="w-full flex items-center gap-3 px-4 py-2 text-sm rounded-lg text-slate-600 dark:text-gray-400 hover:bg-slate-100 dark:hover:bg-gray-800"
        >
            {icon && <span className="w-4 h-4">{icon}</span>}
            <span>{label}</span>
        </button>
    );
};

export default AccountPage;
