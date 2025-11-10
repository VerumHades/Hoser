// src/components/RequireLogin.tsx
import React, { useEffect, useState } from "react";
import { NavLink, useNavigate } from "react-router-dom";
import { getUser, type User } from "../user";
import LoadingIcon from "./prefabs/LoadingIcon";

interface RequireLoginProps {
  children: React.ReactNode;
  condition?: (user: User) => boolean,
  no_access_element?: React.ReactNode;
  redirect?: string
}

const RequireLogin: React.FC<RequireLoginProps> = ({ children, condition, no_access_element, redirect }) => {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();
  
  useEffect(() => {
    async function fetchUser() {
      try {
        const u = await getUser();
        if((!u || (condition && !condition(u as User))) && redirect) navigate(redirect)

        setUser(u ?? null);
      } catch (err) {
        console.error("Error fetching user:", err);
        if(redirect) navigate(redirect)

        setUser(null);
      } finally {
        // Slightly faster load transition
        setTimeout(() => setLoading(false), 200);
      }
    }
    fetchUser();
  }, []);

  if (loading) {
    return (
      <div className="w-full h-full flex flex-col justify-center bg-slate-50">
        <LoadingIcon></LoadingIcon>
      </div>
    );
  }

  const default_required_login_component = <>
      <p className="text-gray-600 dark:text-gray-400 mb-6">
        This feature is only available to some users with an account.
      </p>
      <div className="flex flex-row">
        <NavLink
          to="/login"
          className="
  text-blue-600 dark:text-blue-400 hover:underline
  hover:text-blue-700 dark:hover:text-blue-300 font-medium transition-colors
  mr-5"
        >
          Login
        </NavLink>
        <label>Or</label>
        <NavLink
          to="/register"
          className="
  text-blue-600 dark:text-blue-400 hover:underline
  hover:text-blue-700 dark:hover:text-blue-300 font-medium transition-colors
  ml-5"
        >
          Register
        </NavLink>
      </div>
    </>

  if (!user || (condition && !condition(user))) {
    return (
      <div className="flex items-center justify-center w-full h-full bg-gray-50 dark:bg-gray-950 px-4">
        <div className="w-full text-center flex flex-col items-center">
          {
            no_access_element ? no_access_element : default_required_login_component
          }
        </div>
      </div>
    );
  }

  return <>{children}</>;
};

export default RequireLogin;
