package repositories

import (
	"common/pkg/domain/entities/user"
	"common/pkg/shared"
	"context"
	"time"
)

// SavedListingCommandRepository defines write operations for user saved listings.
type SavedListingCommandRepository interface {
	// Create persists a new saved listing within a transaction.
	Create(
		ctx context.Context,

		item *user.SavedListing,
	) (*user.SavedListing, error)

	// Delete removes a saved listing by its ID within a transaction.
	Delete(
		ctx context.Context,

		itemID shared.SavedListingID,
	) error
}

// / SavedListingCursor represents a stable pagination position
// / for iterating over user saved listings.
type SavedListingCursor struct {
	LastCreatedAt time.Time
}

// SavedListingQueryRepository defines read-only operations for user saved listings.
type SavedListingQueryRepository interface {
	GetByID(
		ctx context.Context,
		itemID shared.SavedListingID,
	) (*user.SavedListing, error)

	GetByUserAndListing(
		ctx context.Context,
		userID shared.UserID,
		listingID shared.ListingID,
	) (*user.SavedListing, error)

	ExistsByUserAndListing(
		ctx context.Context,
		userID shared.UserID,
		listingID shared.ListingID,
	) (bool, error)

	FetchNextBatchByUser(
		ctx context.Context,
		userID shared.UserID,
		request shared.BatchRequest[SavedListingCursor],
	) (savedListings []*user.SavedListing, nextCursor SavedListingCursor, err error)
}
