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

// NumericRange defines a minimum and maximum bound for search attributes.
type NumericRange struct {
	Min *int64
	Max *int64
}

// DateRange defines a time-based window for search queries.
type DateRange struct {
	From *time.Time
	To   *time.Time
}

// SearchQuery encapsulates complex filtering logic for listings.
type SearchQuery struct {
	Text       string
	CPU        NumericRange
	RAMBytes   NumericRange
	DiskBytes  NumericRange
	Price      NumericRange
	CreatedAt  DateRange
	AccessMode *ListingAccessMode
	AuthorID   *shared.UserID
}

// ListingSearchIndex defines operations for indexing and searching listings.
type ListingSearchIndex interface {
	Index(ctx context.Context, listing *Listing) error
	Update(ctx context.Context, listing *ListingUpdateEvent) error
	Remove(ctx context.Context, listingID shared.ListingID) error

	/**
	 * SearchNextBatch executes a granular search with support for
	 * range filters and text matching.
	 */
	SearchNextBatch(
		ctx context.Context,
		query SearchQuery,
		request shared.BatchRequest[ListingSearchCursor],
	) ([]*Listing, ListingSearchCursor, error)
}
