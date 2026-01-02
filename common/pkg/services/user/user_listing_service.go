package userservices

import (
	"common/pkg/domain/listing"
	"common/pkg/domain/user"
	"common/pkg/shared"
	"context"
	"fmt"
)

// UserListingService provides operations for user interactions with listings.
type UserListingService struct {
	listingCommandRepository listing.ListingCommandRepository
	listingQueryRepository   listing.ListingQueryRepository
	savedCommandRepository   user.SavedListingCommandRepository
	savedQueryRepository     user.SavedListingQueryRepository
	listingSearchIndex       listing.ListingSearchIndex
}

// NewUserListingService constructs a new UserListingService instance.
func NewUserListingService(
	listingCommandRepository listing.ListingCommandRepository,
	listingQueryRepository listing.ListingQueryRepository,
	savedCommandRepository user.SavedListingCommandRepository,
	savedQueryRepository user.SavedListingQueryRepository,
	listingSearchIndex listing.ListingSearchIndex,
) *UserListingService {
	return &UserListingService{
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

	createdListing, err := s.savedCommandRepository.Create(ctx, nil, savedListing)
	if err != nil {
		return nil, fmt.Errorf("failed to save listing to library: %w", err)
	}

	return createdListing, nil
}

// FetchNextUserLibraryBatch retrieves the next batch of saved listings for a user.
func (s *UserListingService) FetchNextUserLibraryBatch(
	ctx context.Context,
	userID shared.UserID,
	batchRequest shared.BatchRequest[user.SavedListingCursor],
) ([]*user.SavedListing, user.SavedListingCursor, error) {
	return s.savedQueryRepository.FetchNextBatchByUser(ctx, userID, batchRequest)
}

// RemoveListingFromLibrary removes a listing from a user's library.
func (s *UserListingService) RemoveListingFromLibrary(
	ctx context.Context,
	userID shared.UserID,
	itemID shared.SavedListingID,
) error {

	savedListing, err := s.savedQueryRepository.GetByID(ctx, itemID)
	if err != nil {
		return fmt.Errorf("failed to fetch library item %s: %w", itemID, err)
	}

	if savedListing.UserID() != userID {
		return fmt.Errorf("library item %s does not belong to user %s", itemID, userID)
	}

	return s.savedCommandRepository.Delete(ctx, nil, itemID)
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
	batchRequest shared.BatchRequest[listing.ListingCursor],
) ([]*listing.Listing, listing.ListingCursor, error) {
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
