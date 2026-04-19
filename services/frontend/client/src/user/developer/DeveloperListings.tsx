import { DeveloperListingAPI, type DeveloperListing } from "../../backend/repositories/developer_listing";
import { CursorPaginatedCollection } from "../../components/view/CursorPaginatedCollection";
import { useNavigate } from "react-router-dom";
import { DeveloperListingRow } from "./List/DeveloperListingRow";
import { DeveloperListingCard } from "./List/DeveloperListingCard";
import { Button } from "../../templates/components/Button"; // Assuming this is your standard button
import toast from "react-hot-toast";
import { useState } from "react";
import { SearchableCollection } from "../../components/view/SearchableCollection";
import { DynamicFilterSidebar, type FilterField } from "../../components/querying/FilterSidebar";
import type { ListingSearchQuery } from "../../backend/utils/query";

const SIMPLE_SEARCH_FIELDS: FilterField<ListingSearchQuery>[] = [
    { 
        key: "text", 
        label: "Keywords", 
        type: "text", 
        urlKey: "q" 
    },
];

export default function DeveloperListings() {
    const navigate = useNavigate();
    const [isCreating, setIsCreating] = useState(false);

    const handleCreate = async () => {
        setIsCreating(true);
        const toastId = toast.loading("Creating your listing...");
        try {
            const newListing = await DeveloperListingAPI.create(
                "New Listing", 
                "A short description of your software"
            );
            toast.success("Listing created!", { id: toastId });
            navigate(`/dashboard/developer/listing/${newListing.id}`);
        } catch (err) {
            toast.error(`Failed to create listing: ${err}`, { id: toastId });
        } finally {
            setIsCreating(false);
        }
    };

    return (
        <SearchableCollection<DeveloperListing, ListingSearchQuery>
            filterFields={SIMPLE_SEARCH_FIELDS}
            headerActions={
                <Button onClick={handleCreate} disabled={isCreating}>
                    {isCreating ? "Creating..." : "Create Listing"}
                </Button>
            }
            fetchPage={(query, cursor) => DeveloperListingAPI.search(query, cursor)}
            renderRow={(listing) => (
                <DeveloperListingRow 
                    listing={listing} 
                    onSelect={() => navigate(`/dashboard/developer/listing/${listing.id}`)} 
                />
            )}
            renderCard={(listing) => (
                <DeveloperListingCard 
                    listing={listing} 
                    onSelect={() => navigate(`/dashboard/developer/listing/${listing.id}`)} 
                />
            )}
            emptyState={
                <DeveloperEmptyState onCreate={handleCreate} />
            }
        />
    );
}

/**
 * Clean UI fragment for the empty state
 */
function DeveloperEmptyState({ onCreate }: { onCreate: () => void }) {
    return (
        <div className="flex flex-col items-center justify-center py-12 text-center border-2 border-dashed rounded-xl border-slate-200">
            <p className="text-sm text-slate-500 mb-4">No listings found</p>
            <Button variant="secondary" onClick={onCreate}>
                Create your first listing
            </Button>
        </div>
    );
}