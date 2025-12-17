import React, { createContext, useContext, useState, type ReactNode } from "react";
import { useLocation, useNavigate } from "react-router";
import Navbar from "./Navbar";

interface DashboardContextType {
    page: string;
    setPage: (page: string) => void;
}

const DashboardContext = createContext<DashboardContextType>({
    page: "",
    setPage: () => {},
});

interface PageProps {
    name: string;
    icon?: ReactNode;
    children?: ReactNode; // optional content for the main area
}

export const Page: React.FC<PageProps> = () => null;

interface DashboardProps {
    children: ReactNode;
    logo?: ReactNode;
}

interface DashboardLinkProps {
    name: string;
    icon?: ReactNode;
}

/** Simple flat link without children */
function DashboardLink({ name, icon }: DashboardLinkProps) {
    const { page, setPage } = useContext(DashboardContext);
    const navigate = useNavigate();
    const location = useLocation();

    const handleClick = () => {
        navigate(location.pathname, { state: { dashpage: name } });
        setPage(name);
    };

    return (
        <div
            className={`flex justify-start items-center px-4 py-2 cursor-pointer hover:bg-slate-100 dark:hover:bg-gray-800 transition-colors duration-300 ${
                page === name
                    ? "bg-slate-100 dark:bg-slate-700 border-r-4 border-r-slate-200 dark:border-r-indigo-600"
                    : ""
            }`}
            onClick={handleClick}
        >
            {icon && <span className="mr-2">{icon}</span>}
            <span>{name}</span>
        </div>
    );
}

/** Flatten <Page> children into PageProps */
function flattenPages(children: React.ReactNode): PageProps[] {
    const result: PageProps[] = [];

    React.Children.forEach(children, (child) => {
        if (!React.isValidElement<PageProps>(child)) return;
        if (child.type === React.Fragment) {
            result.push(...flattenPages(child.props.children));
        } else {
            result.push(child.props);
        }
    });

    return result;
}

export default function Dashboard({ children, logo }: DashboardProps) {
    const location = useLocation();
    const currentPage = location.state?.dashpage ?? "";

    const [navbarOpen, setNavbarOpen] = useState<boolean>(false);

    const setPageHandler = (page: string) => {
        setNavbarOpen(false);
    };

    const pages = flattenPages(children);

    return (
        <DashboardContext.Provider value={{ page: currentPage, setPage: setPageHandler }}>
            <div className="flex flex-col w-full h-full bg-white dark:bg-gray-900 pointer-events-auto">
                <Navbar
                    className="z-60 bg-white dark:bg-gray-900 shadow-md w-full pointer-events-auto"
                    open={navbarOpen}
                    onToggleOpen={() => setNavbarOpen(!navbarOpen)}
                >
                    {navbarOpen && (
                        <div className="space-y-2 text-slate-700 dark:text-gray-300 h-full overflow-y-auto">
                            {pages.map((page) => (
                                <DashboardLink key={page.name} name={page.name} icon={page.icon} />
                            ))}
                        </div>
                    )}
                </Navbar>

                <div className="flex md:flex-row flex-col flex-1 min-h-0 bg-slate-50 dark:bg-gray-950 text-slate-800 dark:text-gray-200">
                    <aside className="hidden md:flex md:flex-col md:w-64 md:border-r md:border-slate-100 md:bg-white dark:md:bg-gray-900 dark:md:border-gray-800">
                        <nav className="mt-4 flex flex-col space-y-2 text-slate-700 dark:text-gray-300 h-full overflow-y-auto">
                            {pages.map((page) => (
                                <DashboardLink key={page.name} name={page.name} icon={page.icon} />
                            ))}
                        </nav>
                    </aside>

                    <main className="flex flex-1 min-w-0 overflow-hidden">
                        {pages.map((page) => (page.name === currentPage ? page.children : null))}
                    </main>
                </div>
            </div>
        </DashboardContext.Provider>
    );
}
