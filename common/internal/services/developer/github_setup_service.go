package listing

import (
	"common/internal/domain/listing"
	"common/internal/shared"
)

// DeveloperGitHubSetupService manages GitHub setups as full CRUD entities.
type DeveloperGitHubSetupService struct {
	githubSetupRepository listing.GitHubSetupRepository
}

// NewDeveloperGitHubSetupService constructs the service.
func NewDeveloperGitHubSetupService(
	githubSetupRepository listing.GitHubSetupRepository,
) *DeveloperGitHubSetupService {
	return &DeveloperGitHubSetupService{
		githubSetupRepository: githubSetupRepository,
	}
}

// CreateSetup creates a new GitHub setup for a listing.
func (s *DeveloperGitHubSetupService) CreateSetup(
	listingID shared.ListingID,
	repoURL string,
	accessToken string,
) (*listing.GitHubSetupDefinition, error) {
	setup, err := listing.NewGitHubSetup(listingID, repoURL, accessToken)
	if err != nil {
		return nil, err
	}
	return s.githubSetupRepository.Save(setup)
}

// UpdateSetup updates an existing GitHub setup by ID.
func (s *DeveloperGitHubSetupService) UpdateSetup(
	setupID shared.SetupID,
	repoURL string,
	accessToken string,
) (*listing.GitHubSetupDefinition, error) {
	existing, err := s.githubSetupRepository.GetByID(setupID)
	if err != nil {
		return nil, err
	}
	existing.UpdateRepo(repoURL, accessToken)
	return s.githubSetupRepository.Save(existing)
}

// GetSetupByID retrieves a GitHub setup by ID.
func (s *DeveloperGitHubSetupService) GetSetupByID(setupID shared.SetupID) (*listing.GitHubSetupDefinition, error) {
	return s.githubSetupRepository.GetByID(setupID)
}

// ListSetupsByListing lists all GitHub setups for a given listing.
func (s *DeveloperGitHubSetupService) ListSetupsByListing(
	listingID shared.ListingID,
	lastSeenSetupID shared.SetupID,
	maximumBatchSize int,
) ([]*listing.GitHubSetupDefinition, error) {
	return s.githubSetupRepository.FetchNextBatchByListing(listingID, lastSeenSetupID, maximumBatchSize)
}

// DeleteSetupByID removes a GitHub setup by ID.
func (s *DeveloperGitHubSetupService) DeleteSetupByID(setupID shared.SetupID) error {
	return s.githubSetupRepository.DeleteByID(setupID)
}

// DeleteSetupsByListing removes all GitHub setups for a listing in batches.
func (s *DeveloperGitHubSetupService) DeleteSetupsByListing(
	listingID shared.ListingID,
) error {
	const batchSize = 100
	for {
		setups, err := s.githubSetupRepository.FetchNextBatchByListing(listingID, shared.SetupID(""), batchSize)
		if err != nil {
			return err
		}
		if len(setups) == 0 {
			break
		}
		for _, setup := range setups {
			if err := s.githubSetupRepository.DeleteByID(setup.ID()); err != nil {
				return err
			}
		}
	}
	return nil
}
