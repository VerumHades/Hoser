package listing

import (
	"common/internal/shared"
	"context"
)

// ListingSearchIndex defines operations for indexing and searching listings.
type ListingSearchIndex interface {
	// Index adds a new listing to the search index.
	Index(
		ctx context.Context,
		listing *Listing,
	) error

	// Update modifies an existing listing in the search index.
	Update(
		ctx context.Context,
		listing *Listing,
	) error

	// Remove deletes a listing from the search index by its ID.
	Remove(
		ctx context.Context,
		listingID shared.ListingID,
	) error

	// SearchNextBatch retrieves a batch of listings matching a query,
	// starting after the last seen listing ID, using the provided batch request.
	SearchNextBatch(
		ctx context.Context,
		query string,
		request shared.BatchRequest,
	) (listings []*Listing, nextCursor shared.Cursor, err error)
}
