export function ListingSpecs() {
    return (
        <div className="p-6 rounded-3xl bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800">
            <h3 className="font-bold text-lg mb-4">Technical Specs</h3>
            <ul className="space-y-3 text-sm text-slate-600 dark:text-slate-400">
                <li className="flex justify-between"><span>Version</span> <span className="text-slate-900 dark:text-white font-mono text-xs">v1.4.2</span></li>
                <li className="flex justify-between"><span>Last Updated</span> <span className="text-slate-900 dark:text-white font-medium">Oct 24, 2025</span></li>
                <li className="flex justify-between"><span>License</span> <span className="text-slate-900 dark:text-white font-medium">Commercial</span></li>
            </ul>
            <div className="mt-6 pt-6 border-t border-slate-200 dark:border-slate-800">
                <p className="text-xs leading-relaxed italic">Includes full deployment documentation and 24/7 technical support access for owners.</p>
            </div>
        </div>
    );
}