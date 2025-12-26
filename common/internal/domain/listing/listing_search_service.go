package listing

import "common/internal/shared"

type ListingSearchService interface {
	IndexListing(listing *Listing) error

	RemoveListing(listingID shared.ListingID) error

	SearchNextBatch(
		query string,
		lastSeenListingID shared.ListingID,
		maximumBatchSize int,
	) ([]*Listing, error)
}
