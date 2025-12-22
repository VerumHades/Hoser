import { useEffect, useState, type JSX } from "react";
import { useParams, useSearchParams, type NavigateFunction } from "react-router-dom";
import { API, type Listing } from "../backend";
import { TitleAndDescription } from "../components/prefabs/TitleAndDescription";
import { BookmarkPlus } from "lucide-react";
import { PriceTag } from "../user/developer/prices/PriceTag";

export function gotoListing(
    navigate: NavigateFunction,
    listing: Listing | undefined
) {
    if (!listing) {
        return;
    }

    navigate("/listing/" + listing.id);
}

export function PublicListingView(): JSX.Element {
    const {id} = useParams<{id: string}>();
    const listingID = id;

    const [listing, setListing] = useState<Listing | null>(null);
    const [isLoading, setIsLoading] = useState<boolean>(true);
    const [hasError, setHasError] = useState<boolean>(false);

    useEffect(() => {
        async function loadListing() {
            const result = await API.listing.get(listingID ?? "");

            if (result.ok) {
                setListing(result.json);
            } else {
                setHasError(true);
            }

            setIsLoading(false);
        }

        loadListing();
    }, [listingID]);

    if (isLoading) {
        return (
            <div className="flex items-center justify-center h-full text-slate-500">
                Loading listing…
            </div>
        );
    }

    if (hasError || !listing) {
        return (
            <div className="flex items-center justify-center h-full text-slate-500">
                Listing not found
            </div>
        );
    }

    return (
        <div className="w-full h-full flex flex-col gap-8 px-6 py-8 max-w-7xl mx-auto">
            <section
                className="
                    flex flex-col gap-6
                    border-b
                    border-slate-200
                    dark:border-slate-800
                    pb-6
                "
            >
                <div className="flex flex-col md:flex-row justify-between gap-6">
                    <TitleAndDescription
                        title={listing.title ?? "[No title]"}
                        description={listing.description ?? "[No description]"}
                        titleClassname="text-4xl md:text-5xl"
                        descriptionClassname="max-w-3xl"
                    />

                    <div
                        className="
                            flex flex-col items-start md:items-end
                            gap-3
                            text-slate-600
                            dark:text-slate-400
                        "
                    >
                        {listing.price && (
                            <PriceTag
                                prices={[listing.price]}
                                rate="one time"
                            />
                        )}

                        <button
                            onClick={() => API.user.library.add(listing.id)}
                            className="
                                flex items-center gap-2
                                px-3 py-2
                                rounded-lg
                                border
                                border-slate-300
                                dark:border-slate-700
                                hover:bg-slate-100
                                dark:hover:bg-slate-800
                                transition
                            "
                        >
                            <BookmarkPlus className="w-4 h-4" />
                            <span className="text-sm">Save</span>
                        </button>
                    </div>
                </div>
            </section>

            <section className="flex flex-col flex-1 min-h-0 overflow-auto gap-6">
                <div
                    className="
                        rounded-xl
                        border
                        border-slate-200
                        dark:border-slate-800
                        bg-slate-50
                        dark:bg-slate-900
                        p-6
                        text-slate-700
                        dark:text-slate-300
                    "
                >
                    <p className="text-sm leading-relaxed">
                        Additional listing content goes here.
                        Deployment instructions, screenshots,
                        supported providers, versioning, etc.
                    </p>
                </div>
            </section>
        </div>
    );
}
