package repositories

import (
	"common/pkg/domain/entities/listing"
	"common/pkg/shared"
	"context"
	"time"
)

// ListingCommandRepository defines write operations for listings.
type ListingCommandRepository interface {
	// Create persists a new listing within a transaction.
	Create(
		ctx context.Context,
		listing *listing.Listing,
	) error

	// Update modifies an existing listing within a transaction.
	Update(
		ctx context.Context,
		listing *listing.Listing,
	) error

	// Delete removes a listing by its ID within a transaction.
	Delete(
		ctx context.Context,
		listingID shared.ListingID,
	) error
}

// / ListingCursor represents a stable pagination position
// / for iterating over listings.
type ListingCursor struct {
	LastCreatedAt time.Time
}

// ListingQueryRepository defines read-only operations for listings.
type ListingQueryRepository interface {
	GetByID(
		ctx context.Context,
		listingID shared.ListingID,
	) (*listing.Listing, error)

	GetByIDAndAuthor(
		ctx context.Context,
		listingID shared.ListingID,
		userID shared.UserID,
	) (*listing.Listing, error)

	Exists(
		ctx context.Context,
		listingID shared.ListingID,
	) (bool, error)

	FetchNextBatchByAuthor(
		ctx context.Context,
		authorID shared.UserID,
		request shared.BatchRequest[ListingCursor],
	) (listings []*listing.Listing, nextCursor ListingCursor, err error)

	FetchNextBatchAll(
		ctx context.Context,
		request shared.BatchRequest[ListingCursor],
	) (listings []*listing.Listing, nextCursor ListingCursor, err error)
}
