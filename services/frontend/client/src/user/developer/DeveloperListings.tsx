import { DeveloperListingAPI, type DeveloperListing } from "../../backend/repositories/developer_listing";
import { CursorPaginatedCollection } from "../../components/view/CursorPaginatedCollection";
import { useNavigate } from "react-router-dom";
import { DeveloperListingRow } from "./List/DeveloperListingRow";
import { DeveloperListingCard } from "./List/DeveloperListingCard";
import { Button } from "../../templates/components/Button"; // Assuming this is your standard button
import toast from "react-hot-toast";
import { useState } from "react";

export default function DeveloperListings() {
    const navigate = useNavigate();
    const [isCreating, setIsCreating] = useState(false);

    const handleCreateNew = async () => {
        setIsCreating(true);
        const toastId = toast.loading("Creating your listing...");
        try {
            // Default values handled by the Go service we wrote earlier
            const newListing = await DeveloperListingAPI.create("New Listing", "A short description of your software");
            
            toast.success("Listing created!", { id: toastId });
            
            // Navigate to the editor route
            navigate(`/dashboard/developer/listing/${newListing.id}`);
        } catch (err) {
            toast.error("Failed to create listing: " + err, { id: toastId });
        } finally {
            setIsCreating(false);
        }
    };

    return (
        <div className="flex-1 min-h-0 flex flex-col items-center">
            {/* Header Section with Create Button */}
            <div className="w-full max-w-4xl flex justify-between items-end mt-8 px-6">
                <div>
                    <h1 className="text-2xl font-bold text-slate-900 dark:text-white">Your Listings</h1>
                    <p className="text-slate-500 dark:text-slate-400 text-sm">Manage and monitor your published software.</p>
                </div>
                
                <Button 
                    variant="primary" 
                    onClick={handleCreateNew} 
                    disabled={isCreating}
                >
                    {isCreating ? "Creating..." : "Create Listing"}
                </Button>
            </div>

            <div className="w-full max-w-4xl min-h-0 overflow-hidden mt-5 px-6">
                <CursorPaginatedCollection<DeveloperListing>
                    fetchPage={cursor => DeveloperListingAPI.list(cursor)}
                    renderItemRow={listing => (
                        <DeveloperListingRow
                            key={listing.id}
                            listing={listing}
                            onSelect={() => navigate(`/dashboard/developer/listing/${listing.id}`)}
                        />
                    )}
                    renderItemCard={listing => (
                        <DeveloperListingCard
                            key={listing.id}
                            listing={listing}
                            onSelect={() => navigate(`/dashboard/developer/listing/${listing.id}`)}
                        />
                    )}
                    emptyState={
                        <div className="flex flex-col items-center justify-center py-12 text-center border-2 border-dashed rounded-xl border-slate-200 dark:border-slate-800">
                            <p className="text-sm text-slate-500 mb-4">No listings found</p>
                            <Button variant="secondary" onClick={handleCreateNew}>
                                Create your first listing
                            </Button>
                        </div>
                    }
                />
            </div>
        </div>
    );
}