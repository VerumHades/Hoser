package userservices

import (
	"common/pkg/domain/entities/listing"
	"common/pkg/domain/entities/user"
	"common/pkg/domain/repositories"
	"common/pkg/shared"
	"context"
	"fmt"
	"time"
)

// UserListingService provides operations for user interactions with listings.
type UserListingService struct {
	storageProvider          StorageProvider
	listingCommandRepository repositories.ListingCommandRepository
	listingQueryRepository   repositories.ListingQueryRepository
	savedCommandRepository   repositories.SavedListingCommandRepository
	savedQueryRepository     repositories.SavedListingQueryRepository
	listingSearchIndex       listing.ListingSearchIndex
}

type StorageProvider interface {
	GetReadLink(ctx context.Context, key string, expires time.Duration) (string, error)
}

// NewUserListingService constructs a new UserListingService contract.
func NewUserListingService(
	storageProvider StorageProvider,
	listingCommandRepository repositories.ListingCommandRepository,
	listingQueryRepository repositories.ListingQueryRepository,
	savedCommandRepository repositories.SavedListingCommandRepository,
	savedQueryRepository repositories.SavedListingQueryRepository,
	listingSearchIndex listing.ListingSearchIndex,
) *UserListingService {
	return &UserListingService{
		storageProvider:          storageProvider,
		listingCommandRepository: listingCommandRepository,
		listingQueryRepository:   listingQueryRepository,
		savedCommandRepository:   savedCommandRepository,
		savedQueryRepository:     savedQueryRepository,
		listingSearchIndex:       listingSearchIndex,
	}
}

//////////////////////////////
// User Library Operations  //
//////////////////////////////

// SaveListingToLibrary saves a public listing to a user's library.
func (s *UserListingService) SaveListingToLibrary(
	ctx context.Context,
	userID shared.UserID,
	listingID shared.ListingID,
) (*user.SavedListing, error) {

	listing, err := s.listingQueryRepository.GetByID(ctx, listingID)
	if err != nil {
		return nil, fmt.Errorf("listing %s not found: %w", listingID, err)
	}

	exists, err := s.savedQueryRepository.ExistsByUserAndListing(ctx, userID, listingID)
	if err != nil {
		return nil, fmt.Errorf("failed to check library existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("listing %s already saved in user %s library", listingID, userID)
	}

	savedListing, err := user.NewSavedListing(userID, listing.ID())
	if err != nil {
		return nil, err
	}

	createdListing, err := s.savedCommandRepository.Create(ctx, savedListing)
	if err != nil {
		return nil, fmt.Errorf("failed to save listing to library: %w", err)
	}

	return createdListing, nil
}

// FetchNextUserLibraryBatch retrieves the next batch of saved listings for a user.
func (s *UserListingService) FetchNextUserLibraryBatch(
	ctx context.Context,
	userID shared.UserID,
	batchRequest shared.BatchRequest[repositories.SavedListingCursor],
) ([]*user.SavedListing, repositories.SavedListingCursor, error) {
	return s.savedQueryRepository.FetchNextBatchByUser(ctx, userID, batchRequest)
}

// RemoveListingFromLibrary removes a listing from a user's library.
func (s *UserListingService) RemoveListingFromLibrary(
	ctx context.Context,
	userID shared.UserID,
	listingID shared.ListingID,
) error {
	savedListing, err := s.savedQueryRepository.GetByUserAndListing(ctx, userID, listingID)
	if err != nil {
		return fmt.Errorf("failed to fetch library item %s: %w", listingID, err)
	}

	return s.savedCommandRepository.Delete(ctx, savedListing.ID())
}

//////////////////////////////
// Listing Retrieval        //
//////////////////////////////

// GetListing retrieves a single listing by ID.
func (s *UserListingService) GetListing(
	ctx context.Context,
	listingID shared.ListingID,
) (*listing.Listing, error) {
	return s.listingQueryRepository.GetByID(ctx, listingID)
}

// FetchNextListingsByAuthor retrieves listings by a specific author in batches.
func (s *UserListingService) FetchNextListingsByAuthor(
	ctx context.Context,
	authorID shared.UserID,
	batchRequest shared.BatchRequest[repositories.ListingCursor],
) ([]*listing.Listing, repositories.ListingCursor, error) {
	return s.listingQueryRepository.FetchNextBatchByAuthor(ctx, authorID, batchRequest)
}

// SearchNextListingsBatch retrieves listings matching a query in batches.
func (s *UserListingService) SearchNextListingsBatch(
	ctx context.Context,
	query string,
	batchRequest shared.BatchRequest[listing.ListingSearchCursor],
) ([]*listing.Listing, listing.ListingSearchCursor, error) {
	return s.listingSearchIndex.SearchNextBatch(ctx, query, batchRequest)
}

func (s *UserListingService) GetListingScreenshotReadUrl(
	ctx context.Context,
	listingID shared.ListingID,
	screenshotID shared.ListingScreenshotID,
) (string, error) {
	_, err := s.listingQueryRepository.GetByID(ctx, listingID)
	if err != nil {
		return "", err
	}

	fileKey := shared.GenerateFileKey(listingID, screenshotID)
	return s.storageProvider.GetReadLink(ctx, fileKey, 1*time.Hour)
}
