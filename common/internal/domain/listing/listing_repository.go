package listing

import (
	"common/internal/shared"
	"context"
)

// ListingCommandRepository defines write operations for listings.
type ListingCommandRepository interface {
	// Create persists a new listing within a transaction.
	Create(
		ctx context.Context,
		transaction shared.Transaction,
		listing *Listing,
	) (*Listing, error)

	// Update modifies an existing listing within a transaction.
	Update(
		ctx context.Context,
		transaction shared.Transaction,
		listing *Listing,
	) (*Listing, error)

	// Delete removes a listing by its ID within a transaction.
	Delete(
		ctx context.Context,
		transaction shared.Transaction,
		listingID shared.ListingID,
	) error
}

// ListingQueryRepository defines read-only operations for listings.
type ListingQueryRepository interface {
	// GetByID retrieves a listing by its ID.
	GetByID(
		ctx context.Context,
		listingID shared.ListingID,
	) (*Listing, error)

	// FetchNextBatchByAuthor returns listings for a given author in batches.
	FetchNextBatchByAuthor(
		ctx context.Context,
		authorID shared.UserID,
		request shared.BatchRequest,
	) (listings []*Listing, nextCursor shared.Cursor, err error)

	// FetchNextBatchAll returns all listings in batches.
	FetchNextBatchAll(
		ctx context.Context,
		request shared.BatchRequest,
	) (listings []*Listing, nextCursor shared.Cursor, err error)
}
