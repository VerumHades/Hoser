import React from "react";
import { useNavigate } from "react-router-dom";
import {
    Rocket,
    ShieldCheck,
    Cpu,
    Globe,
    Zap,
    ArrowRight,
    Code2,
    Terminal
} from "lucide-react";
import MainNavbar from "../components/navigation/MainNavbar";

/**
 * PublicLandingPage serves as the primary entry point for the marketplace.
 * It is structured into distinct vertical sections to guide the user from 
 * value proposition to social proof and finally to action.
 */
export default function PublicLandingPage() {
    return (
        <div className="min-h-screen bg-slate-50 dark:bg-slate-950 text-slate-900 dark:text-slate-100 font-sans selection:bg-indigo-500/30">
            <MainNavbar></MainNavbar>
            <HeroSection />
            <ValuePropositionGrid />
            <FeaturedListingsSection />
            <CallToActionSection />
            <MainFooter />
        </div>
    );
}

/**
 * Renders the primary hero area with the main call to action.
 * Uses a subtle background gradient to draw focus to the headline.
 */
function HeroSection() {
    const navigate = useNavigate();

    return (
        <section className="relative w-full max-w-7xl mx-auto px-6 pt-32 pb-20 flex flex-col items-center text-center gap-8 overflow-hidden">
            <div className="absolute top-0 left-1/2 -translate-x-1/2 w-full h-[500px] bg-indigo-500/10 blur-[120px] rounded-full -z-10" />

            <span className="px-4 py-1.5 rounded-full border border-slate-200 dark:border-slate-800 text-xs font-semibold tracking-widest uppercase bg-white dark:bg-slate-900 shadow-sm">
                Marketplace v1.0 is Live
            </span>

            <h1 className="text-5xl md:text-7xl font-extrabold tracking-tight bg-clip-text text-transparent bg-gradient-to-b from-slate-900 to-slate-600 dark:from-white dark:to-slate-400">
                Infrastructure for those <br /> who ship real software.
            </h1>

            <p className="max-w-2xl text-xl text-slate-600 dark:text-slate-400 leading-relaxed">
                Stop wrestling with YAML files. Browse pre-configured, production-ready
                hosting environments designed by senior engineers. Deploy in seconds.
            </p>

            <div className="flex flex-col sm:flex-row gap-4 mt-4">
                <PrimaryButton onClick={() => navigate("/search")}>
                    Explore Marketplace
                    <ArrowRight className="w-4 h-4" />
                </PrimaryButton>
                <SecondaryButton onClick={() => navigate("/developer")}>
                    Become a Creator
                </SecondaryButton>
            </div>

            <HeroVisualGraphic />
        </section>
    );
}

/**
 * Displays a mock interface or abstract graphic to provide visual weight to the hero.
 */
function HeroVisualGraphic() {
    return (
        <div className="mt-16 w-full max-w-5xl rounded-2xl border border-slate-200 dark:border-slate-800 bg-white/50 dark:bg-slate-900/50 backdrop-blur-sm p-2 shadow-2xl">
            <div className="rounded-xl overflow-hidden aspect-[16/9] relative">
                <img
                    src="carbon.png"
                    alt="Data Center Visualization"
                    className="object-fit opacity-80"
                />
            </div>
        </div>
    );
}

/**
 * A layout for core features using a Bento-grid inspired design.
 */
function ValuePropositionGrid() {
    return (
        <section className="w-full max-w-7xl mx-auto px-6 py-24">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                <FeatureCard
                    icon={<Cpu className="w-6 h-6 text-indigo-500" />}
                    title="Optimized Runtimes"
                    description="Every setup is benchmarked for performance. No bloated containers, just pure efficiency."
                />
                <FeatureCard
                    icon={<ShieldCheck className="w-6 h-6 text-emerald-500" />}
                    title="Security First"
                    description="Immutable infrastructure patterns that follow industry best practices for data protection."
                />
                <FeatureCard
                    icon={<Zap className="w-6 h-6 text-amber-500" />}
                    title="Instant Provisioning"
                    description="Connect your cloud provider and watch your stack come to life in under 60 seconds."
                />
            </div>
        </section>
    );
}

/**
 * Reusable card component for features.
 * Adheres to the principle of small, focused components.
 */
function FeatureCard({ icon, title, description }: { icon: React.ReactNode, title: string, description: string }) {
    return (
        <div className="group p-8 rounded-2xl bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 hover:border-indigo-500/50 transition-all duration-300 shadow-sm">
            <div className="mb-4 p-3 rounded-lg bg-slate-50 dark:bg-slate-800 w-fit group-hover:scale-110 transition-transform">
                {icon}
            </div>
            <h3 className="text-xl font-bold mb-2 text-slate-900 dark:text-slate-100">{title}</h3>
            <p className="text-slate-600 dark:text-slate-400 leading-relaxed">{description}</p>
        </div>
    );
}

function SocialProofSection() {
    return (
        <div className="w-full border-y border-slate-200 dark:border-slate-800 bg-slate-100/50 dark:bg-slate-900/20 py-12">
            <div className="max-w-7xl mx-auto px-6 flex flex-wrap justify-center gap-12 opacity-50 grayscale">
                <BrandPlaceholder name="AWS" />
                <BrandPlaceholder name="DigitalOcean" />
                <BrandPlaceholder name="Terraform" />
                <BrandPlaceholder name="Docker" />
                <BrandPlaceholder name="Kubernetes" />
            </div>
        </div>
    );
}

