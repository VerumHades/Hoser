package listing

import (
	"fmt"
	"time"

	"common/pkg/util"
)

// =================== MODEL ===================

// GitHubSetupDefinition represents a repository configuration that can be attached to a listing.
// Each listing can have zero or more setups (future-proofed).
type GitHubSetupDefinition struct {
	ID            string
	ListingID     string
	GitHubRepoURL string
	AccessToken   string // optional, for private repos
	CreatedAt     time.Time
	UpdatedAt     time.Time
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
		ID:            util.GenerateUUID(),
		ListingID:     listingID,
		GitHubRepoURL: repoURL,
		AccessToken:   accessToken,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := s.repo.Save(setup); err != nil {
		return nil, fmt.Errorf("failed to save GitHub setup: %w", err)
	}
	return setup, nil
}

// UpdateSetup updates an existing GitHub setup. Only non-empty fields are updated.
func (s *GitHubSetupService) UpdateSetup(id string, repoURL *string, accessToken *string) (*GitHubSetupDefinition, error) {
	setup, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch GitHub setup: %w", err)
	}
	if setup == nil {
		return nil, fmt.Errorf("GitHub setup with ID %s not found", id)
	}

	if repoURL != nil && *repoURL != "" {
		setup.GitHubRepoURL = *repoURL
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

// GetSetupByID returns a single GitHub setup by ID.
func (s *GitHubSetupService) GetSetupByID(id string) (*GitHubSetupDefinition, error) {
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}
	setup, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch GitHub setup: %w", err)
	}
	if setup == nil {
		return nil, fmt.Errorf("GitHub setup with ID %s not found", id)
	}
	return setup, nil
}

// ListSetupsByListing returns all GitHub setups attached to a listing.
func (s *GitHubSetupService) ListSetupsByListing(listingID string) ([]*GitHubSetupDefinition, error) {
	if listingID == "" {
		return nil, fmt.Errorf("listingID cannot be empty")
	}
	setups, err := s.repo.GetByListingID(listingID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch setups for listing %s: %w", listingID, err)
	}
	return setups, nil
}

// DeleteSetupByID deletes a setup by ID.
func (s *GitHubSetupService) DeleteSetupByID(id string) error {
	if id == "" {
		return fmt.Errorf("id cannot be empty")
	}
	return s.repo.DeleteByID(id)
}

// DeleteSetupsByListing deletes all setups attached to a listing.
func (s *GitHubSetupService) DeleteSetupsByListing(listingID string) error {
	if listingID == "" {
		return fmt.Errorf("listingID cannot be empty")
	}
	return s.repo.DeleteByListingID(listingID)
}
