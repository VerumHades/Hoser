// src/components/listings/DeveloperListingLoader.tsx
import { useEffect, useState } from "react"
import DeveloperListingEditor from "./DeveloperListingEditor"
import { useParams } from "react-router-dom"
import { DeveloperListingAPI, type DeveloperListing } from "../../backend/repositories/developer_listing"

interface DeveloperListingLoaderProps {
    onClose?: () => void
}

export default function DeveloperListingLoader({
    onClose
}: DeveloperListingLoaderProps) {
    const {id: listingId} = useParams<{id: string}>()

    
    if (!listingId) {
        return (
            <div className="flex items-center justify-center h-full w-full">
                <span className="text-gray-500">Invalid listing id...</span>
            </div>
        )
    }

    const [listing, setListing] = useState<DeveloperListing | null>(null)
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState<string | null>(null)

    useEffect(() => {
        let isMounted = true

        const fetchListing = async () => {
            
            setLoading(true)
            setError(null)
            try {
                const response = await DeveloperListingAPI.get(listingId)
                if (isMounted) setListing(response)
            } catch (err) {
                if (isMounted) setError((err as Error).message || "Unknown error")
            } finally {
                if (isMounted) setLoading(false)
            }
        }

        fetchListing()

        return () => {
            isMounted = false
        }
    }, [listingId])

    if (loading) {
        return (
            <div className="flex items-center justify-center h-full w-full">
                <span className="text-gray-500">Loading listing...</span>
            </div>
        )
    }

    if (error) {
        return (
            <div className="flex flex-col items-center justify-center h-full w-full text-red-500">
                <span>Error: {error}</span>
                {onClose && (
                    <button
                        className="mt-4 px-3 py-1 bg-gray-200 rounded hover:bg-gray-300"
                        onClick={onClose}
                    >
                        Close
                    </button>
                )}
            </div>
        )
    }

    if (!listing) {
        return (
            <div className="flex items-center justify-center h-full w-full text-gray-500">
                No listing found.
            </div>
        )
    }

    return <DeveloperListingEditor listing={listing} onShouldClose={onClose} />
}
