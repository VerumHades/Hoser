package listing

import "common/internal/shared"

// GitHubSetupRepository defines persistence operations for GitHub setups.
type GitHubSetupRepository interface {
	Save(setup *GitHubSetupDefinition) (*GitHubSetupDefinition, error)

	GetByID(setupID shared.SetupID) (*GitHubSetupDefinition, error)

	DeleteByID(setupID shared.SetupID) error
	DeleteByListingID(listingID shared.ListingID) error

	FetchNextBatchByListing(
		listingID shared.ListingID,
		lastSeenSetupID shared.SetupID,
		maximumBatchSize int,
	) ([]*GitHubSetupDefinition, error)
}
