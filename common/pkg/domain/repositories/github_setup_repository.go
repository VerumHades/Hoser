package repositories

import (
	"common/pkg/domain/entities/githubsetups"
	"common/pkg/shared"
	"context"
	"time"
)

// GitHubSetupCommandRepository defines write operations for GitHub setups.
type GitHubSetupCommandRepository interface {
	// Create persists a new GitHub setup within a transaction.
	Create(
		ctx context.Context,
		setup *githubsetups.GitHubSetupDefinition,
	) (*githubsetups.GitHubSetupDefinition, error)

	// Update modifies an existing GitHub setup within a transaction.
	Update(
		ctx context.Context,
		setup *githubsetups.GitHubSetupDefinition,
	) (*githubsetups.GitHubSetupDefinition, error)

	// DeleteByID removes a GitHub setup by its ID within a transaction.
	DeleteByID(
		ctx context.Context,
		setupID shared.SetupID,
	) error

	// DeleteByListingID removes all GitHub setups for a listing within a transaction.
	DeleteByListingID(
		ctx context.Context,
		listingID shared.ListingID,
	) error
}

// / GitHubSetupCursor represents a stable pagination position
// / for iterating over GitHub setups.
type GitHubSetupCursor struct {
	LastCreatedAt time.Time
}

// GitHubSetupQueryRepository defines read-only operations for GitHub setups.
type GitHubSetupQueryRepository interface {
	GetByID(
		ctx context.Context,
		setupID shared.SetupID,
	) (*githubsetups.GitHubSetupDefinition, error)

	FetchNextBatchByListing(
		ctx context.Context,
		listingID shared.ListingID,
		request shared.BatchRequest[GitHubSetupCursor],
	) (setups []*githubsetups.GitHubSetupDefinition, nextCursor GitHubSetupCursor, err error)
}
