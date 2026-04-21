
import { useNavigate } from "react-router-dom";
import { ListingAPI, type Listing } from "../backend/repositories/listing";
import { PublicListingRow } from "./List/PublicListingRow";
import { PublicListingCard } from "./List/PublicListingCard";
import { SearchableCollection } from "../components/view/SearchableCollection";
import type { ListingSearchQuery } from "../backend/utils/query";
import { type FilterField } from "../components/querying/FilterSidebar";
import MainNavbar from "../components/navigation/MainNavbar";

const PUBLIC_FILTER_FIELDS: FilterField<ListingSearchQuery>[] = [
    { 
        key: "text", 
        label: "Search", 
        type: "text", 
        urlKey: "q" 
    },
    { 
        key: "cpu", 
        label: "CPU Cores", 
        type: "range", 
        urlKeys: { min: "cpu_min", max: "cpu_max" } 
    },
    { 
        key: "ramBytes", 
        label: "RAM (Bytes)", 
        type: "range", 
        urlKeys: { min: "ram_min", max: "ram_max" } 
    },
    { 
        key: "diskBytes", 
        label: "Disk (Bytes)", 
        type: "range", 
        urlKeys: { min: "disk_min", max: "disk_max" } 
    },
    { 
        key: "price", 
        label: "Price", 
        type: "range", 
        urlKeys: { min: "price_min", max: "price_max" } 
    },
    { 
        key: "createdAt", 
        label: "Creation Date", 
        type: "date", 
        urlKeys: { from: "date_from", to: "date_to" } 
    },
    { 
        key: "authorId", 
        label: "Author ID", 
        type: "text", 
        urlKey: "author_id" 
    },
];

export default function PublicListingSearch() {
	const navigate = useNavigate();

	return (<div className="h-full flex flex-col">
		<MainNavbar></MainNavbar>
		<SearchableCollection<Listing, ListingSearchQuery>
			initialQuery={{}}
			filterFields={PUBLIC_FILTER_FIELDS}
			fetchPage={(q, cursor) => ListingAPI.search(q, cursor)}
			renderRow={(l) => <PublicListingRow listing={l} onSelect={() => navigate(`/listing/${l.id}`)} />}
			renderCard={(l) => <PublicListingCard listing={l} onSelect={() => navigate(`/listing/${l.id}`)} />}
		/>
	</div>
		
	);
}