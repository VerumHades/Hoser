// src/components/DashboardApp.tsx
import { AnimatePresence, motion } from "framer-motion";
import { ChevronDown, ChevronLeft } from "lucide-react";
import React, { createContext, useContext, useEffect, useState, type ReactNode } from "react";
import { useLocation, useNavigate } from "react-router";
import Navbar from "./Navbar";

interface DashboardContextType {
    setPage: (page: string) => void,
    page: string,
}

const DashboardContext = createContext<DashboardContextType>({ setPage: () => { }, page: "" })

interface PageProps {
    name: string;
    subpage_of?: string,
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
    logo: React.ReactNode;
}

function buildPageMap(pages: PageProps[]) {
    let hold: Record<string, any> = {

    }
    let map: Record<string, any> = {

    }

    for (let page of pages) {
        const { name, subpage_of } = page;
        if (subpage_of) {
            if (!(subpage_of in hold))
                throw Error(`Pages is a subpage of non existent page: ${subpage_of}. Did you declare it before this one?`)

            hold[name] = { name: page.name, icon: page.icon, children: {} }
            hold[subpage_of].children[name] = hold[name]
            continue
        }

        map[name] = { name: page.name, icon: page.icon, children: {} }
        hold[name] = map[name]
    }

    return { map, flatMap: hold }
}

interface DashboardLinkProps extends PageProps {
    cathegory: boolean,
    childMap: Record<string, any>
}

function hasNestedKeyValue(obj: any, key: string, value: any): boolean {
    if (obj && typeof obj === "object") {
        for (const k in obj) {
            if (!Object.prototype.hasOwnProperty.call(obj, k)) continue;

            if (k === key && obj[k] === value) return true;

            if (typeof obj[k] === "object") {
                if (k == "icon") continue;
                if (hasNestedKeyValue(obj[k], key, value)) return true;
            }
        }
    }
    return false;
}

function DashboardLink(props: DashboardLinkProps) {
    const dashcontext = useContext(DashboardContext);

    const hasActiveChild = hasNestedKeyValue(props.childMap, "name", dashcontext.page);
    const [open, setOpen] = useState<boolean>(hasActiveChild);

    const clickHandler = () => {
        if (!props.cathegory)
            dashcontext.setPage(props.name)
        else
            setOpen(!open)
    }


    return <div
        className="flex flex-col transition-all md:justify-between justify-center md:flex-0 flex-1 text-center"
    >
        <div className={"flex flex-row justify-between px-4 py-2 hover:bg-slate-100 dark:hover:bg-gray-800 transition-colors duration-300 " + (props.name == dashcontext.page ? "bg-slate-100 border-r-slate-200 dark:bg-slate-700 dark:border-r-indigo-600 border-r-4" : "")} onClick={clickHandler}>
            <button>{props.name}</button>
            {props.cathegory ? (open ? <ChevronDown /> : <ChevronLeft />) : props.icon}
        </div>

        <AnimatePresence>
            {open && (
                <motion.div
                    key="sidebar"
                    initial={{ y: -10, opacity: 0 }}
                    animate={{ y: 0, opacity: 1 }}
                    exit={{ y: -10, opacity: 0 }}
                    transition={{ type: "tween", duration: 0.3 }}
                    className="flex flex-col ml-5"
                >
                    {props.children}
                </motion.div>
            )}
        </AnimatePresence>
    </div>
}

interface PageLinksProps {
    linkMap: Record<string, any>,
}

function PageLinks({ linkMap }: PageLinksProps) {
    return Object.values(linkMap).map(({ name, icon, children }) => {
        return <DashboardLink childMap={children} icon={icon} name={name} cathegory={Object.values(children).length != 0}>
            <PageLinks linkMap={children} ></PageLinks>
        </DashboardLink>
    })
}
function flattenChildren(
    children: React.ReactNode
): React.ReactElement[] {
    const result: React.ReactElement[] = [];

    React.Children.forEach(children, (child) => {
        if (child === null || child === undefined || typeof child === "boolean") {
            return; // skip invisible nodes
        }

        if (
            React.isValidElement(child) &&
            child.type === React.Fragment
        ) {
            result.push(...flattenChildren((child as React.ReactElement<any>).props.children));
            return;
        }

        if (Array.isArray(child)) {
            result.push(...flattenChildren(child));
            return;
        }

        if (React.isValidElement(child)) {
            result.push(child);
            return;
        }
    });

    return result;
}
export default function Dashboard({ children }: DashboardProps) {
    const [navbar_open, setNavbarOpen] = useState<boolean>(false);

    const location = useLocation()
    const navigate = useNavigate()
    const current_page = location.state?.dashpage ?? ""

    const pages = flattenChildren(children).map((child) => {
        if (!React.isValidElement<PageProps>(child)) {
            throw new Error("<Dashboard> children must be <Page> components");
        }

        return child.props;
    })

    const { map: linkMap, flatMap } = buildPageMap(pages)

    useEffect(
        () => {
            if (!(current_page in flatMap) && Object.keys(flatMap).length != 0) {
                navigate(location.pathname, { state: { dashpage: Object.keys(flatMap)[0] } })
            }
        })

    const setPageHandler = (page: string) => {
        navigate(location.pathname, { state: { dashpage: page } })
        setNavbarOpen(false)
    }

    return (
        <DashboardContext.Provider value={{ setPage: setPageHandler, page: current_page }}>
            <div className="flex flex-col w-full h-full bg-white dark:bg-gray-900 pointer-events-auto">
                <Navbar className="z-60 bg-white dark:bg-gray-900 shadow-md w-full pointer-events-auto" open={navbar_open} onToggleOpen={() => setNavbarOpen(!navbar_open)}>
                    {navbar_open && (
                        <div className="space-y-2 text-slate-700 dark:text-gray-300 h-full overflow-y-auto">
                            <PageLinks linkMap={linkMap} />
                        </div>
                    )}
                </Navbar>
                <div className="flex md:flex-row flex-col flex-1 min-h-0 bg-slate-50 dark:bg-gray-950 text-slate-800 dark:text-gray-200">
                    <aside className="hidden md:flex md:flex-col md:w-64 md:relative md:border-r md:border-slate-100 md:bg-white dark:md:bg-gray-900 dark:md:border-gray-800">
                        <nav className="mt-4 flex flex-col space-y-2 text-slate-700 dark:text-gray-300 h-full overflow-y-auto">
                            <PageLinks linkMap={linkMap} />
                        </nav>
                    </aside>

                    <main className="flex-1 min-h-0">
                        {pages.map(page => {
                            if (page.name == current_page) return page.children
                            return null
                        })}
                    </main>
                </div>
            </div>
        </DashboardContext.Provider>
    );
}
