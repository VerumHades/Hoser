package listing

import (
	"common/internal/domain/listing"
	"common/internal/shared"
)

// DeveloperListingService exposes all operations including create, update, delete, and reindexing.
type DeveloperListingService struct {
	repository    listing.ListingRepository
	searchService listing.ListingSearchService
}

func NewDeveloperListingService(repo listing.ListingRepository, search listing.ListingSearchService) *DeveloperListingService {
	return &DeveloperListingService{
		repository:    repo,
		searchService: search,
	}
}

// CreateListing creates a listing and indexes it.
func (s *DeveloperListingService) CreateListing(listing *listing.Listing) error {
	if err := s.repository.Save(listing); err != nil {
		return err
	}

	return s.searchService.IndexListing(listing)
}

// UpdateListing updates a listing and reindexes it.
func (s *DeveloperListingService) UpdateListing(listing *listing.Listing) error {
	if err := s.repository.Save(listing); err != nil {
		return err
	}

	return s.searchService.IndexListing(listing)
}

// DeleteListing removes a listing and deletes it from the search index.
func (s *DeveloperListingService) DeleteListing(listingID shared.ListingID) error {
	listingView, err := s.repository.GetByID(listingID)
	if err != nil {
		return err
	}

	if err := s.repository.Delete(listingID); err != nil {
		return err
	}

	return s.searchService.RemoveListing(listingView.ID())
}

// GetSetupForListing returns the GitHub setup definition for a listing (guarded).
func (s *DeveloperListingService) GetSetupForListing(listingID string) (*githubsetups.GitHubSetupDefinition, error) {
	listingView, err := s.listingFacade.GetListing(listingID)
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
func (s *DeveloperListingService) AttachOrUpdateSetup(userID, listingID, repoURL, accessToken string) (*githubsetups.GitHubSetupDefinition, error) {
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
func (s *DeveloperListingService) RemoveSetup(userID, listingID string) error {
	if _, err := s.GetOwnedListing(userID, listingID); err != nil {
		return err
	}
	return s.githubSetupService.DeleteSetupsByListing(listingID)
}
