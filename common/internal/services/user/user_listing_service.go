package user

import (
	"common/internal/domain/listing"
	"common/internal/domain/user"
	"common/internal/shared"
	"context"
	"fmt"
)

type UserListingService struct {
	listingCommandRepo listing.ListingCommandRepository
	listingQueryRepo   listing.ListingQueryRepository
	savedCommandRepo   user.SavedListingCommandRepository
	savedQueryRepo     user.SavedListingQueryRepository
	searchIndex        listing.ListingSearchIndex
}

func NewUserListingService(
	listingCommandRepo listing.ListingCommandRepository,
	listingQueryRepo listing.ListingQueryRepository,
	savedCommandRepo user.SavedListingCommandRepository,
	savedQueryRepo user.SavedListingQueryRepository,
	searchIndex listing.ListingSearchIndex,
) *UserListingService {
	return &UserListingService{
		listingCommandRepo: listingCommandRepo,
		listingQueryRepo:   listingQueryRepo,
		savedCommandRepo:   savedCommandRepo,
		savedQueryRepo:     savedQueryRepo,
		searchIndex:        searchIndex,
	}
}

// SaveListingToLibrary saves a public listing to a user's library.
func (s *UserListingService) SaveListingToLibrary(ctx context.Context, userID shared.UserID, listingID shared.ListingID) (*user.SavedListing, error) {
	listing, err := s.listingQueryRepo.GetByID(ctx, listingID)
	if err != nil {
		return nil, fmt.Errorf("listing %s not found: %w", listingID, err)
	}

	exists, err := s.savedQueryRepo.ExistsByUserAndListing(ctx, userID, listingID)
	if err != nil {
		return nil, fmt.Errorf("failed to check library existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("listing %s already saved in user %s library", listingID, userID)
	}

	item, err := user.NewSavedListing(userID, listing.ID())
	if err != nil {
		return nil, err
	}

	savedItem, err := s.savedCommandRepo.Create(ctx, nil, item)
	if err != nil {
		return nil, fmt.Errorf("failed to save listing to library: %w", err)
	}

	return savedItem, nil
}

// FetchNextUserLibraryBatch retrieves the next batch of saved listings for a user.
func (s *UserListingService) FetchNextUserLibraryBatch(ctx context.Context, userID shared.UserID, batchRequest shared.BatchRequest) ([]*user.SavedListing, shared.Cursor, error) {
	return s.savedQueryRepo.FetchNextBatchByUser(ctx, userID, batchRequest)
}

// RemoveListingFromLibrary removes a listing from a user's library.
func (s *UserListingService) RemoveListingFromLibrary(ctx context.Context, userID shared.UserID, itemID shared.SavedListingID, transaction shared.Transaction) error {
	item, err := s.savedQueryRepo.GetByID(ctx, itemID)
	if err != nil {
		return fmt.Errorf("failed to fetch library item %s: %w", itemID, err)
	}

	if item.UserID() != userID {
		return fmt.Errorf("library item %s does not belong to user %s", itemID, userID)
	}

	return s.savedCommandRepo.Delete(ctx, transaction, itemID)
}

// GetListing retrieves a single listing by ID.
func (s *UserListingService) GetListing(ctx context.Context, listingID shared.ListingID) (*listing.Listing, error) {
	return s.listingQueryRepo.GetByID(ctx, listingID)
}

// FetchNextListingsByAuthor retrieves listings by a specific author in batches.
func (s *UserListingService) FetchNextListingsByAuthor(ctx context.Context, authorID shared.UserID, batchRequest shared.BatchRequest) ([]*listing.Listing, shared.Cursor, error) {
	return s.listingQueryRepo.FetchNextBatchByAuthor(ctx, authorID, batchRequest)
}

// SearchNextListingsBatch retrieves listings matching a query in batches.
func (s *UserListingService) SearchNextListingsBatch(ctx context.Context, query string, batchRequest shared.BatchRequest) ([]*listing.Listing, shared.Cursor, error) {
	return s.searchIndex.SearchNextBatch(ctx, query, batchRequest)
}
