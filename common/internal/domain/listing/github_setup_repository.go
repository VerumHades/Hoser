package listing

import (
	"common/internal/shared"
	"context"
)

// GitHubSetupCommandRepository defines write operations for GitHub setups.
type GitHubSetupCommandRepository interface {
	// Create persists a new GitHub setup within a transaction.
	Create(
		ctx context.Context,
		transaction shared.Transaction,
		setup *GitHubSetupDefinition,
	) (*GitHubSetupDefinition, error)

	// Update modifies an existing GitHub setup within a transaction.
	Update(
		ctx context.Context,
		transaction shared.Transaction,
		setup *GitHubSetupDefinition,
	) (*GitHubSetupDefinition, error)

	// DeleteByID removes a GitHub setup by its ID within a transaction.
	DeleteByID(
		ctx context.Context,
		transaction shared.Transaction,
		setupID shared.SetupID,
	) error

	// DeleteByListingID removes all GitHub setups for a listing within a transaction.
	DeleteByListingID(
		ctx context.Context,
		transaction shared.Transaction,
		listingID shared.ListingID,
	) error
}

// GitHubSetupQueryRepository defines read-only operations for GitHub setups.
type GitHubSetupQueryRepository interface {
	// GetByID retrieves a GitHub setup by its ID.
	GetByID(
		ctx context.Context,
		setupID shared.SetupID,
	) (*GitHubSetupDefinition, error)

	// FetchNextBatchByListing returns GitHub setups for a listing in batches.
	FetchNextBatchByListing(
		ctx context.Context,
		listingID shared.ListingID,
		request shared.BatchRequest,
	) (setups []*GitHubSetupDefinition, nextCursor shared.Cursor, err error)
}
