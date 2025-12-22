package library

import "common/pkg/util"

// SavedListing represents a listing saved by a user in their library.
type SavedListing struct {
	ID        string // unique ID for this saved item
	UserID    string // the owner of the library
	ListingID string // the referenced listing
}

// SavedListingRepository defines persistence operations for user library items.
type SavedListingRepository interface {
	Save(item *SavedListing) error
	GetByUserID(userID string) ([]*SavedListing, error)
	GetByID(itemID string) (*SavedListing, error)
	Delete(itemID string) error
	ExistsByUserAndListing(userID string, listingID string) (bool, error)
}

// LibraryService handles application-level operations for user libraries.
type LibraryService struct {
	repository SavedListingRepository
}

// NewLibraryService creates a new instance of LibraryService.
func NewLibraryService(repository SavedListingRepository) *LibraryService {
	return &LibraryService{
		repository: repository,
	}
}

// SaveListing stores a listing in the user's library.
func (service *LibraryService) SaveListing(userID string, listingID string) (*SavedListing, error) {
	savedListing := &SavedListing{
		ID:        util.GenerateUUID(),
		UserID:    userID,
		ListingID: listingID,
	}

	if err := service.repository.Save(savedListing); err != nil {
		return nil, err
	}

	return savedListing, nil
}

// HasUserListingInLibrary checks whether a user already saved a given listing.
func (service *LibraryService) HasUserListingInLibrary(userID string, listingID string) (bool, error) {
	return service.repository.ExistsByUserAndListing(userID, listingID)
}

// ListUserLibrary returns all saved listings for a user.
func (service *LibraryService) ListUserLibrary(userID string) ([]*SavedListing, error) {
	return service.repository.GetByUserID(userID)
}

// RemoveListing removes a saved listing from the library.
func (service *LibraryService) RemoveListing(itemID string) error {
	return service.repository.Delete(itemID)
}
