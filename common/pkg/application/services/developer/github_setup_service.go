package developer

import (
	"common/pkg/domain/entities/githubsetups"
	"common/pkg/domain/repositories"
	"common/pkg/shared"
	"common/pkg/util"
	"context"
)

// DeveloperGitHubSetupService manages GitHub setups as full CRUD entities.
type DeveloperGitHubSetupService struct {
	commandRepository repositories.GitHubSetupCommandRepository
	queryRepository   repositories.GitHubSetupQueryRepository
}

// NewDeveloperGitHubSetupService constructs the service.
func NewDeveloperGitHubSetupService(
	commandRepository repositories.GitHubSetupCommandRepository,
	queryRepository repositories.GitHubSetupQueryRepository,
) *DeveloperGitHubSetupService {
	return &DeveloperGitHubSetupService{
		commandRepository: commandRepository,
		queryRepository:   queryRepository,
	}
}

// CreateSetup creates a new GitHub setup for a githubsetups.
func (s *DeveloperGitHubSetupService) CreateSetup(
	ctx context.Context,
	listingID shared.ListingID,
	repoURL string,
	accessToken string,
) (*githubsetups.GitHubSetupDefinition, error) {
	setup, err := githubsetups.NewGitHubSetup(listingID, repoURL, accessToken)
	if err != nil {
		return nil, err
	}
	return s.commandRepository.Create(ctx, setup)
}

// UpdateSetup updates an existing GitHub setup by ID.
func (s *DeveloperGitHubSetupService) UpdateSetup(
	ctx context.Context,
	setupID shared.GithubSetupID,
	repoURL string,
	accessToken string,
) (*githubsetups.GitHubSetupDefinition, error) {
	existing, err := s.queryRepository.GetByID(ctx, setupID)
	if err != nil {
		return nil, err
	}
	existing.UpdateRepo(repoURL, accessToken)
	return s.commandRepository.Update(ctx, existing)
}

// GetSetupByID retrieves a GitHub setup by ID.
func (s *DeveloperGitHubSetupService) GetSetupByID(
	ctx context.Context,
	setupID shared.GithubSetupID,
) (*githubsetups.GitHubSetupDefinition, error) {
	return s.queryRepository.GetByID(ctx, setupID)
}

// ListSetupsByListing lists all GitHub setups for a given listing in batches.
func (s *DeveloperGitHubSetupService) ListSetupsByListing(
	ctx context.Context,
	listingID shared.ListingID,
	batchRequest shared.BatchRequest[repositories.GitHubSetupCursor],
) ([]*githubsetups.GitHubSetupDefinition, repositories.GitHubSetupCursor, error) {
	return s.queryRepository.FetchNextBatchByListing(ctx, listingID, batchRequest)
}

// DeleteSetupByID removes a GitHub setup by ID.
func (s *DeveloperGitHubSetupService) DeleteSetupByID(
	ctx context.Context,
	setupID shared.GithubSetupID,
) error {
	return s.commandRepository.DeleteByID(ctx, setupID)
}

func (s *DeveloperGitHubSetupService) DeleteSetupsByListing(
	ctx context.Context,
	listingID shared.ListingID,
) error {
	fetchNextBatch := func(
		ctx context.Context,
		request shared.BatchRequest[repositories.GitHubSetupCursor],
	) (items []*githubsetups.GitHubSetupDefinition, nextCursor repositories.GitHubSetupCursor, err error) {
		return s.queryRepository.FetchNextBatchByListing(ctx, listingID, request)
	}

	for setup := range util.GenerateInBatches(ctx, 100, fetchNextBatch) {
		if err := s.commandRepository.DeleteByID(ctx, setup.ID()); err != nil {
			return err
		}
	}

	return nil
}
