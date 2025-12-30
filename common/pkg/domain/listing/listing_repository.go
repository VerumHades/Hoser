package listing

import (
	"common/pkg/shared"
	"context"
	"time"
)

// ListingCommandRepository defines write operations for listings.
type ListingCommandRepository interface {
	// Create persists a new listing within a transaction.
	Create(
		ctx context.Context,
		transaction shared.Transaction,
		listing *Listing,
	) error

	// Update modifies an existing listing within a transaction.
	Update(
		ctx context.Context,
		transaction shared.Transaction,
		listing *Listing,
	) error

	// Delete removes a listing by its ID within a transaction.
	Delete(
		ctx context.Context,
		transaction shared.Transaction,
		listingID shared.ListingID,
	) error
}

// / ListingCursor represents a stable pagination position
// / for iterating over listings.
type ListingCursor struct {
	LastCreatedAt time.Time
	LastListingID shared.ListingID
}

// ListingQueryRepository defines read-only operations for listings.
type ListingQueryRepository interface {
	GetByID(
		ctx context.Context,
		listingID shared.ListingID,
	) (*Listing, error)

	Exists(
		ctx context.Context,
		listingID shared.ListingID,
	) (bool, error)

	FetchNextBatchByAuthor(
		ctx context.Context,
		authorID shared.UserID,
		request shared.BatchRequest[ListingCursor],
	) (listings []*Listing, nextCursor ListingCursor, err error)

	FetchNextBatchAll(
		ctx context.Context,
		request shared.BatchRequest[ListingCursor],
	) (listings []*Listing, nextCursor ListingCursor, err error)
}
