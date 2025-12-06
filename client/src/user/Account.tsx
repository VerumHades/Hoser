// src/components/AccountPage.tsx
import React from "react";
import backend_constants from "../backend_constants";

import { BarChart3, Images, ShoppingCart, User as UserIcon } from "lucide-react";
import Logo from "../components/prefabs/Logo";

import { useNavigate } from "react-router-dom";
import RequireLogin from "../components/restriction/RequireLogin";
import Dashboard, { Page } from "../components/navigation/Dashboard";
import UserRentalList from "./rentals/Rentals";
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
                        <div className="md:max-w-md w-full mx-auto p-6 bg-white dark:bg-gray-800">
                            <h1 className="text-2xl font-bold mb-4 text-gray-900 dark:text-gray-100">
                                Account Details
                            </h1>
                            <p className="mb-2 text-gray-700 dark:text-gray-300">
                                <strong>Username:</strong> {session.user?.Username}
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
                        <UserRentalList></UserRentalList>
                    </Page>
                    {session.user?.isDeveloper && <>
                        <Page name="Developer" icon={<ShoppingCart />} children={undefined}>

                        </Page>,
                        <Page name="Listings" subpage_of="Developer" icon={<Images />}>
                            <DeveloperListings></DeveloperListings>
                        </Page>,
                        <Page name="Earnings" subpage_of="Developer" icon={<BarChart3 />}>
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
