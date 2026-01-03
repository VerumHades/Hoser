package listing

import (
	"common/pkg/shared"
	"context"
	"time"
)

// / ListingSearchCursor represents a stable pagination position
// / for iterating over search results for listings.
type ListingSearchCursor struct {
	LastCreatedAt time.Time
}

// ListingSearchIndex defines operations for indexing and searching listings.
type ListingSearchIndex interface {
	Index(
		ctx context.Context,
		listing *Listing,
	) error

	Update(
		ctx context.Context,
		listing *Listing,
	) error

	Remove(
		ctx context.Context,
		listingID shared.ListingID,
	) error

	SearchNextBatch(
		ctx context.Context,
		query string,
		request shared.BatchRequest[ListingSearchCursor],
	) (listings []*Listing, nextCursor ListingSearchCursor, err error)
}
