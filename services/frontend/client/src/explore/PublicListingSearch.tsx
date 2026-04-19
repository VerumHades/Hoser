import { useState, useMemo, useEffect } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { ListingAPI, type Listing } from "../backend/repositories/listing";
import { PublicListingRow } from "./List/PublicListingRow";
import { PublicListingCard } from "./List/PublicListingCard";
import { SearchableCollection } from "../components/view/SearchableCollection";
import type { DateRange, ListingSearchQuery, NumericRange } from "../backend/utils/query";

export default function PublicListingSearch() {
	const navigate = useNavigate();

	return (
		<SearchableCollection<Listing, ListingSearchQuery>
			initialQuery={{}}
			fetchPage={(q, cursor) => ListingAPI.search(q, cursor)}
			parseParams={parseParamsToQuery}
			buildParams={buildUrlParams}
			FilterComponent={ListingFilterSidebar}
			renderRow={(l) => <PublicListingRow listing={l} onSelect={() => navigate(`/listing/${l.id}`)} />}
			renderCard={(l) => <PublicListingCard listing={l} onSelect={() => navigate(`/listing/${l.id}`)} />}
		/>
	);
}

/**
 * ListingFilterSidebar contains the guarded numeric and date inputs.
 */
export function ListingFilterSidebar({ query, onApply }: { query: ListingSearchQuery; onApply: (q: ListingSearchQuery) => void }) {
	const [local, setLocal] = useState(query || {});

	useEffect(() => setLocal(query), [query]);

	const isValid = useMemo(() => validateRange(local.cpu) && validateRange(local.ramBytes) && validateRange(local.diskBytes) && validateRange(local.price), [local]);

	return (
		<div className="p-8 flex flex-col gap-10" onKeyDown={(e) => e.key === "Enter" && isValid && onApply(local)}>
			<RangeInput title="CPU Cores" range={local.cpu || {}} onChange={(cpu) => setLocal({ ...local, cpu })} />
			<RangeInput title="RAM (Bytes)" range={local.ramBytes || {}} onChange={(ramBytes) => setLocal({ ...local, ramBytes })} />
			<RangeInput title="Disk (Bytes)" range={local.diskBytes || {}} onChange={(diskBytes) => setLocal({ ...local, diskBytes })} />
			<RangeInput title="Price" range={local.price || {}} onChange={(price) => setLocal({ ...local, price })} />
			<DateInput title="Creation Date" range={local.createdAt || {}} onChange={(createdAt) => setLocal({ ...local, createdAt })} />

			<div className="p-4 bg-amber-50 border border-amber-100 rounded-xl">
				<label className="text-[10px] font-black uppercase text-amber-600 block mb-2">Author ID (Temp)</label>
				<input type="text" className="w-full bg-white border border-amber-200 p-2 text-xs rounded shadow-inner" value={local.authorId || ""} onChange={(e) => setLocal({ ...local, authorId: e.target.value })} />
			</div>

			<button disabled={!isValid} onClick={() => onApply(local)} className="w-full bg-slate-900 text-white py-4 rounded-2xl font-bold shadow-xl active:scale-95 transition-all disabled:opacity-30">Apply Filters</button>
		</div>
	);
}

/**
 * RangeInput handles numeric guards for min/max logic.
 */
function RangeInput({ title, range, onChange }: { title: string; range: NumericRange; onChange: (r: NumericRange) => void }) {
	const error = range.min !== undefined && range.max !== undefined && range.min > range.max;

	return (
		<section className="flex flex-col gap-3">
			<h4 className="text-[11px] font-black text-slate-400 uppercase tracking-widest">{title}</h4>
			<div className={`flex items-center gap-2 p-1 rounded-lg transition-colors ${error ? "bg-red-50" : ""}`}>
				<input type="number" placeholder="Min" className="w-full bg-white border border-slate-200 p-2 rounded-md text-sm outline-none focus:border-indigo-500" value={range.min ?? ""} onChange={(e) => onChange({ ...range, min: e.target.value ? parseInt(e.target.value) : undefined })} />
				<span className="text-slate-300">/</span>
				<input type="number" placeholder="Max" className="w-full bg-white border border-slate-200 p-2 rounded-md text-sm outline-none focus:border-indigo-500" value={range.max ?? ""} onChange={(e) => onChange({ ...range, max: e.target.value ? parseInt(e.target.value) : undefined })} />
			</div>
		</section>
	);
}

function DateInput({ title, range, onChange }: { title: string; range: DateRange; onChange: (r: DateRange) => void }) {
	return (
		<section className="flex flex-col gap-3">
			<h4 className="text-[11px] font-black text-slate-400 uppercase tracking-widest">{title}</h4>
			<input type="date" className="w-full border border-slate-200 p-2 rounded-md text-sm mb-2" value={range.from || ""} onChange={(e) => onChange({ ...range, from: e.target.value })} />
			<input type="date" className="w-full border border-slate-200 p-2 rounded-md text-sm" value={range.to || ""} onChange={(e) => onChange({ ...range, to: e.target.value })} />
		</section>
	);
}
/**
 * Utility: URL Serialization
 */
function parseParamsToQuery(p: URLSearchParams): ListingSearchQuery {
	const parse = (k: string) => (p.get(k) ? parseInt(p.get(k)!, 10) : undefined);
	return {
		text: p.get("q") || undefined,
		cpu: { min: parse("cpu_min"), max: parse("cpu_max") },
		ramBytes: { min: parse("ram_min"), max: parse("ram_max") },
		diskBytes: { min: parse("disk_min"), max: parse("disk_max") },
		price: { min: parse("price_min"), max: parse("price_max") },
		createdAt: { from: p.get("date_from") || undefined, to: p.get("date_to") || undefined },
		authorId: p.get("author_id") || undefined,
	};
}

function buildUrlParams(q: ListingSearchQuery): URLSearchParams {
	const p = new URLSearchParams();
	if (q.text) p.set("q", q.text);
	const add = (k: string, v?: number) => v !== undefined && p.set(k, v.toString());
	add("cpu_min", q.cpu?.min); add("cpu_max", q.cpu?.max);
	add("ram_min", q.ramBytes?.min); add("ram_max", q.ramBytes?.max);
	add("disk_min", q.diskBytes?.min); add("disk_max", q.diskBytes?.max);
	add("price_min", q.price?.min); add("price_max", q.price?.max);
	if (q.createdAt?.from) p.set("date_from", q.createdAt.from);
	if (q.createdAt?.to) p.set("date_to", q.createdAt.to);
	if (q.authorId) p.set("author_id", q.authorId);
	return p;
}

function validateRange(r?: NumericRange) {
	return !r || r.min === undefined || r.max === undefined || r.min <= r.max;
}