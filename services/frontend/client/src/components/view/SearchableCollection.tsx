import { useState, useMemo, useEffect, type ReactNode } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { CursorPaginatedCollection } from "./CursorPaginatedCollection";

interface SearchableCollectionProps<T, Q> {
	title?: string;
	description?: string;
	headerActions?: ReactNode;
	emptyState?: ReactNode;
	initialQuery: Q;
	fetchPage: (query: Q, cursor?: string) => Promise<any>;
	renderRow: (item: T) => ReactNode;
	renderCard: (item: T) => ReactNode;
	parseParams: (params: URLSearchParams) => Q;
	buildParams: (query: Q) => URLSearchParams;
	FilterComponent?: React.ComponentType<{
		query: Q;
		onApply: (query: Q) => void;
	}>;
}

/**
 * A standardized high-level wrapper for paginated lists with search and filter capabilities.
 */
export function SearchableCollection<T, Q>(props: SearchableCollectionProps<T, Q>) {
	const [searchParams, setSearchParams] = useSearchParams();
	const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

	const activeQuery = useMemo(
		() => props.parseParams(searchParams),
		[searchParams, props.parseParams]
	);

	const handleUpdate = (newFields: Partial<Q>) => {
		const nextQuery = { ...activeQuery, ...newFields };
		setSearchParams(props.buildParams(nextQuery));
		setIsMobileMenuOpen(false);
	};

	return (
		<div className="flex flex-col h-full w-full bg-white overflow-hidden">
			<CollectionHeader
				title={props.title}
				description={props.description}
				actions={props.headerActions}
				onOpenFilters={() => setIsMobileMenuOpen(true)}
				hasFilters={!!props.FilterComponent}
			/>

			<div className="flex flex-1 overflow-hidden max-w-[1600px] mx-auto w-full">
				<SidebarArea
					query={activeQuery}
					onApply={handleUpdate}
					FilterComponent={props.FilterComponent}
				/>

				<MobileDrawerArea
					isOpen={isMobileMenuOpen}
					query={activeQuery}
					onApply={handleUpdate}
					onClose={() => setIsMobileMenuOpen(false)}
					FilterComponent={props.FilterComponent}
				/>

				<MainContentArea
					query={activeQuery}
					searchParams={searchParams}
					fetchPage={props.fetchPage}
					renderRow={props.renderRow}
					renderCard={props.renderCard}
					emptyState={props.emptyState}
				/>
			</div>
		</div>
	);
}

/**
 * Renders the top bar with titles and actions.
 */
function CollectionHeader({
	title,
	description,
	actions,
	onOpenFilters,
	hasFilters,
}: {
	title?: string;
	description?: string;
	actions?: ReactNode;
	onOpenFilters: () => void;
	hasFilters: boolean;
}) {
	if (!title && !actions) return null;

	return (
		<header className="flex-shrink-0 border-b border-slate-100 px-6 py-6 bg-white sticky top-0 z-30">
			<div className="flex justify-between items-end max-w-[1600px] mx-auto w-full">
				<div>
					{title && <h1 className="text-2xl font-bold text-slate-900">{title}</h1>}
					{description && <p className="text-slate-500 text-sm">{description}</p>}
				</div>
				<div className="flex gap-2">
					{actions}
					{hasFilters && (
						<button
							onClick={onOpenFilters}
							className="md:hidden p-2 bg-slate-100 rounded-lg"
						>
							Filters
						</button>
					)}
				</div>
			</div>
		</header>
	);
}

/**
 * Handles the desktop sidebar visibility logic.
 */
function SidebarArea<Q>({
	query,
	onApply,
	FilterComponent,
}: {
	query: Q;
	onApply: (q: Q) => void;
	FilterComponent?: React.ComponentType<{ query: Q; onApply: (q: Q) => void }>;
}) {
	if (!FilterComponent) return null;

	return (
		<aside className="hidden md:block w-80 h-full border-r border-slate-100 overflow-y-auto bg-slate-50/30">
			<FilterComponent query={query} onApply={onApply} />
		</aside>
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
	FilterComponent,
}: {
	isOpen: boolean;
	query: Q;
	onApply: (q: Q) => void;
	onClose: () => void;
	FilterComponent?: React.ComponentType<{ query: Q; onApply: (q: Q) => void }>;
}) {
	if (!isOpen || !FilterComponent) return null;

	return (
		<div className="fixed inset-0 z-50 md:hidden">
			<div className="absolute inset-0 bg-slate-900/60 backdrop-blur-sm" onClick={onClose} />
			<aside className="absolute right-0 top-0 h-full w-80 bg-white shadow-2xl flex flex-col overflow-y-auto">
				<div className="p-6 border-b flex justify-between items-center">
					<span className="font-black">Filters</span>
					<button onClick={onClose}>✕</button>
				</div>
				<FilterComponent query={query} onApply={onApply} />
			</aside>
		</div>
	);
}

/**
 * The scrolling area for the actual collection items.
 */
function MainContentArea<T, Q>({
	query,
	searchParams,
	fetchPage,
	renderRow,
	renderCard,
	emptyState,
}: {
	query: Q;
	searchParams: URLSearchParams;
	fetchPage: (q: Q, cursor?: string) => Promise<any>;
	renderRow: (item: T) => ReactNode;
	renderCard: (item: T) => ReactNode;
	emptyState?: ReactNode;
}) {
	return (
		<main className="flex-1 h-full overflow-y-auto scroll-smooth">
			<div className="px-6 py-8 md:px-12">
				<CursorPaginatedCollection<T>
					key={searchParams.toString()}
					fetchPage={(cursor) => fetchPage(query, cursor)}
					renderItemRow={renderRow}
					renderItemCard={renderCard}
					emptyState={emptyState}
				/>
			</div>
		</main>
	);
}