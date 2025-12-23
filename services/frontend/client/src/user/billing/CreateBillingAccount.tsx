import { NavLink } from "react-router-dom";

export default function CreateBillingAccount() {
    return (
        <div className="flex justify-center items-center h-full">
            <NavLink
                to="/dashboard/billing"
                className="
                    inline-flex items-center justify-center
                    px-6 py-3
                    bg-blue-600 text-white font-semibold
                    rounded-lg
                    hover:bg-blue-700
                    dark:bg-blue-500 dark:text-white dark:hover:bg-blue-600
                    focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2
                    transition
                "
            >
                Create a billing account
            </NavLink>
        </div>
    );
}
