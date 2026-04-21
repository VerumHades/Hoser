import { DeveloperListingAPI, type DeveloperListing } from "../../backend/repositories/developer_listing";
import { useNavigate } from "react-router-dom";
import { DeveloperListingRow } from "./List/DeveloperListingRow";
import { DeveloperListingCard } from "./List/DeveloperListingCard";
import { Button } from "../../templates/components/Button"; // Assuming this is your standard button
import toast from "react-hot-toast";
import { useState } from "react";
import { SearchableCollection } from "../../components/view/SearchableCollection";
import { type FilterField } from "../../components/querying/FilterSidebar";
import type { ListingSearchQuery } from "../../backend/utils/query";

/**
 * A comprehensive collection of filter fields available for developer listings.
 * Includes technical specifications, metadata, and internal identifiers.
 */
const DEVELOPER_COMPLETE_FILTER_FIELDS: FilterField<ListingSearchQuery>[] = [
    { 
        key: "text", 
        label: "Text", 
        type: "text", 
        urlKey: "title" 
    },
    { 
        key: "authorId", 
        label: "Author Reference", 
        type: "text", 
        urlKey: "author" 
    },
    { 
        key: "accessMode", 
        label: "System Access Mode", 
        type: "select", 
        urlKey: "access_mode" 
    },
    { 
        key: "price", 
        label: "Listing Price", 
        type: "range", 
        urlKeys: { 
            min: "price_min", 
            max: "price_max" 
        } 
    },
    { 
        key: "cpu", 
        label: "Recommended CPU Cores", 
        type: "range", 
        urlKeys: { 
            min: "cpu_min", 
            max: "cpu_max" 
        } 
    },
    { 
        key: "ramBytes", 
        label: "Recommended RAM Bytes", 
        type: "range", 
        urlKeys: { 
            min: "ram_min", 
            max: "ram_max" 
        } 
    },
    { 
        key: "diskBytes", 
        label: "Recommended Disk Bytes", 
        type: "range", 
        urlKeys: { 
            min: "disk_min", 
            max: "disk_max" 
        } 
    }
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
            navigate(`/dashboard/developer/listing/${newListing?.id}`);
        } catch (err) {
            toast.error(`Failed to create listing: ${err}`, { id: toastId });
        } finally {
            setIsCreating(false);
        }
    };

    return (
        <SearchableCollection<DeveloperListing, ListingSearchQuery>
            filterFields={DEVELOPER_COMPLETE_FILTER_FIELDS}
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