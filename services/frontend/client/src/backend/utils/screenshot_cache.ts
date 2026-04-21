import { ListingAPI } from "../repositories/listing";

// Store promises so that simultaneous requests for the same ID only trigger one API call
const urlCache = new Map<string, Promise<string>>();

export const ScreenshotCache = {
    async getUrl(listingId: string, screenshotId: string): Promise<string> {
        const cacheKey = `${listingId}:${screenshotId}`;

        if (urlCache.has(cacheKey)) {
            return urlCache.get(cacheKey)!;
        }

        const fetchPromise = (async () => {
            try {
                const response = await ListingAPI.getScreenshotReadURL(listingId, screenshotId);
                return response?.readUrl || "";
            } catch (err) {
                urlCache.delete(cacheKey); // Don't cache failures
                throw err;
            }
        })();

        urlCache.set(cacheKey, fetchPromise);
        return fetchPromise;
    }
};