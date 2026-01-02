import { UserIcon } from "lucide-react";
import React, { useEffect } from "react";
import { createContext, useContext, useState } from "react";

import { NavLink } from "react-router-dom";
import { UserAPI, type User } from "../../backend/repositories/user";

interface UserSessionType {
    user: User | undefined,
    refresh: () => void
}

const UserSessionContext = createContext<UserSessionType>({ user: undefined, refresh: () => {} });

interface UserSessionProps {
    children?: React.ReactNode
}

export function UserSession({ children }: UserSessionProps) {
    const [user, setUser] = useState<User | undefined>(undefined);

    const refreshUser = () => {
        async function fetchUser() {
            const u = await UserAPI.getData();
            setUser(u);
        }
        fetchUser();
    };
    useEffect(refreshUser, []);

    return (
        <UserSessionContext.Provider value={{ user, refresh: refreshUser }}>
            {children}
        </UserSessionContext.Provider>
    );
}

export function UserSessionDisplay() {
    let { user } = useContext(UserSessionContext);
    return <NavLink
        to={user ? "/dashboard" : "/login"}
        className="flex items-center gap-2 p-2 m-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors"
    >
        <UserIcon size={20} className="text-gray-700 dark:text-gray-300" />
        <span className="font-medium text-gray-900 dark:text-gray-100">
            {user ? user.username : "Login"}
        </span>
    </NavLink>
}

export function useUserSession() {
    const ctx = useContext(UserSessionContext);
    if (!ctx) throw new Error("useUserSession must be used within a UserSession");
    return ctx;
}