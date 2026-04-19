import { useState, useMemo, type ReactNode } from "react";
import { useSearchParams } from "react-router-dom";
import { CursorPaginatedCollection } from "./CursorPaginatedCollection";
import { DynamicFilterSidebar, type FilterField } from "../querying/FilterSidebar";

interface SearchableCollectionProps<T, Q> {
    filterFields?: FilterField<Q>[];
    headerActions?: ReactNode;
    useUrlParams?: boolean;
    initialQuery?: Q;
    fetchPage: (query: Q, cursor?: string) => Promise<any>;
    renderRow: (item: T) => ReactNode;
    renderCard: (item: T) => ReactNode;
    title?: string;
    description?: string;
    emptyState?: ReactNode;
}

/**
 * A standardized high-level wrapper for paginated lists with search and filter capabilities.
 */
export function SearchableCollection<T, Q>(props: SearchableCollectionProps<T, Q>) {
    const { filterFields = [], useUrlParams = true, initialQuery = {} as Q } = props;
    const [searchParameters, setSearchParameters] = useSearchParams();
    const [localQuery, setLocalQuery] = useState<Q>(initialQuery);
    const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

    const activeQuery = useMemo(() => {
        return useUrlParams ? parseUrlToQuery(searchParameters, filterFields) : localQuery;
    }, [searchParameters, localQuery, filterFields, useUrlParams]);

    const handleQueryUpdate = (nextQuery: Q) => {
        const nextParameters = buildUrlFromQuery(nextQuery, filterFields);
        if (useUrlParams) {
            setSearchParameters(nextParameters);
        } else {
            setLocalQuery(nextQuery);
        }
        setIsMobileMenuOpen(false);
    };

    return (
        <div className="flex flex-col h-full w-full bg-white overflow-hidden md:flex-row">
            <MobileTopBar 
                title={props.title} 
                onOpenFilters={() => setIsMobileMenuOpen(true)} 
            />
            
            <DesktopSidebar 
                title={props.title}
                description={props.description}
                actions={props.headerActions}
                query={activeQuery}
                fields={filterFields}
                onApply={handleQueryUpdate}
            />

            <MobileDrawerArea 
                isOpen={isMobileMenuOpen} 
                query={activeQuery} 
                onApply={handleQueryUpdate} 
                onClose={() => setIsMobileMenuOpen(false)} 
                fields={filterFields} 
            />

            <MainContentArea 
                query={activeQuery} 
                searchParamsKey={useUrlParams ? searchParameters.toString() : JSON.stringify(activeQuery)}
                {...props} 
            />
        </div>
    );
}

/**
 * Renders a sticky top bar visible only on mobile devices.
 */
function MobileTopBar({ title, onOpenFilters }: { title?: string; onOpenFilters: () => void }) {
    return (
        <header className="md:hidden flex items-center justify-between px-6 py-4 border-b border-slate-100 bg-white sticky top-0 z-30">
            <h1 className="font-bold text-lg truncate">{title}</h1>
            <button 
                type="button"
                onClick={onOpenFilters} 
                className="p-2 bg-slate-100 rounded-lg text-sm font-medium"
            >
                Filters
            </button>
        </header>
    );
}

/**
 * Renders the sidebar containing the header info and filters for desktop.
 */
function DesktopSidebar<Q>({
    title,
    description,
    actions,
    query,
    fields,
    onApply
}: {
    title?: string;
    description?: string;
    actions?: ReactNode;
    query: Q;
    fields: FilterField<Q>[];
    onApply: (q: Q) => void;
}) {
    return (
        <aside className="hidden md:flex flex-col w-80 h-full border-r border-slate-100 bg-slate-50/30 overflow-y-auto">
            <SidebarHeader title={title} description={description} />
            <SidebarActions actions={actions} />
            <DynamicFilterSidebar query={query} onApply={onApply} fields={fields} />
        </aside>
    );
}

/**
 * Renders the title and description inside the sidebar.
 */
function SidebarHeader({ title, description }: { title?: string; description?: string }) {
    if (!title && !description) return null;

    return (
        <div className="px-6 pt-8 pb-4">
            {title && <h1 className="text-xl font-bold text-slate-900">{title}</h1>}
            {description && <p className="mt-1 text-sm text-slate-500">{description}</p>}
        </div>
    );
}

