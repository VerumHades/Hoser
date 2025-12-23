package app

import (
	"fmt"

	"common/pkg/listing"
	githubsetups "common/pkg/setups/github"
)

// PublicListingService handles listing access rules, ownership checks, and GitHub setups.
type PublicListingService struct {
	listingFacade      *ListingFacadeService
	githubSetupService *githubsetups.GitHubSetupService
}

// NewPublicListingService creates a new PublicListingService.
func NewPublicListingService(listingFacade *ListingFacadeService, githubSetupService *githubsetups.GitHubSetupService) *PublicListingService {
	return &PublicListingService{
		listingFacade:      listingFacade,
		githubSetupService: githubSetupService,
	}
}

// =================== LISTING FUNCTIONS ===================

func (s *PublicListingService) ListByAuthor(authorID string) ([]*listing.Listing, error) {
	return s.listingFacade.ListByAuthor(authorID)
}

func (s *PublicListingService) SearchPublicListings(query string) ([]*listing.Listing, error) {
	listings, err := s.listingFacade.SearchListings(query)
	if err != nil {
		return nil, err
	}

	var publicListings []*listing.Listing
	for _, listingView := range listings {
		if listingView.AccessMode == listing.Public {
			publicListings = append(publicListings, listingView)
		}
	}

	return publicListings, nil
}

func (s *PublicListingService) GetOwnedListing(userID, listingID string) (*listing.Listing, error) {
	listingView, err := s.listingFacade.GetListing(listingID)
	if err != nil {
		return nil, err
	}

	if listingView.AuthorID != userID {
		return nil, fmt.Errorf("user %s is not the owner of listing %s", userID, listingID)
	}

	return listingView, nil
}

func (s *PublicListingService) GetPublicListing(listingID string) (*listing.Listing, error) {
	listingView, err := s.listingFacade.GetListing(listingID)
	if err != nil {
		return nil, err
	}

	if listingView.AccessMode != listing.Public {
		return nil, fmt.Errorf("listing %s is not public", listingID)
	}

	return listingView, nil
}

func (s *PublicListingService) CreateListing(authorID, title, description string, accessMode listing.ListingAccessMode) (*listing.Listing, error) {
	return s.listingFacade.CreateListing(authorID, title, description, accessMode)
}

func (s *PublicListingService) UpdateListing(userID, listingID string, update listing.ListingUpdate) (*listing.Listing, error) {
	if _, err := s.GetOwnedListing(userID, listingID); err != nil {
		return nil, err
	}
	return s.listingFacade.UpdateListing(listingID, update)
}

func (s *PublicListingService) DeleteListing(userID, listingID string) error {
	if _, err := s.GetOwnedListing(userID, listingID); err != nil {
		return err
	}
	return s.listingFacade.DeleteListing(listingID)
}

// =================== GITHUB SETUP FUNCTIONS ===================

// GetSetupForListing returns the GitHub setup definition for a listing (guarded).
func (s *PublicListingService) GetSetupForListing(userID, listingID string) (*githubsetups.GitHubSetupDefinition, error) {
	listingView, err := s.GetOwnedListing(userID, listingID)
	if err != nil {
		return nil, err
	}

	definition, err := s.githubSetupService.GetSetupByListing(listingView.ID)
	if err != nil {
		return nil, err
	}
	return definition, nil
}

// AttachOrUpdateSetup attaches or updates a GitHub setup for an owned listing.
func (s *PublicListingService) AttachOrUpdateSetup(userID, listingID, repoURL, accessToken string) (*githubsetups.GitHubSetupDefinition, error) {
	listingView, err := s.GetOwnedListing(userID, listingID)
	if err != nil {
		return nil, err
	}

	// Check if a setup already exists
	existingSetup, err := s.githubSetupService.GetSetupByListing(listingView.ID)
	if err != nil {
		def, err := s.githubSetupService.CreateSetup(listingID, repoURL, accessToken)

		if err != nil {
			return nil, err
		}

		return def, nil
	}

	def, err := s.githubSetupService.UpdateSetup(existingSetup.ID, &repoURL, &accessToken)

	if err != nil {
		return nil, err
	}

	return def, nil
}

// RemoveSetup removes a GitHub setup from an owned listing.
func (s *PublicListingService) RemoveSetup(userID, listingID string) error {
	if _, err := s.GetOwnedListing(userID, listingID); err != nil {
		return err
	}
	return s.githubSetupService.DeleteSetupsByListing(listingID)
}
