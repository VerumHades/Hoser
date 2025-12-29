package listing

import (
	"common/pkg/domain/listing"
	"common/pkg/shared"
	"common/pkg/util"
	"context"
)

// DeveloperListingService exposes operations for managing listings and GitHub setups.
type DeveloperListingService struct {
	transactionProvider shared.TransactionProvider

	listingRepository      listing.ListingCommandRepository
	listingQueryRepository listing.ListingQueryRepository
	searchIndex            listing.ListingSearchIndex
	githubSetupRepository  listing.GitHubSetupCommandRepository
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
func (s *DeveloperListingService) CreateListing(ctx context.Context, listing *listing.Listing) error {
	err := s.listingRepository.Create(ctx, nil, listing)
	if err != nil {
		return err
	}
	if err := s.searchIndex.Index(ctx, listing); err != nil {
		return err
	}
	return nil
}

// UpdateListing updates a listing and reindexes it.
func (s *DeveloperListingService) UpdateListing(ctx context.Context, listingID shared.ListingID, updateFunction func(listing *listing.Listing) error) error {
	existingListing, err := util.GetExistingEntity(ctx, listingID, s.listingQueryRepository)
	if err != nil {
		return err
	}

	return shared.WithTransaction(ctx, s.transactionProvider, func(transaction shared.Transaction) error {
		if err := updateFunction(existingListing); err != nil {
			return err
		}

		if err := s.listingRepository.Update(ctx, transaction, existingListing); err != nil {
			return err
		}

		if err := s.searchIndex.Index(ctx, existingListing); err != nil {
			return err
		}

		return nil
	})
}

// DeleteListing removes a listing and deletes it from the search index.
func (s *DeveloperListingService) DeleteListing(ctx context.Context, listingID shared.ListingID) error {
	err := s.listingRepository.Delete(ctx, nil, listingID)
	if err != nil {
		return err
	}
	return s.searchIndex.Remove(ctx, listingID)
}
