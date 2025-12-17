import { useEffect, useState, type JSX } from "react";
import { useSearchParams, type NavigateFunction } from "react-router-dom";
import { API, type Listing } from "../backend";
import { TitleAndDescription } from "../components/prefabs/TitleAndDescription";
import { BookmarkPlus } from "lucide-react";
import { PriceTag } from "../user/developer/prices/PriceTag";

export function gotoListing(navigate: NavigateFunction, listing: Listing | undefined) {
    if(!listing) return;

    const url = new URLSearchParams()
    url.append("id", listing.id)
    navigate("/listing?" + url.toString())
}

export function PublicListingView(): JSX.Element {
    const [searchParameters] = useSearchParams();
    const listingID = searchParameters.get("id");
    
    const [listing, setListing] = useState<Listing | null>(null);
    const [isLoading, setIsLoading] = useState<boolean>(true);
    const [hasError, setHasError] = useState<boolean>(false);

    useEffect(() => {

        async function load(){
            const result = await API.listing.get(listingID ?? "")
            console.log(result)
            if(result.ok){
                setListing(result.json)
            }
            else {
                setHasError(true)
            }
            setIsLoading(false)
        }

        load()
    },[]);

    if (isLoading) {
        return <div className="p-6">Loading…</div>;
    }

    if (hasError || !listing) {
        return <div className="p-6">Listing not found</div>;
    }

    return (
        <div className="w-full h-full flex flex-col p-6 min-h-0">
            <div className="flex flex-row justify-between items-end">
                <TitleAndDescription
                    title={listing.title ?? "[No title]"}
                    description={listing.description ?? "[No description]"}
                    titleClassname="text-7xl"
                />
                <div className="text-slate-500 dark:text-slate-400 text-sm flex flex-col items-end gap-2">
                    {listing.price && <PriceTag prices={[listing.price]} rate={"one time"} />}
                    <BookmarkPlus onClick={() => API.user.library.add(listing.id)} />
                </div>
            </div>

            <div className="flex flex-col flex-1 min-h-0 overflow-auto">
                {/* additional listing content */}
            </div>
        </div>
    );
}
