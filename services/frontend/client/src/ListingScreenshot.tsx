import React, { useEffect, useState } from "react";
import { motion } from "framer-motion";
import { ListingAPI } from "./backend/repositories/listing";

interface ListingScreenshotProps {
    listingId: string;
    screenshotId: string;
    className?: string;
}

export function ListingScreenshot({ listingId, screenshotId, className }: ListingScreenshotProps) {
    const [url, setUrl] = useState<string | null>(null);
    const [hasError, setHasError] = useState(false);

    useEffect(() => {
        let isMounted = true;
        
        async function fetchUrl() {
            try {
                const response = await ListingAPI.getScreenshotReadURL(listingId, screenshotId);
                if (isMounted) setUrl(response.readUrl);
            } catch (err) {
                console.error("Failed to fetch screenshot URL", err);
                if (isMounted) setHasError(true);
            }
        }

        fetchUrl();
        return () => { isMounted = false; };
    }, [listingId, screenshotId]);

    if (hasError) {
        return (
            <div className={`flex items-center justify-center bg-slate-200 dark:bg-slate-800 text-xs text-slate-400 ${className}`}>
                Failed to load
            </div>
        );
    }

    if (!url) {
        return (
            <div className={`animate-pulse bg-slate-200 dark:bg-slate-800 ${className}`} />
        );
    }

    return (
        <motion.img
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            src={url}
            alt="Listing Screenshot"
            className={`object-cover ${className}`}
        />
    );
}