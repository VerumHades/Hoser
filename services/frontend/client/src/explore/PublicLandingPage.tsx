import backend_constants from "../backend_constants";
import { useNavigate } from "react-router-dom";
import Table from "../components/table/Table";
import TableRow from "../components/table/TableRow";
import Query from "../components/querying/Query";
import { type Listing } from "../backend";
import { PriceTag } from "../user/developer/prices/PriceTag";
import { gotoListing } from "./PublicListingView";

export default function PublicLandingPage() {
    const navigate = useNavigate();

    const queryBuilder = (q: string) => {
        return { q };
    };

    return (
        <div className="w-full h-full flex flex-col items-center overflow-y-auto">
            <section className="w-full max-w-7xl px-6 py-24 flex flex-col gap-8">
                <h1 className="text-4xl md:text-6xl font-bold tracking-tight text-slate-900 dark:text-slate-100">
                    Deploy real software.
                    <br />
                    One click.
                </h1>

                <p className="max-w-2xl text-lg text-slate-700 dark:text-slate-400">
                    A developer marketplace for production-ready hosting setups.
                    Buy, deploy, and run infrastructure created by developers who
                    actually ship.
                </p>

                <div className="flex gap-4 mt-6">
                    <button
                        onClick={() => navigate("/search")}
                        className="px-6 py-3 rounded-lg bg-slate-900 text-white dark:bg-slate-100 dark:text-slate-900 font-medium hover:opacity-90 transition"
                    >
                        Browse setups
                    </button>

                    <button
                        onClick={() => navigate("/developer")}
                        className="px-6 py-3 rounded-lg border border-slate-300 dark:border-slate-700 text-slate-700 dark:text-slate-300 font-medium hover:bg-slate-100 dark:hover:bg-slate-800 transition"
                    >
                        Sell your setup
                    </button>
                </div>
            </section>

            <section className="w-full max-w-7xl px-6 pb-24 flex flex-col gap-12">
                <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
                    <Feature
                        title="Built by developers"
                        description="Every listing is a real, deployable setup created by someone who understands production."
                    />
                    <Feature
                        title="One-click hosting"
                        description="No tutorials. No guessing. Buy a setup and deploy instantly."
                    />
                    <Feature
                        title="Open marketplace"
                        description="Set your own pricing, ship updates, and get paid for real infrastructure work."
                    />
                </div>

                <div className="flex flex-col gap-4">
                    <h2 className="text-2xl font-semibold text-slate-900 dark:text-slate-200">
                        Featured setups
                    </h2>
                </div>
            </section>

            <footer className="w-full border-t border-slate-200 dark:border-slate-800 py-8 flex justify-center text-sm text-slate-500">
                © {new Date().getFullYear()} — Built for developers who ship
            </footer>
        </div>
    );
}

function Feature(props: { title: string; description: string }) {
    return (
        <div className="flex flex-col gap-2 p-6 rounded-xl bg-slate-100 dark:bg-slate-900 border border-slate-200 dark:border-slate-800">
            <h3 className="font-semibold text-slate-900 dark:text-slate-200">
                {props.title}
            </h3>
            <p className="text-sm text-slate-600 dark:text-slate-400">
                {props.description}
            </p>
        </div>
    );
}
