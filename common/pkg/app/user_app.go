package app

import (
	"common/pkg/library"
	"common/pkg/listing"
	"common/pkg/user"
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
func (s *UserAppService) ListUserListings(userID string) ([]*listing.ListingView, error) {
	return s.listingFacade.ListByAuthor(userID)
}

// CreateListing creates a new listing authored by the user.
func (s *UserAppService) CreateListing(authorID, title, description string, accessMode listing.ListingAccessMode) (*listing.ListingView, error) {
	listingView, err := s.listingFacade.CreateListing(authorID, title, description, accessMode)
	if err != nil {
		return nil, err
	}
	return listingView, nil
}

// SaveListingToLibrary adds a listing to a user's library.
func (s *UserAppService) SaveListingToLibrary(userID string, listingID string) (*library.SavedListingView, error) {
	savedItem, err := s.libraryService.SaveListing(userID, listingID)
	if err != nil {
		return nil, err
	}
	return savedItem.ToView(), nil
}

// ListUserLibrary returns all saved listings for a user.
func (s *UserAppService) ListUserLibrary(userID string) ([]*library.SavedListingView, error) {
	savedItems, err := s.libraryService.ListUserLibrary(userID)
	if err != nil {
		return nil, err
	}

	views := make([]*library.SavedListingView, len(savedItems))
	for i, item := range savedItems {
		views[i] = item.ToView()
	}
	return views, nil
}

// RemoveListingFromLibrary removes a saved listing from a user's library.
func (s *UserAppService) RemoveListingFromLibrary(itemID string) error {
	return s.libraryService.RemoveListing(itemID)
}
