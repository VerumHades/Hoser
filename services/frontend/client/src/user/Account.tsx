// src/components/AccountPage.tsx
import React from "react";
import backend_constants from "../backend_constants";

import { BarChart3, Images, ShoppingCart, User as UserIcon } from "lucide-react";
import Logo from "../components/prefabs/Logo";

import { useNavigate } from "react-router-dom";
import RequireLogin from "../components/restriction/RequireLogin";
import Dashboard, { Page } from "../components/navigation/Dashboard";
import UserLibrary from "./library/Library";
import DeveloperListings from "./developer/DeveloperListings";
import { useUserSession } from "../components/restriction/UserSession";

const AccountPage: React.FC = () => {
    const session = useUserSession()
    const navigate = useNavigate()

    const handleLogout = async () => {
        try {
            await fetch(`${backend_constants.address}/logout`, { method: "POST", credentials: "include" });
            session.refresh()
            navigate("/")
        } catch (err) {
            console.error("Logout failed", err);
        }
    };

    return (
        <RequireLogin>
            <div className="w-full h-full flex flex-col">
                <Dashboard logo={<Logo></Logo>}>
                    <Page name="Account Information" icon={<UserIcon />}>
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
                    </Page>
                    <Page name="My Rentals" icon={<ShoppingCart />}>
                        <UserLibrary></UserLibrary>
                    </Page>
                    {session.user?.isDeveloper && <>
                        <Page name="Listings" icon={<Images />}>
                            <DeveloperListings></DeveloperListings>
                        </Page>,
                        <Page name="Earnings" icon={<BarChart3 />}>
                            <div>Earnings content goes here</div>
                        </Page>
                    </>
                    }
                </Dashboard>
            </div>
        </RequireLogin>
    );
};

export default AccountPage;
