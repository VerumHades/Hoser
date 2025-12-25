package listing

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

func (s *UserAppService) HasUserBoughtListing(userID string, listingID string) (bool, error) {
	return s.HasUserBoughtListing(userID, listingID)
}

// IsUserDeveloper checks whether the user is a developer.
func (s *UserAppService) IsUserDeveloper(userID string) (bool, error) {
	return s.userService.CheckDeveloper(userID)
}

// GetUser returns a user.
func (s *UserAppService) GetUser(userID string) (*user.User, error) {
	return s.userService.GetUser(userID)
}
