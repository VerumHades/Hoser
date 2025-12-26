package user

import (
	"common/internal/domain/listing"
	"common/internal/domain/user"
	"common/internal/shared"
	"fmt"
)

type UserListingService struct {
	listingRepository      listing.ListingRepository
	savedListingRepository user.SavedListingRepository
	searchService          listing.ListingSearchService
}

func NewUserListingService(
	listingRepository listing.ListingRepository,
	savedListingRepository user.SavedListingRepository,
	searchService listing.ListingSearchService,
) *UserListingService {
	return &UserListingService{
		listingRepository:      listingRepository,
		savedListingRepository: savedListingRepository,
		searchService:          searchService,
	}
}

// SaveListingToLibrary saves a public listing to the user's library.
func (s *UserListingService) SaveListingToLibrary(userID shared.UserID, listingID shared.ListingID) (*user.SavedListing, error) {
	// Check if listing exists
	listing, err := s.listingRepository.GetByID(listingID)
	if err != nil {
		return nil, fmt.Errorf("listing %s not found: %w", listingID, err)
	}

	// Prevent duplicate entries
	exists, err := s.savedListingRepository.ExistsByUserAndListing(userID, listingID)
	if err != nil {
		return nil, fmt.Errorf("failed to check library existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("listing %s already saved in user %s library", listingID, userID)
	}

	// Save listing to library
	item, err := user.NewSavedListing(userID, listing.ID())
	if err != nil {
		return nil, err
	}
	item, err = s.savedListingRepository.Save(item)
	if err != nil {
		return nil, fmt.Errorf("failed to save listing to library: %w", err)
	}

	return item, nil
}

func (s *UserListingService) FetchNextUserLibraryBatch(
	userID shared.UserID,
	lastSeenSavedListingID shared.SavedListingID,
	maximumBatchSize int,
) ([]*user.SavedListing, error) {
	return s.savedListingRepository.FetchNextBatchByUser(
		userID,
		lastSeenSavedListingID,
		maximumBatchSize,
	)
}

// RemoveListingFromLibrary removes a listing from a user's library.
func (s *UserListingService) RemoveListingFromLibrary(userID shared.UserID, itemID shared.SavedListingID) error {
	// Fetch the library item
	item, err := s.savedListingRepository.GetByID(itemID)
	if err != nil {
		return fmt.Errorf("failed to fetch library item %s: %w", itemID, err)
	}

	// Verify ownership
	if item.UserID() != userID {
		return fmt.Errorf("library item %s does not belong to user %s", itemID, userID)
	}

	return s.savedListingRepository.Delete(itemID)
}

// GetListing returns a single listing by ID.
func (s *UserListingService) GetListing(listingID shared.ListingID) (*listing.Listing, error) {
	return s.listingRepository.GetByID(listingID)
}

// ListByAuthor returns listings authored by a specific user.
func (s *UserListingService) FetchNextListingsByAuthor(
	authorID shared.UserID,
	lastSeenListingID shared.ListingID,
	maximumBatchSize int,
) ([]*listing.Listing, error) {
	return s.listingRepository.FetchNextBatchByAuthor(
		authorID,
		lastSeenListingID,
		maximumBatchSize,
	)
}

func (s *UserListingService) SearchNextListingsBatch(
	query string,
	lastSeenListingID shared.ListingID,
	maximumBatchSize int,
) ([]*listing.Listing, error) {
	return s.searchService.SearchNextBatch(
		query,
		lastSeenListingID,
		maximumBatchSize,
	)
}
