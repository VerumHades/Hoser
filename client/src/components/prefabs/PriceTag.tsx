interface Price {
    name: string,
    value: number,
    currency_short: string
}

interface PriceTagProps {
    prices: Price[],
    rate?: string,
    onClick?: () => void
}

export function PriceTag({ prices, rate, onClick }: PriceTagProps) {
    const totalValue = prices.reduce((sum, p) => sum + p.value, 0);
    const currency = prices[0]?.currency_short || "";

    return (
        <div
            onClick={onClick}
            className="inline-flex items-center gap-2 px-4 py-2 rounded-lg
                bg-slate-100 dark:bg-slate-800
                border border-slate-300 dark:border-slate-700
                text-slate-800 dark:text-slate-100
                shadow-sm cursor-pointer relative overflow-hidden transition-all duration-300"
        >
            {/* Total is always visible */}
            <span className="font-medium">{totalValue}{currency}</span>
            {rate && (
                <span className="text-slate-500 dark:text-slate-400 text-sm whitespace-nowrap">
                    / {rate}
                </span>
            )}
        </div>
    );
}