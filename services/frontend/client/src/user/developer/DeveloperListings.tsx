import { DeveloperListingAPI, type DeveloperListing } from "../../backend/repositories/developer_listing";
import { CursorPaginatedCollection } from "../../components/view/CursorPaginatedCollection";
import { useNavigate } from "react-router-dom";
import { DeveloperListingRow } from "./List/DeveloperListingRow";
import { DeveloperListingCard } from "./List/DeveloperListingCard";
import { Button } from "../../templates/components/Button"; // Assuming this is your standard button
import toast from "react-hot-toast";
import { useState } from "react";
import { SearchableCollection } from "../../components/view/SearchableCollection";

export default function DeveloperListings() {
    const navigate = useNavigate();
    const [isCreating, setIsCreating] = useState(false);

    const handleCreate = async () => {
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
		<SearchableCollection<DeveloperListing, { q?: string }>
			title="Your Listings"
			description="Manage and monitor your published software."
            
			headerActions={
				<Button onClick={handleCreate} disabled={isCreating}>
					{isCreating ? "Creating..." : "Create Listing"}
				</Button>
			}
			initialQuery={{}}

			fetchPage={(q, cursor) => DeveloperListingAPI.search( {text: q as string}, cursor)}
			parseParams={(p) => ({ q: p.get("q") || undefined })}
			buildParams={(q) => new URLSearchParams(q.q ? { q: q.q } : {})}

			renderRow={(l) => <DeveloperListingRow listing={l} onSelect={() => navigate(`/editor/${l.id}`)} />}
			renderCard={(l) => <DeveloperListingCard listing={l} onSelect={() => navigate(`/editor/${l.id}`)} />}

			emptyState={
				<div className="text-center py-12 border-2 border-dashed rounded-xl">
					<p className="text-sm text-slate-500 mb-4">No listings found</p>
					<Button variant="secondary" onClick={handleCreate}>Create your first</Button>
				</div>
			}
		/>
	);
}