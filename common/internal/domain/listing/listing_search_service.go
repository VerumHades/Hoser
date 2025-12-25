package listing

import "common/internal/shared"

type ListingSearchService interface {
	IndexListing(listing *Listing) error
	RemoveListing(listingID shared.ListingID) error
	SearchListings(query string) ([]*Listing, error)
}
