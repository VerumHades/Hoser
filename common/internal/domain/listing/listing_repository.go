package listing

import "common/internal/shared"

type ListingRepository interface {
	Save(listing *Listing) (*Listing, error)

	GetByID(listingID shared.ListingID) (*Listing, error)

	Delete(listingID shared.ListingID) error

	FetchNextBatchByAuthor(
		authorID shared.UserID,
		lastSeenListingID shared.ListingID,
		maximumBatchSize int,
	) ([]*Listing, error)

	FetchNextBatchAll(
		lastSeenListingID shared.ListingID,
		maximumBatchSize int,
	) ([]*Listing, error)
}
