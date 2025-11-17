// src/components/DashboardApp.tsx
import { AnimatePresence, motion } from "framer-motion";
import { ChevronDown, ChevronLeft, Menu } from "lucide-react";
import React, { Children, useEffect, useState, type ReactNode } from "react";
import { useLocation } from "react-router";
import Navbar from "../Navbar";

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

            hold[name] = { props: page, children: {} }
            hold[subpage_of].children[name] = hold[name]
            continue
        }

        map[name] = { props: page, children: {} }
        hold[name] = map[name]
    }

    return { map, flatMap: hold }
}

interface DashboardLinkProps extends PageProps {
    setPage: (name: string) => void,
    cathegory: boolean
}
function DashboardLink(props: DashboardLinkProps) {
    const [open, setOpen] = useState<boolean>(false);

    const clickHandler = () => {
        if (!props.cathegory)
            props.setPage(props.name)
        else
            setOpen(!open)
    }

    return <div
        className="flex flex-col transition-all md:justify-between justify-center md:flex-0 flex-1 text-center ml-5"
    >
        <div className="flex flex-row justify-between  px-4 py-2 hover:bg-slate-100 dark:hover:bg-gray-800 rounded" onClick={clickHandler}>
            <button>{props.name}</button>
            {props.cathegory ? (open ? <ChevronDown /> : <ChevronLeft />) : props.icon}
        </div>

        <AnimatePresence>
            {open && (
                <motion.aside
                    key="sidebar"
                    initial={{ y: -10, opacity: 0 }}
                    animate={{ y: 0, opacity: 1 }}
                    exit={{ y: -10, opacity: 0 }}
                    transition={{ type: "tween", duration: 0.3 }}
                    className="flex flex-col"
                >
                    {props.children}
                </motion.aside>
            )}
        </AnimatePresence>
    </div>
}

interface PageLinksProps {
    linkMap: Record<string, any>,
    setPage: (name: string) => void
}

function PageLinks({ linkMap, setPage }: PageLinksProps) {
    return Object.values(linkMap).map(({ props, children }) => {
        return <DashboardLink icon={props.icon} name={props.name} cathegory={Object.values(children).length != 0} setPage={setPage}>
            <PageLinks linkMap={children} setPage={setPage} ></PageLinks>
        </DashboardLink>
    })
}
export function flattenChildren(
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

export default function Dashboard({ children, logo }: DashboardProps) {
    const [current_page, setPage] = useState<string>("listings");
    const [open, setOpen] = useState<boolean>(false)

    const pages = flattenChildren(children).map((child) => {
        if (!React.isValidElement<PageProps>(child)) {
            throw new Error("<Dashboard> children must be <Page> components");
        }

        return child.props;
    })
    const { map: linkMap, flatMap } = buildPageMap(pages)

    const location = useLocation()

    useEffect(() => {
        if (location.state?.dashpage) {
            const name = location.state?.dashpage
            if (name in flatMap) {
                setPage(name)
                return;
            }
        }
        setPage(pages[0].name)
    }, [])

    return (
        <div className="flex flex-col w-full h-full bg-white dark:bg-gray-900 pointer-events-auto">
            <Navbar>
                <div className="text-indigo-700 dark:text-gray-300 px-4 py-2 rounded transition-colors">
                    <PageLinks linkMap={linkMap} setPage={setPage} />
                </div>
            </Navbar>
            <div className="flex md:flex-row flex-col flex-1 bg-slate-50 dark:bg-gray-950 text-slate-800 dark:text-gray-200">
                <aside className="hidden md:flex md:flex-col md:w-64 md:relative md:border-r md:border-slate-100 md:bg-white dark:md:bg-gray-900 dark:md:border-gray-800">
                    <nav className="mt-4 flex flex-col space-y-2 text-slate-700 dark:text-gray-300 h-full overflow-y-auto">
                        <PageLinks linkMap={linkMap} setPage={setPage} />
                    </nav>
                </aside>

                <main className="flex-1 p-6 overflow-auto">
                    {pages.map(page => {
                        if (page.name == current_page) return page.children
                        return null
                    })}
                </main>
            </div>
        </div>
    );
}
