// src/components/AccountPage.tsx
import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { getUser, type User } from "../user";
import backend_constants from "../backend_constants";
import UserRentalList from "./Rentals";
import RequireLogin from "../components/RequireLogin";
import Dashboard, { Page } from "../components/Dashboard";
import { ShoppingCart, User as UserIcon} from "lucide-react";
import LoadingIcon from "../components/prefabs/LoadingIcon";
import Logo from "../components/prefabs/Logo";

const AccountPage: React.FC = () => {
    const [user, setUser] = useState<User | undefined>(undefined);
    const [loading, setLoading] = useState(true);
    const navigate = useNavigate();

    useEffect(() => {
        async function fetchUser() {
            const u = await getUser();
            setUser(u);
            setLoading(false);
        }

        fetchUser();
    }, []);

    const handleLogout = async () => {
        try {
            await fetch(`${backend_constants.address}/logout`, { method: "POST", credentials: "include" });
            setUser(undefined);
            navigate("/user_logged_out");
        } catch (err) {
            console.error("Logout failed", err);
        }
    };

    if (loading) return <LoadingIcon></LoadingIcon>

    if (!user)
        return (
            <div className="p-5">
                <p>You are not logged in.</p>
            </div>
        );

    return (
        <RequireLogin>
            <div className="w-full h-full">
                <Dashboard logo={<Logo></Logo>}>
                    <Page name="Account Information" icon={<UserIcon />}>
                        <div className="md:max-w-md w-full mx-auto p-6 bg-white dark:bg-gray-800">
                            <h1 className="text-2xl font-bold mb-4 text-gray-900 dark:text-gray-100">
                                Account Details
                            </h1>
                            <p className="mb-2 text-gray-700 dark:text-gray-300">
                                <strong>Username:</strong> {user.Username}
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
                    <Page name="Developer" icon={<ShoppingCart />}>
                        <UserRentalList></UserRentalList>
                    </Page>
                    <Page name="My Listings" subpage_of="Developer" icon={<ShoppingCart />}>
                        <UserRentalList></UserRentalList>
                    </Page>
                </Dashboard>
            </div>
        </RequireLogin>
    );
};

export default AccountPage;
