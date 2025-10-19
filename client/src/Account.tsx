// src/components/AccountPage.tsx
import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { getUser, type User } from "./user";
import backend_constants from "./backend_constants";

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

  if (loading) return <div className="p-5">Loading...</div>;

  if (!user)
    return (
      <div className="p-5">
        <p>You are not logged in.</p>
      </div>
    );

  return (
    <div className="max-w-md mx-auto mt-10 p-6 bg-white dark:bg-gray-800 shadow-md rounded-lg">
      <h1 className="text-2xl font-bold mb-4 text-gray-900 dark:text-gray-100">
        Account Details
      </h1>
      <p className="mb-2 text-gray-700 dark:text-gray-300">
        <strong>Username:</strong> {user.Username}
      </p>
      <p className="mb-4 text-gray-700 dark:text-gray-300">
        <strong>Email:</strong> {user.Email || "Not provided"}
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

export default AccountPage;
