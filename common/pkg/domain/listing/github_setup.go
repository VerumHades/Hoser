package listing

import (
	"common/pkg/shared"
	"errors"
	"time"
)

// GitHubSetupDefinition represents a GitHub repository configuration attached to a listing.
type GitHubSetupDefinition struct {
	id          shared.SetupID
	listingID   shared.ListingID
	repoURL     string
	accessToken string
	createdAt   time.Time
	updatedAt   time.Time
}

// NewGitHubSetup creates a new setup with a generated ID and timestamps.
func NewGitHubSetup(listingID shared.ListingID, repoURL string, accessToken string) (*GitHubSetupDefinition, error) {
	id := shared.SetupID(shared.GenerateUUID())
	now := time.Now().UTC()
	return NewGitHubSetupWithID(id, listingID, repoURL, accessToken, now, now)
}

// NewGitHubSetupWithID creates a setup with an existing ID and timestamps (for repository hydration).
func NewGitHubSetupWithID(id shared.SetupID, listingID shared.ListingID, repoURL string, accessToken string, createdAt, updatedAt time.Time) (*GitHubSetupDefinition, error) {
	setup := &GitHubSetupDefinition{
		id:          id,
		listingID:   listingID,
		repoURL:     repoURL,
		accessToken: accessToken,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
	if err := setup.Validate(); err != nil {
		return nil, err
	}
	return setup, nil
}

// ID returns the setup's unique identifier.
func (s *GitHubSetupDefinition) ID() shared.SetupID {
	return s.id
}

// shared.ListingID returns the associated listing ID.
func (s *GitHubSetupDefinition) ListingID() shared.ListingID {
	return s.listingID
}

// RepoURL returns the repository URL.
func (s *GitHubSetupDefinition) RepoURL() string {
	return s.repoURL
}

// AccessToken returns the access token (if any).
func (s *GitHubSetupDefinition) AccessToken() string {
	return s.accessToken
}

// CreatedAt returns the timestamp when the setup was created.
func (s *GitHubSetupDefinition) CreatedAt() time.Time {
	return s.createdAt
}

// UpdatedAt returns the timestamp when the setup was last updated.
func (s *GitHubSetupDefinition) UpdatedAt() time.Time {
	return s.updatedAt
}

// UpdateRepo updates the repository URL and optionally the access token, updating the timestamp.
func (s *GitHubSetupDefinition) UpdateRepo(repoURL string, accessToken string) error {
	if repoURL == "" {
		return errors.New("repo URL cannot be empty")
	}
	s.repoURL = repoURL
	s.accessToken = accessToken
	s.updatedAt = time.Now().UTC()
	return nil
}

// Validate ensures the setup entity is in a valid state.
func (s *GitHubSetupDefinition) Validate() error {
	if s.id == "" {
		return errors.New("setup ID cannot be empty")
	}
	if s.listingID == "" {
		return errors.New("listing ID cannot be empty")
	}
	if s.repoURL == "" {
		return errors.New("repo URL cannot be empty")
	}
	if s.createdAt.IsZero() {
		return errors.New("createdAt timestamp cannot be zero")
	}
	if s.updatedAt.IsZero() {
		return errors.New("updatedAt timestamp cannot be zero")
	}
	return nil
}
