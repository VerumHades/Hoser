package shared

import (
	"fmt"
)

func GenerateFileKey(listingID ListingID, screenshotID ListingScreenshotID) string {
	return fmt.Sprintf("listings/%s/screenshots/%s",
		listingID,
		screenshotID,
	)
}