function BrandPlaceholder({ name }: { name: string }) {
    return <span className="text-xl font-bold tracking-tighter">{name}</span>;
}

function FeaturedListingsSection() {
    return (
        <section className="w-full max-w-7xl mx-auto px-6 py-24">
            <div className="flex flex-col md:flex-row justify-between items-end mb-12 gap-4">
                <div>
                    <h2 className="text-3xl font-bold mb-4">Trending Architectures</h2>
                    <p className="text-slate-600 dark:text-slate-400 max-w-xl">
                        Hand-picked hosting configurations that are currently dominating the developer landscape.
                    </p>
                </div>
                <button className="text-indigo-500 font-semibold hover:text-indigo-400 transition flex items-center gap-2">
                    View all listings <ArrowRight className="w-4 h-4" />
                </button>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
                {/* Sample Static Listings */}
                <ListingCard title="High-Availability K8s" author="infra_guru" price="$49" tags={['Kubernetes', 'AWS']} />
                <ListingCard title="Edge-Ready Next.js" author="vercel_fan" price="$19" tags={['Vercel', 'Edge']} />
                <ListingCard title="Bare Metal Database" author="db_admin" price="$89" tags={['Postgres', 'Linux']} />
            </div>
        </section>
    );
}

function ListingCard({ title, author, price, tags }: any) {
    return (
        <div className="flex flex-col rounded-xl border border-slate-200 dark:border-slate-800 overflow-hidden hover:shadow-xl transition-shadow cursor-pointer">
            <div className="h-48 bg-slate-200 dark:bg-slate-800 flex items-center justify-center">
                <Terminal className="w-12 h-12 opacity-20" />
            </div>
            <div className="p-6">
                <div className="flex justify-between items-start mb-4">
                    <h4 className="font-bold text-lg">{title}</h4>
                    <span className="font-mono text-indigo-500">{price}</span>
                </div>
                <div className="flex gap-2 mb-4">
                    {tags.map((tag: string) => (
                        <span key={tag} className="text-[10px] px-2 py-0.5 rounded-md bg-slate-100 dark:bg-slate-800 border border-slate-200 dark:border-slate-700">
                            {tag}
                        </span>
                    ))}
                </div>
                <div className="text-xs text-slate-500 italic">by @{author}</div>
            </div>
        </div>
    );
}

function CallToActionSection() {
    return (
        <section className="w-full max-w-5xl mx-auto px-6 py-24">
            <div className="rounded-3xl bg-indigo-600 p-8 md:p-16 text-center text-white flex flex-col items-center gap-6 relative overflow-hidden">
                <div className="absolute top-0 right-0 w-64 h-64 bg-white/10 blur-3xl rounded-full translate-x-1/2 -translate-y-1/2" />
                <h2 className="text-3xl md:text-5xl font-bold">Ready to ship your next big idea?</h2>
                <p className="text-indigo-100 max-w-xl text-lg">
                    Join over 2,000 developers who have stopped configuring and started launching.
                </p>
                <button className="px-8 py-4 rounded-full bg-white text-indigo-600 font-bold hover:bg-indigo-50 transition shadow-lg">
                    Get Started Now
                </button>
            </div>
        </section>
    );
}

function MainFooter() {
    return (
        <footer className="w-full border-t border-slate-200 dark:border-slate-800 pt-16 pb-8">
            <div className="max-w-7xl mx-auto px-6 grid grid-cols-2 md:grid-cols-4 gap-12 mb-12">
                <div className="col-span-2">
                    <div className="flex items-center gap-2 font-bold text-xl mb-4">
                        <Rocket className="text-indigo-500" /> ShipMarket
                    </div>
                    <p className="text-slate-500 max-w-xs">
                        The premiere marketplace for professional infrastructure and hosting configurations.
                    </p>
                </div>
                <FooterColumn title="Platform" links={['Marketplace', 'Developers', 'Security', 'Pricing']} />
                <FooterColumn title="Company" links={['About', 'Blog', 'Careers', 'Contact']} />
            </div>
            <div className="text-center text-sm text-slate-500 border-t border-slate-200 dark:border-slate-800 pt-8">
                © {new Date().getFullYear()} — Built for developers who ship
            </div>
        </footer>
    );
}

function FooterColumn({ title, links }: { title: string, links: string[] }) {
    return (
        <div className="flex flex-col gap-4">
            <h5 className="font-bold text-slate-900 dark:text-slate-100">{title}</h5>
            {links.map(link => (
                <a key={link} href="#" className="text-slate-500 hover:text-indigo-500 transition">{link}</a>
            ))}
        </div>
    );
}

function PrimaryButton({ children, onClick }: { children: React.ReactNode, onClick: () => void }) {
    return (
        <button
            onClick={onClick}
            className="px-8 py-4 rounded-full bg-slate-900 text-white dark:bg-white dark:text-slate-900 font-bold hover:scale-105 transition-all flex items-center gap-2 shadow-xl shadow-indigo-500/10"
        >
            {children}
        </button>
    );
}

function SecondaryButton({ children, onClick }: { children: React.ReactNode, onClick: () => void }) {
    return (
        <button
            onClick={onClick}
            className="px-8 py-4 rounded-full border border-slate-300 dark:border-slate-700 text-slate-700 dark:text-slate-300 font-bold hover:bg-white dark:hover:bg-slate-900 transition-all"
        >
            {children}
        </button>
    );
}