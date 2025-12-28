package listing

import (
	"common/internal/domain/listing"
	"common/internal/shared"
	"context"
)

// DeveloperListingService exposes operations for managing listings and GitHub setups.
type DeveloperListingService struct {
	listingRepository     listing.ListingCommandRepository
	searchIndex           listing.ListingSearchIndex
	githubSetupRepository listing.GitHubSetupCommandRepository
}

// NewDeveloperListingService constructs a new DeveloperListingService.
func NewDeveloperListingService(
	listingRepository listing.ListingCommandRepository,
	searchService listing.ListingSearchIndex,
	githubSetupRepository listing.GitHubSetupCommandRepository,
) *DeveloperListingService {
	return &DeveloperListingService{
		listingRepository:     listingRepository,
		searchIndex:           searchService,
		githubSetupRepository: githubSetupRepository,
	}
}

// CreateListing creates a listing, persists it, and indexes it.
func (s *DeveloperListingService) CreateListing(ctx context.Context, listing *listing.Listing) (*listing.Listing, error) {
	savedListing, err := s.listingRepository.Create(ctx, nil, listing)
	if err != nil {
		return nil, err
	}
	if err := s.searchIndex.Index(ctx, savedListing); err != nil {
		return nil, err
	}
	return savedListing, nil
}

// UpdateListing updates a listing and reindexes it.
func (s *DeveloperListingService) UpdateListing(ctx context.Context, listing *listing.Listing) (*listing.Listing, error) {
	savedListing, err := s.listingRepository.Update(ctx, nil, listing)
	if err != nil {
		return nil, err
	}
	if err := s.searchIndex.Index(ctx, savedListing); err != nil {
		return nil, err
	}
	return savedListing, nil
}

// DeleteListing removes a listing and deletes it from the search index.
func (s *DeveloperListingService) DeleteListing(ctx context.Context, listingID shared.ListingID) error {
	err := s.listingRepository.Delete(ctx, nil, listingID)
	if err != nil {
		return err
	}
	return s.searchIndex.Remove(ctx, listingID)
}
