package listing

import (
	"common/pkg/domain/listing"
	"common/pkg/shared"
	"common/pkg/util"
	"context"
)

// DeveloperGitHubSetupService manages GitHub setups as full CRUD entities.
type DeveloperGitHubSetupService struct {
	commandRepository listing.GitHubSetupCommandRepository
	queryRepository   listing.GitHubSetupQueryRepository
}

// NewDeveloperGitHubSetupService constructs the service.
func NewDeveloperGitHubSetupService(
	commandRepository listing.GitHubSetupCommandRepository,
	queryRepository listing.GitHubSetupQueryRepository,
) *DeveloperGitHubSetupService {
	return &DeveloperGitHubSetupService{
		commandRepository: commandRepository,
		queryRepository:   queryRepository,
	}
}

// CreateSetup creates a new GitHub setup for a listing.
func (s *DeveloperGitHubSetupService) CreateSetup(
	ctx context.Context,
	listingID shared.ListingID,
	repoURL string,
	accessToken string,
) (*listing.GitHubSetupDefinition, error) {
	setup, err := listing.NewGitHubSetup(listingID, repoURL, accessToken)
	if err != nil {
		return nil, err
	}
	return s.commandRepository.Create(ctx, nil, setup)
}

// UpdateSetup updates an existing GitHub setup by ID.
func (s *DeveloperGitHubSetupService) UpdateSetup(
	ctx context.Context,
	setupID shared.SetupID,
	repoURL string,
	accessToken string,
) (*listing.GitHubSetupDefinition, error) {
	existing, err := s.queryRepository.GetByID(ctx, setupID)
	if err != nil {
		return nil, err
	}
	existing.UpdateRepo(repoURL, accessToken)
	return s.commandRepository.Update(ctx, nil, existing)
}

// GetSetupByID retrieves a GitHub setup by ID.
func (s *DeveloperGitHubSetupService) GetSetupByID(
	ctx context.Context,
	setupID shared.SetupID,
) (*listing.GitHubSetupDefinition, error) {
	return s.queryRepository.GetByID(ctx, setupID)
}

// ListSetupsByListing lists all GitHub setups for a given listing in batches.
func (s *DeveloperGitHubSetupService) ListSetupsByListing(
	ctx context.Context,
	listingID shared.ListingID,
	batchRequest shared.BatchRequest,
) ([]*listing.GitHubSetupDefinition, shared.Cursor, error) {
	return s.queryRepository.FetchNextBatchByListing(ctx, listingID, batchRequest)
}

// DeleteSetupByID removes a GitHub setup by ID.
func (s *DeveloperGitHubSetupService) DeleteSetupByID(
	ctx context.Context,
	setupID shared.SetupID,
) error {
	return s.commandRepository.DeleteByID(ctx, nil, setupID)
}

// DeleteSetupsByListing removes all GitHub setups for a listing in batches.
func (s *DeveloperGitHubSetupService) DeleteSetupsByListing(
	ctx context.Context,
	listingID shared.ListingID,
) error {
	return util.ProcessInBatches(
		ctx,
		100,
		func(ctx context.Context, request shared.BatchRequest) (items []*listing.GitHubSetupDefinition, nextCursor shared.Cursor, err error) {
			return s.queryRepository.FetchNextBatchByListing(ctx, listingID, request)
		},
		func(ctx context.Context, setup *listing.GitHubSetupDefinition) error {
			if err := s.commandRepository.DeleteByID(ctx, nil, setup.ID()); err != nil {
				return err
			}
			return nil
		},
	)
}
