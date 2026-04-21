import { useState, useEffect } from "react";
import type { DateRange, NumericRange } from "../../backend/utils/query";

export type FilterField<Q> = {
	key: keyof Q;
	label: string;
} & (
	| { type: "text"; urlKey: string }
    | { type: "select"; urlKey: string }
	| { type: "date"; urlKeys: { from: string; to: string } }
	| { type: "range"; urlKeys: { min: string; max: string } }
);

interface DynamicFilterSidebarProps<Q> {
    query: Q;
    onApply: (query: Q) => void;
    fields: FilterField<Q>[];
}

/**
 * Renders a dynamic set of filter inputs based on a configuration schema.
 */
export function DynamicFilterSidebar<Q>({ query, onApply, fields }: DynamicFilterSidebarProps<Q>) {
    const [local, setLocal] = useState<Q>(query);

    useEffect(() => setLocal(query), [query]);

    const handleFieldChange = (key: keyof Q, value: unknown) => {
        setLocal((prev) => ({ ...prev, [key]: value }));
    };

    return (
        <div 
            className="p-8 flex flex-col gap-10" 
            onKeyDown={(e) => e.key === "Enter" && onApply(local)}
        >
            {fields.map((field) => (
                <FilterFieldSwitch
                    key={String(field.key)}
                    config={field}
                    value={local[field.key]}
                    onChange={(val) => handleFieldChange(field.key, val)}
                />
            ))}

            <button 
                onClick={() => onApply(local)} 
                className="w-full bg-slate-900 text-white py-4 rounded-2xl font-bold shadow-xl active:scale-95 transition-all"
            >
                Apply Filters
            </button>
        </div>
    );
}

/**
 * Determines which input type to render based on configuration.
 */
function FilterFieldSwitch({ config, value, onChange }: { 
    config: FilterField<any>; 
    value: any; 
    onChange: (val: any) => void 
}) {
    switch (config.type) {
        case "range":
            return <RangeInput title={config.label} range={value || {}} onChange={onChange} />;
        case "date":
            return <DateInput title={config.label} range={value || {}} onChange={onChange} />;
        case "text":
            return <TextInput title={config.label} value={value || ""} onChange={onChange} />;
        default:
            return null;
    }
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


function TextInput({ title, value, onChange }: { title: string; value: string; onChange: (v: string) => void }) {
    return (
        <section className="flex flex-col gap-3">
            <h4 className="text-[11px] font-black text-slate-400 uppercase tracking-widest">{title}</h4>
            <input 
                type="text" 
                className="w-full bg-slate-50 border border-slate-200 p-2 text-sm rounded-md" 
                value={value} 
                onChange={(e) => onChange(e.target.value)} 
            />
        </section>
    );
}

function RangeInput({ title, range, onChange }: { title: string; range: NumericRange; onChange: (r: NumericRange) => void }) {
    return (
        <section className="flex flex-col gap-3">
            <h4 className="text-[11px] font-black text-slate-400 uppercase tracking-widest">{title}</h4>
            <div className="flex items-center gap-2">
                <input 
                    type="number" 
                    placeholder="Min" 
                    className="w-full border border-slate-200 p-2 rounded-md text-sm" 
                    value={range.min ?? ""} 
                    onChange={(e) => onChange({ ...range, min: e.target.value ? parseInt(e.target.value) : undefined })} 
                />
                <input 
                    type="number" 
                    placeholder="Max" 
                    className="w-full border border-slate-200 p-2 rounded-md text-sm" 
                    value={range.max ?? ""} 
                    onChange={(e) => onChange({ ...range, max: e.target.value ? parseInt(e.target.value) : undefined })} 
                />
            </div>
        </section>
    );
}