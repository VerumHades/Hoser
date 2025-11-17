// src/components/DashboardApp.tsx
import { ChevronDown, ChevronLeft } from "lucide-react";
import React, { Children, useEffect, useState, type ReactNode } from "react";

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

            hold[name] = {props: page, children: {}}
            hold[subpage_of].children[name] = hold[name]
            continue
        }

        map[name] = {props: page, children: {}}
        hold[name] = map[name]
    }

    return map
}

function hasChildren(node: React.ReactNode): boolean {
    if (!React.isValidElement(node)) return false;
    const element = node as React.ReactElement<any>;
    return React.Children.count(element.props.children) > 0;
}

interface DashboardLinkProps extends PageProps{
    setPage: (name: string) => void,
    cathegory: boolean
}
function DashboardLink(props: DashboardLinkProps) {
    const [open, setOpen] = useState<boolean>(false);

    const clickHandler = () => {
        if(!props.cathegory)
            props.setPage(props.name)
        else
            setOpen(!open)
    }

    return <div
        className="flex flex-col transition-all md:justify-between justify-center md:flex-0 flex-1 text-center"
    >   
        <div className="flex flex-row justify-between  px-4 py-2 hover:bg-slate-100 dark:hover:bg-gray-800 rounded" onClick={clickHandler}>
            <button className="md:flex hidden">{props.name}</button>
            {props.cathegory ? (open ? <ChevronDown /> : <ChevronLeft />) : props.icon}
        </div>
        
        {open ? <div className="flex flex-col">{props.children}</div> : <></>}
    </div>
}

interface PageLinksProps{
    linkMap: Record<string, any>,
    setPage: (name: string) => void
}

function PageLinks({linkMap, setPage}: PageLinksProps) {
    return Object.values(linkMap).map(({props, children}) => {
        return <DashboardLink icon={props.icon} name={props.name} cathegory={Object.values(children).length != 0} setPage={setPage}>
            <PageLinks linkMap={children} setPage={setPage} ></PageLinks>
        </DashboardLink>
    })
}

export default function Dashboard({ children, logo }: DashboardProps) {
    const [current_page, setPage] = useState<string>("listings");

    const pages = React.Children.toArray(children).map((child) => {
        if (!React.isValidElement<PageProps>(child)) {
            throw new Error("<Dashboard> children must be <Page> components");
        }

        return child.props;
    })
    const linkMap = buildPageMap(pages)

    useEffect(() => {
        setPage(pages[0].name)
    }, [])

    return (
        <div className="h-full bg-white dark:bg-gray-900 shadow-md w-full pointer-events-auto">
            <div className="flex w-full flex-row items-center pt-4">
                {logo}
            </div>
            <div className="flex md:flex-row flex-col h-full bg-slate-50 dark:bg-gray-950 text-slate-800 dark:text-gray-200">
                {/* Sidebar */}
                <aside className="bg-white dark:bg-gray-900 
                border-r border-slate-100 dark:border-gray-800 flex-shrink-0
                md:relative md:w-64 md:order-1 w-full order-2">
                    <nav className="mt-4 flex md:flex-col flex-row space-y-2 text-slate-700 dark:text-gray-300">
                        <PageLinks linkMap={linkMap} setPage={setPage}></PageLinks>
                    </nav>
                </aside>

                <main className="flex-1 p-6 overflow-auto md:order-2 order-1">
                    {pages.map(page => {
                        if (page.name == current_page) return page.children
                        return null
                    })}
                </main>
            </div>
        </div>
    );
}
