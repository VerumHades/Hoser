import React, { useState, type ReactNode } from "react";
import { useLocation, useNavigate, Routes, Route } from "react-router";
import Navbar from "./Navbar";

interface PageProps {
    name: string;
    icon?: ReactNode;
    children?: ReactNode;
    path?: string; // optional, allows subpaths like "billing/*"
}

export const Page: React.FC<PageProps> = () => null;

interface DashboardProps {
    children: ReactNode;
}

interface DashboardLinkProps {
    name: string;
    icon?: ReactNode;
}

function DashboardLink({ name, icon }: DashboardLinkProps) {
    const navigate = useNavigate();
    const location = useLocation();
    const isActive = location.pathname === `/dashboard/${name.toLowerCase()}`;

    const handleClick = () => {
        navigate(`/dashboard/${name.toLowerCase()}`);
    };

    return (
        <button
            onClick={handleClick}
            className={`
                w-full flex items-center gap-3
                px-4 py-2 text-sm rounded-lg transition
                ${isActive
                    ? "bg-slate-100 dark:bg-gray-800 text-slate-900 dark:text-gray-100"
                    : "text-slate-600 dark:text-gray-400 hover:bg-slate-100 dark:hover:bg-gray-800"
                }
            `}
        >
            {icon && <span className="w-4 h-4">{icon}</span>}
            <span className="truncate">{name}</span>
        </button>
    );
}

function flattenPages(children: React.ReactNode): PageProps[] {
    const pages: PageProps[] = [];

    React.Children.forEach(children, (child) => {
        if (!React.isValidElement<PageProps>(child)) return;

        if (child.type === React.Fragment) {
            pages.push(...flattenPages(child.props.children));
        } else {
            pages.push(child.props);
        }
    });

    return pages;
}

export default function Dashboard({ children }: DashboardProps) {
    const [navbarOpen, setNavbarOpen] = useState(false);
    const pages = flattenPages(children);

    return (
        <div className="flex flex-col w-full h-full bg-slate-50 dark:bg-gray-950">
            <Navbar
                open={navbarOpen}
                onToggleOpen={() => setNavbarOpen(!navbarOpen)}
            >
                <div className="flex flex-col gap-1">
                    {pages.map((page) => (
                        <DashboardLink
                            key={page.name}
                            name={page.name}
                            icon={page.icon}
                        />
                    ))}
                </div>
            </Navbar>

            <div className="flex flex-1 min-h-0">
                <aside className="
                        hidden md:flex md:flex-col md:w-64
                        border-r border-slate-200 dark:border-gray-800
                        bg-white dark:bg-gray-900
                        p-3
                    ">
                    <nav className="flex flex-col gap-1">
                        {pages.map((page) => (
                            <DashboardLink
                                key={page.name}
                                name={page.name}
                                icon={page.icon}
                            />
                        ))}
                    </nav>
                </aside>

                <main className="flex-1 min-w-0 overflow-hidden p-4">
                    <Routes>
                        {pages.map((page) => (
                            <Route
                                key={page.name}
                                path={page.path ?? page.name.toLowerCase() + "/*"}
                                element={page.children}
                            />
                        ))}
                    </Routes>

                </main>
            </div>
        </div>
    );
}
