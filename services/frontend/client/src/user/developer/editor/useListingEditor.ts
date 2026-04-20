import { useReducer, useState } from "react";
import { DeveloperListingAPI, type DeveloperListing } from "../../../backend/repositories/developer_listing";
import toast from "react-hot-toast";
import type { HardwareSpecification, ListingAccessMode } from "../../../backend/types";


type ListingAction =
    | { type: "set"; listing: DeveloperListing }
    | { type: "setTitle"; title: string }
    | { type: "setDescription"; description: string }
    | { type: "setAccessMode"; mode: ListingAccessMode }
    | { type: "setHardware"; hardware: HardwareSpecification }
    | { type: "setPrice"; currency: number }
    | { type: "reset"; backup: DeveloperListing }
    | { type: "addScreenshot"; id: string }
    | { type: "removeScreenshot"; id: string }
    | { type: "setDocumentationMarkdown"; markdown: string };

function listingReducer(state: DeveloperListing, action: ListingAction): DeveloperListing {
    switch (action.type) {
        case "setTitle": return { ...state, title: action.title };
        case "setDescription": return { ...state, description: action.description };
        case "setAccessMode": return { ...state, accessMode: action.mode };
        case "setHardware": return { ...state, recommended_hardware: { ...state.recommended_hardware, ...action.hardware } };
        case "setPrice": return { ...state, price: action.currency };
        case "reset": return action.backup;
        case "set": return action.listing;
        case "addScreenshot": 
            return { ...state, screenshotIds: [...(state.screenshotIds || []), action.id] };
        case "removeScreenshot":
            return { ...state, screenshotIds: (state.screenshotIds || []).filter(id => id !== action.id) };
        case "setDocumentationMarkdown":
            return {...state, documentation_markdown: action.markdown}
        default: return state;
    }
}

// hooks/useListingEditor.ts
export function useListingEditor(sourceListing: DeveloperListing) {
    const [backupListing, setBackupListing] = useState(sourceListing);
    const [listing, dispatch] = useReducer(listingReducer, sourceListing);
    const [hasUnsavedChanges, setHasUnsavedChanges] = useState(false);

    const markChanged = () => setHasUnsavedChanges(true);
    
    const resetChanges = () => {
        dispatch({ type: "reset", backup: backupListing });
        setHasUnsavedChanges(false);
    };

    const saveChanges = async () => {
        try {
            const a = new Set(backupListing.screenshotIds)
            const b = new Set(listing.screenshotIds)

            console.log(a,b,Array.from(a.difference(b)))
            const toRemove = Array.from(a.difference(b));

            for(const i in toRemove) {
                await DeveloperListingAPI.screenshots.delete(listing.id, toRemove[i])
            }

            const updated = await DeveloperListingAPI.update(listing);
            dispatch({ type: "set", listing: updated });
            setBackupListing(updated);
            setHasUnsavedChanges(false);
            
            toast.success("Changes saved");
        } catch (err) {
            toast.error(err + "");
        }
    };

    return { listing, dispatch, hasUnsavedChanges, markChanged, resetChanges, saveChanges };
}