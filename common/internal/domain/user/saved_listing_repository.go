package user

import "common/internal/shared"

// SavedListingRepository defines persistence operations for user saved listings.
type SavedListingRepository interface {
	Save(item *SavedListing) (*SavedListing, error)

	GetByID(itemID shared.SavedListingID) (*SavedListing, error)

	Delete(itemID shared.SavedListingID) error

	ExistsByUserAndListing(
		userID shared.UserID,
		listingID shared.ListingID,
	) (bool, error)

	FetchNextBatchByUser(
		userID shared.UserID,
		lastSeenSavedListingID shared.SavedListingID,
		maximumBatchSize int,
	) ([]*SavedListing, error)
}
