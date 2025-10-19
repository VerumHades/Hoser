// src/components/DashboardApp.tsx
import React, { useEffect, useState, type ReactNode } from "react";

interface PageProps {
    name: string;
    icon: ReactNode;
    children: ReactNode;
}

/**
 * <Page> is a wrapper component that only provides props to the parent.
 * It does not render anything itself.
 */
export const Page: React.FC<PageProps> = () => {
    // returns nothing
    return null;
};

type DashboardProps = {
    children: ReactNode;
}

export default function Dashboard({ children }: DashboardProps) {
    const [current_page, setPage] = useState<string>("listings");

    const pages = React.Children.toArray(children).map((child) => {
        console.log(child)
        if (!React.isValidElement<PageProps>(child)) {
            throw new Error("<Dashboard> children must be <Page> components");
        }

        return child.props;
    });
    useEffect(() => {
        setPage(pages[0].name)
    }, [])

    return (
        <div className="flex flex-row h-full bg-slate-50 dark:bg-gray-950 text-slate-800 dark:text-gray-200">
            {/* Sidebar */}
            <aside className="bg-white dark:bg-gray-900 
                border-r border-slate-100 dark:border-gray-800 flex-shrink-0
                sm:relative sm:w-64 fixed bottom-0 w-full">
                <nav className="mt-4 flex sm:flex-col flex-row space-y-2 text-slate-700 dark:text-gray-300">
                    {pages.map((page, i) => {
                        return <div
                            className="flex flex-row sm:justify-between justify-center sm:flex-0 flex-1 text-center px-4 py-2 hover:bg-slate-100 dark:hover:bg-gray-800 rounded"
                            onClick={() => setPage(page.name)}
                        >
                            <button className="sm:flex hidden" key={i}>{page.name}</button>
                            {page.icon}
                        </div>
                    })}
                </nav>
            </aside>

            <main className="flex-1 p-6 overflow-auto">
                {pages.map(page => {
                    if (page.name == current_page) return page.children
                    return null
                })}
            </main>
        </div>
    );
}
