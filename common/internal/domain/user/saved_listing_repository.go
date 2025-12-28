package user

import (
	"common/internal/shared"
	"context"
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

// SavedListingQueryRepository defines read-only operations for user saved listings.
type SavedListingQueryRepository interface {
	// GetByID retrieves a saved listing by its ID.
	GetByID(
		ctx context.Context,
		itemID shared.SavedListingID,
	) (*SavedListing, error)

	// ExistsByUserAndListing checks if a user has saved a specific listing.
	ExistsByUserAndListing(
		ctx context.Context,
		userID shared.UserID,
		listingID shared.ListingID,
	) (bool, error)

	// FetchNextBatchByUser returns a batch of saved listings for a user.
	FetchNextBatchByUser(
		ctx context.Context,
		userID shared.UserID,
		request shared.BatchRequest,
	) (savedListings []*SavedListing, nextCursor shared.Cursor, err error)
}
