import { useState, type FormEvent } from "react";
import { toast, Toaster } from 'react-hot-toast';
import { useNavigate } from "react-router-dom";

import backend_constants from "../backend_constants";

export default function LoginForm() {
    const [username, setUsername] = useState<string>("");
    const [password, setPassword] = useState<string>("");
    
    const navigate = useNavigate(); 

    const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
        e.preventDefault();

        try {
            const response = await fetch(`${backend_constants.address}/login`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/x-www-form-urlencoded",
                },
                body: new URLSearchParams({ username, password }),
                credentials: "include" 
            });

            if (response.ok) {
                navigate("/user_logged_in");
            } else {
                const text = await response.text();
                toast.error(text)
            }
        } catch (err) {
            toast.error("Network error");
            console.error(err);
        }
    };

    return (
        <div className="flex items-center justify-center h-screen w-screen margin-0 bg-gray-100 dark:bg-gray-900">
            <div className="w-96 p-8 text-center max-w-11/12 text-gray-800 dark:text-gray-10">
                <h2 className="text-2xl font-bold mb-6">Login</h2>

                <form onSubmit={handleSubmit} className="space-y-4">
                    <input
                        type="text"
                        placeholder="Username"
                        value={username}
                        onChange={(e) => setUsername(e.target.value)}
                        required
                        className="w-full p-3 rounded border border-gray-300 dark:border-gray-600 focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700 dark:text-gray-100"
                    />

                    <input
                        type="password"
                        placeholder="Password"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        required
                        className="w-full p-3 rounded border border-gray-300 dark:border-gray-600 focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700 dark:text-gray-100"
                    />

                    <button
                        type="submit"
                        className="w-full p-3 bg-blue-500 text-white rounded hover:bg-blue-600 transition-colors"
                    >
                        Login
                    </button>
                </form>

                <Toaster position="top-right" />
            </div>
        </div>
    );
}