/**
 * Renders custom action buttons inside the sidebar.
 */
function SidebarActions({ actions }: { actions?: ReactNode }) {
    if (!actions) return null;

    return (
        <div className="px-6 py-4 flex flex-wrap gap-2">
            {actions}
        </div>
    );
}

/**
 * The scrolling area for the actual collection items.
 */
function MainContentArea<T, Q>({
    query,
    searchParamsKey,
    fetchPage,
    renderRow,
    renderCard,
    emptyState,
}: {
    query: Q;
    searchParamsKey: string;
    fetchPage: (q: Q, cursor?: string) => Promise<any>;
    renderRow: (item: T) => ReactNode;
    renderCard: (item: T) => ReactNode;
    emptyState?: ReactNode;
}) {
    return (
        <main className="flex flex-1 h-full overflow-y-auto scroll-smooth">
			<CursorPaginatedCollection<T>
				key={searchParamsKey}
				fetchPage={(cursor) => fetchPage(query, cursor)}
				renderItemRow={renderRow}
				renderItemCard={renderCard}
				emptyState={emptyState}
			/>
        </main>
    );
}

/**
 * Manages the mobile drawer overlay.
 */
function MobileDrawerArea<Q>({
    isOpen,
    query,
    onApply,
    onClose,
    fields
}: {
    isOpen: boolean;
    query: Q;
    onApply: (q: Q) => void;
    onClose: () => void;
    fields: FilterField<Q>[] 
}) {
    if (!isOpen) return null;

    return (
        <div className="fixed inset-0 z-50 md:hidden">
            <div 
                className="absolute inset-0 bg-slate-900/60 backdrop-blur-sm" 
                onClick={onClose} 
            />
            <aside className="absolute right-0 top-0 h-full w-80 bg-white shadow-2xl flex flex-col overflow-y-auto">
                <DrawerHeader onClose={onClose} />
                <DynamicFilterSidebar query={query} onApply={onApply} fields={fields} />
            </aside>
        </div>
    );
}

/**
 * Renders the header for the mobile drawer.
 */
function DrawerHeader({ onClose }: { onClose: () => void }) {
    return (
        <div className="p-6 border-b flex justify-between items-center">
            <span className="font-black text-slate-900">Filters</span>
            <button type="button" onClick={onClose} className="p-2 text-slate-500 hover:text-slate-900">✕</button>
        </div>
    );
}

/**
 * Transforms URL parameters into a structured query object.
 */
function parseUrlToQuery<Q>(searchParameters: URLSearchParams, fields: FilterField<Q>[]): Q {
    const query = {} as any;
    fields.forEach((field) => {
        if (field.type === "text") {
            query[field.key] = searchParameters.get(field.urlKey) || undefined;
        }
        if (field.type === "range") {
            const minimum = searchParameters.get(field.urlKeys.min);
            const maximum = searchParameters.get(field.urlKeys.max);
            query[field.key] = { 
                min: minimum ? parseInt(minimum) : undefined, 
                max: maximum ? parseInt(maximum) : undefined 
            };
        }
        if (field.type === "date") {
            query[field.key] = { 
                from: searchParameters.get(field.urlKeys.from) || undefined, 
                to: searchParameters.get(field.urlKeys.to) || undefined 
            };
        }
    });
    return query;
}

/**
 * Transforms a query object into URL parameters.
 */
function buildUrlFromQuery<Q>(query: Q, fields: FilterField<Q>[]): URLSearchParams {
    const params = new URLSearchParams();
    fields.forEach((field) => {
        const value = query[field.key] as any;
        if (!value) return;

        if (field.type === "text") params.set(field.urlKey, value);
        if (field.type === "range") {
            if (value.min !== undefined) params.set(field.urlKeys.min, String(value.min));
            if (value.max !== undefined) params.set(field.urlKeys.max, String(value.max));
        }
        if (field.type === "date") {
            if (value.from) params.set(field.urlKeys.from, value.from);
            if (value.to) params.set(field.urlKeys.to, value.to);
        }
    });
    return params;
}