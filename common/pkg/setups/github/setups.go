package githubsetups

import (
	"fmt"
	"time"

	"common/pkg/util"
)

// =================== MODEL ===================

// GitHubSetupDefinition represents a repository configuration that can be attached to a listing.
// Each listing can have zero or more setups (future-proofed).
type GitHubSetupDefinition struct {
	ID          string
	ListingID   string
	RepoURL     string
	AccessToken string // optional, for private repos
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// =================== REPOSITORY ===================

// GitHubSetupRepository defines persistence operations for GitHub setups.
type GitHubSetupRepository interface {
	Save(setup *GitHubSetupDefinition) error
	GetByID(id string) (*GitHubSetupDefinition, error)
	GetByListingID(listingID string) ([]*GitHubSetupDefinition, error)
	DeleteByID(id string) error
	DeleteByListingID(listingID string) error
}

// =================== SERVICE ===================

// GitHubSetupService provides higher-level business logic around GitHub setups.
type GitHubSetupService struct {
	repo GitHubSetupRepository
}

// NewGitHubSetupService creates a new service instance.
func NewGitHubSetupService(repo GitHubSetupRepository) *GitHubSetupService {
	return &GitHubSetupService{repo: repo}
}

// CreateSetup attaches a new GitHub setup to a listing.
func (s *GitHubSetupService) CreateSetup(listingID, repoURL, accessToken string) (*GitHubSetupDefinition, error) {
	if listingID == "" {
		return nil, fmt.Errorf("listingID cannot be empty")
	}
	if repoURL == "" {
		return nil, fmt.Errorf("GitHubRepoURL cannot be empty")
	}

	now := time.Now()
	setup := &GitHubSetupDefinition{
		ID:          util.GenerateUUID(),
		ListingID:   listingID,
		RepoURL:     repoURL,
		AccessToken: accessToken,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.Save(setup); err != nil {
		return nil, fmt.Errorf("failed to save GitHub setup: %w", err)
	}
	return setup, nil
}

// ensureSetupExists fetches a setup by ID and errors if it does not exist.
func (s *GitHubSetupService) ensureSetupExists(id string) (*GitHubSetupDefinition, error) {
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}
	setup, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch GitHub setup: %w", err)
	}
	if setup == nil {
		return nil, fmt.Errorf("GitHub setup with ID %s does not exist", id)
	}
	return setup, nil
}

// ensureSetupsExistForListing fetches setups for a listing and errors if none exist.
func (s *GitHubSetupService) ensureSetupsExistForListing(listingID string) ([]*GitHubSetupDefinition, error) {
	if listingID == "" {
		return nil, fmt.Errorf("listingID cannot be empty")
	}
	setups, err := s.repo.GetByListingID(listingID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch GitHub setups for listing %s: %w", listingID, err)
	}
	if len(setups) == 0 {
		return nil, fmt.Errorf("no GitHub setups exist for listing %s", listingID)
	}
	return setups, nil
}

// UpdateSetup updates an existing GitHub setup.
func (s *GitHubSetupService) UpdateSetup(id string, repoURL *string, accessToken *string) (*GitHubSetupDefinition, error) {
	setup, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if repoURL != nil && *repoURL != "" {
		setup.RepoURL = *repoURL
	}
	if accessToken != nil {
		setup.AccessToken = *accessToken
	}
	setup.UpdatedAt = time.Now()

	if err := s.repo.Save(setup); err != nil {
		return nil, fmt.Errorf("failed to save GitHub setup: %w", err)
	}
	return setup, nil
}

// GetSetupByListing retrieves the single GitHub setup for a listing.
// Returns an error if none exist.
func (s *GitHubSetupService) GetSetupByListing(listingID string) (*GitHubSetupDefinition, error) {
	if listingID == "" {
		return nil, fmt.Errorf("listingID cannot be empty")
	}
	setups, err := s.repo.GetByListingID(listingID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch GitHub setup for listing %s: %w", listingID, err)
	}
	if len(setups) == 0 {
		return nil, fmt.Errorf("no GitHub setup exists for listing %s", listingID)
	}
	// enforce single setup per listing
	return setups[0], nil
}

// GetSetupByID returns a single GitHub setup by ID.
func (s *GitHubSetupService) GetSetupByID(id string) (*GitHubSetupDefinition, error) {
	return s.ensureSetupExists(id)
}

// ListSetupsByListing returns all GitHub setups attached to a listing.
func (s *GitHubSetupService) ListSetupsByListing(listingID string) ([]*GitHubSetupDefinition, error) {
	return s.ensureSetupsExistForListing(listingID)
}

// DeleteSetupByID deletes a setup by ID.
func (s *GitHubSetupService) DeleteSetupByID(id string) error {
	_, err := s.ensureSetupExists(id)
	if err != nil {
		return err
	}
	return s.repo.DeleteByID(id)
}

// DeleteSetupsByListing deletes all setups attached to a listing.
func (s *GitHubSetupService) DeleteSetupsByListing(listingID string) error {
	_, err := s.ensureSetupsExistForListing(listingID)
	if err != nil {
		return err
	}
	return s.repo.DeleteByListingID(listingID)
}
