package app

import (
	"common/pkg/library"
	"common/pkg/listing"
	"common/pkg/user"
	"fmt"
)

// UserAppService handles cross-domain operations involving users.
type UserAppService struct {
	userService    *user.UserService
	listingFacade  *ListingFacadeService
	libraryService *library.LibraryService
}

// NewUserAppService creates a new instance of UserAppService.
func NewUserAppService(
	userService *user.UserService,
	listingFacade *ListingFacadeService,
	libraryService *library.LibraryService,
) *UserAppService {
	return &UserAppService{
		userService:    userService,
		listingFacade:  listingFacade,
		libraryService: libraryService,
	}
}

// ListUserListings returns all listings authored by a given user.
func (s *UserAppService) ListUserListings(userID string) ([]*listing.Listing, error) {
	return s.listingFacade.ListByAuthor(userID)
}

// CreateListing creates a new listing authored by the user.
func (s *UserAppService) CreateListing(authorID, title, description string, accessMode listing.ListingAccessMode) (*listing.Listing, error) {
	if _, err := s.userService.GetUser(authorID); err != nil {
		return nil, fmt.Errorf("user %s does not exist", authorID)
	}
	listingView, err := s.listingFacade.CreateListing(authorID, title, description, accessMode)
	if err != nil {
		return nil, err
	}
	return listingView, nil
}

// UpdateUserListing updates a listing authored by the user.
// Only the owner can update their listings.
func (s *UserAppService) UpdateUserListing(userID, listingID string, update listing.ListingUpdate) (*listing.Listing, error) {
	// Fetch the listing to check ownership
	_, err := s.GetUserListing(userID, listingID)
	if err != nil {
		return nil, err
	}

	listingView, err := s.listingFacade.UpdateListing(listingID, update)
	if err != nil {
		return nil, err
	}

	return listingView, nil
}

// UpdateUserListing updates a listing authored by the user.
// Only the owner can update their listings.
func (s *UserAppService) DeleteUserListing(userID, listingID string) error {
	// Fetch the listing to check ownership
	_, err := s.GetUserListing(userID, listingID)
	if err != nil {
		return err
	}

	return s.listingFacade.DeleteListing(listingID)
}

// SaveListingToLibrary adds a listing to a user's library.
// SaveListingToLibrary adds a listing to a user's library.
// It ensures the user exists and the listing is public.
func (s *UserAppService) SaveListingToLibrary(userID string, listingID string) (*library.SavedListing, error) {
	// Check that the user exists
	if _, err := s.userService.GetUser(userID); err != nil {
		return nil, fmt.Errorf("user %s does not exist", userID)
	}

	// Check that the listing exists and is public
	listingItem, err := s.GetPublicListing(listingID)
	if err != nil {
		return nil, fmt.Errorf("listing %s is not public or does not exist", listingID)
	}

	// Save to library
	savedItem, err := s.libraryService.SaveListing(userID, listingItem.ID)
	if err != nil {
		return nil, err
	}

	return savedItem, nil
}

// ListUserLibrary returns all saved listings for a user.
func (s *UserAppService) ListUserLibrary(userID string) ([]*library.SavedListing, error) {
	return s.libraryService.ListUserLibrary(userID)
}

// RemoveListingFromLibrary removes a saved listing from a user's library.
// Only the owner of the library item can remove it.
func (s *UserAppService) RemoveListingFromLibrary(userID, itemID string) error {
	// Fetch the user's library
	userLibrary, err := s.libraryService.ListUserLibrary(userID)
	if err != nil {
		return fmt.Errorf("failed to fetch user's library: %w", err)
	}

	// Verify that the item belongs to this user
	var found bool
	for _, item := range userLibrary {
		if item.ID == itemID {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("library item %s does not belong to user %s", itemID, userID)
	}

	// Remove the item
	return s.libraryService.RemoveListing(itemID)
}
func (s *UserAppService) IsUserDeveloper(userID string) (bool, error) {
	return s.userService.CheckDeveloper(userID)
}

func (s *UserAppService) GetUser(userID string) (*user.User, error) {
	return s.userService.GetUser(userID)
}

func (s *UserAppService) GetUserListings(userID string) ([]*listing.Listing, error) {
	return s.listingFacade.ListByAuthor(userID)
}

func (s *UserAppService) GetUserListing(userID string, listingID string) (*listing.Listing, error) {
	// Fetch the listing to check ownership
	listingView, err := s.listingFacade.GetListing(listingID)
	if err != nil {
		return nil, err
	}

	if listingView.AuthorID != userID {
		return nil, fmt.Errorf("user %s is not the owner of listing %s", userID, listingID)
	}

	return listingView, nil
}

func (s *UserAppService) GetPublicListing(listingID string) (*listing.Listing, error) {
	// Fetch the listing to check ownership
	listingView, err := s.listingFacade.GetListing(listingID)
	if err != nil {
		return nil, err
	}
	if listingView.AccessMode != listing.Public {
		return nil, fmt.Errorf("listing is not public.")
	}

	return listingView, nil
}
