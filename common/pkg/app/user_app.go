package app

import (
	"common/pkg/library"
	"common/pkg/listing"
	"common/pkg/user"
	"fmt"
)

// UserAppService handles cross-domain user operations.
type UserAppService struct {
	userService          *user.UserService
	publicListingService *PublicListingService
	libraryService       *library.LibraryService
}

// NewUserAppService creates a new UserAppService.
func NewUserAppService(
	userService *user.UserService,
	publicListingService *PublicListingService,
	libraryService *library.LibraryService,
) *UserAppService {
	return &UserAppService{
		userService:          userService,
		publicListingService: publicListingService,
		libraryService:       libraryService,
	}
}

// CreateListing creates a new listing authored by the user.
func (s *UserAppService) CreateListing(
	authorID string,
	title string,
	description string,
	accessMode listing.ListingAccessMode,
) (*listing.Listing, error) {
	if _, err := s.userService.GetUser(authorID); err != nil {
		return nil, fmt.Errorf("user %s does not exist", authorID)
	}

	return s.publicListingService.CreateListing(authorID, title, description, accessMode)
}

// SaveListingToLibrary saves a public listing to the user's library.
func (s *UserAppService) SaveListingToLibrary(userID string, listingID string) (*library.SavedListing, error) {
	if _, err := s.userService.GetUser(userID); err != nil {
		return nil, fmt.Errorf("user %s does not exist", userID)
	}

	publicListing, err := s.publicListingService.GetPublicListing(listingID)
	if err != nil {
		return nil, err
	}

	return s.libraryService.SaveListing(userID, publicListing.ID)
}

// ListUserLibrary returns all saved listings for a user.
func (s *UserAppService) ListUserLibrary(userID string) ([]*library.SavedListing, error) {
	return s.libraryService.ListUserLibrary(userID)
}

// RemoveListingFromLibrary removes a listing from a user's library.
func (s *UserAppService) RemoveListingFromLibrary(userID string, itemID string) error {
	userLibrary, err := s.libraryService.ListUserLibrary(userID)
	if err != nil {
		return fmt.Errorf("failed to fetch user's library: %w", err)
	}

	for _, item := range userLibrary {
		if item.ID == itemID {
			return s.libraryService.RemoveListing(itemID)
		}
	}

	return fmt.Errorf("library item %s does not belong to user %s", itemID, userID)
}

// IsUserDeveloper checks whether the user is a developer.
func (s *UserAppService) IsUserDeveloper(userID string) (bool, error) {
	return s.userService.CheckDeveloper(userID)
}

// GetUser returns a user.
func (s *UserAppService) GetUser(userID string) (*user.User, error) {
	return s.userService.GetUser(userID)
}
