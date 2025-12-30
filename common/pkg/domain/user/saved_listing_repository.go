package user

import (
	"common/pkg/shared"
	"context"
	"time"
)

// SavedListingCommandRepository defines write operations for user saved listings.
type SavedListingCommandRepository interface {
	// Create persists a new saved listing within a transaction.
	Create(
		ctx context.Context,
		transaction shared.Transaction,
		item *SavedListing,
	) (*SavedListing, error)

	// Update modifies an existing saved listing within a transaction.
	Update(
		ctx context.Context,
		transaction shared.Transaction,
		item *SavedListing,
	) (*SavedListing, error)

	// Delete removes a saved listing by its ID within a transaction.
	Delete(
		ctx context.Context,
		transaction shared.Transaction,
		itemID shared.SavedListingID,
	) error
}

// / SavedListingCursor represents a stable pagination position
// / for iterating over user saved listings.
type SavedListingCursor struct {
	LastCreatedAt   time.Time
	LastSavedItemID shared.SavedListingID
}

// SavedListingQueryRepository defines read-only operations for user saved listings.
type SavedListingQueryRepository interface {
	GetByID(
		ctx context.Context,
		itemID shared.SavedListingID,
	) (*SavedListing, error)

	ExistsByUserAndListing(
		ctx context.Context,
		userID shared.UserID,
		listingID shared.ListingID,
	) (bool, error)

	FetchNextBatchByUser(
		ctx context.Context,
		userID shared.UserID,
		request shared.BatchRequest[SavedListingCursor],
	) (savedListings []*SavedListing, nextCursor SavedListingCursor, err error)
}
