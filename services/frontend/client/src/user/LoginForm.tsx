import { useState, type FormEvent } from "react";
import { toast, Toaster } from "react-hot-toast";
import { useNavigate } from "react-router-dom";

import backend_constants from "../backend/constants";
import { useUserSession } from "../components/restriction/UserSession";

export default function LoginForm() {
    const [username, setUsername] = useState<string>("");
    const [password, setPassword] = useState<string>("");

    const userSession = useUserSession();
    const navigate = useNavigate();

    /**
     * Handles submission of login credentials and session initialization.
     */
    const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
        event.preventDefault();

        try {
            const response = await fetch(`${backend_constants.address}/login`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    username,
                    password,
                }),
                credentials: "include",
            });

            if (!response.ok) {
                const errorMessage = await response.text();
                toast.error(errorMessage);
                return;
            }

            userSession.refresh();
            navigate("/");
        } catch (error) {
            toast.error("Network error");
            console.error(error);
        }
    };

    return (
        <div className="min-h-screen w-full flex items-center justify-center bg-slate-50 dark:bg-slate-950 px-6">
            <div className="w-full max-w-5xl grid grid-cols-1 md:grid-cols-2 gap-12 items-center">
                <section className="flex flex-col gap-6">
                    <h1 className="text-4xl md:text-5xl font-bold tracking-tight text-slate-900 dark:text-slate-100">
                        Welcome back.
                    </h1>

                    <p className="max-w-md text-lg text-slate-600 dark:text-slate-400">
                        Sign in to deploy, manage, and ship real infrastructure.
                    </p>

                    <img
                        src="/undraw_authentication_1evl.svg"
                        alt="Authentication illustration"
                        className="w-full max-w-md hidden md:block opacity-90 dark:opacity-80"
                    />
                </section>

                <section className="w-full">
                    <div className="w-full max-w-md mx-auto p-8 rounded-2xl bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 shadow-sm">
                        <h2 className="text-2xl font-semibold text-slate-900 dark:text-slate-100 mb-6">
                            Login
                        </h2>

                        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
                            <input
                                type="text"
                                value={username}
                                onChange={(event) => setUsername(event.target.value)}
                                placeholder="Username"
                                required
                                className="w-full px-4 py-3 rounded-lg bg-slate-50 dark:bg-slate-800 border border-slate-300 dark:border-slate-700 text-slate-900 dark:text-slate-100 placeholder:text-slate-400 focus:outline-none focus:ring-2 focus:ring-slate-900 dark:focus:ring-slate-100"
                            />

                            <input
                                type="password"
                                value={password}
                                onChange={(event) => setPassword(event.target.value)}
                                placeholder="Password"
                                required
                                className="w-full px-4 py-3 rounded-lg bg-slate-50 dark:bg-slate-800 border border-slate-300 dark:border-slate-700 text-slate-900 dark:text-slate-100 placeholder:text-slate-400 focus:outline-none focus:ring-2 focus:ring-slate-900 dark:focus:ring-slate-100"
                            />

                            <button
                                type="submit"
                                className="mt-2 w-full px-6 py-3 rounded-lg bg-slate-900 text-white dark:bg-slate-100 dark:text-slate-900 font-medium hover:opacity-90 transition"
                            >
                                Sign in
                            </button>
                        </form>
                    </div>
                </section>

                <Toaster position="top-right" />
            </div>
        </div>
    );
}
