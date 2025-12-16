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
	Delete(itemID string) error
	GetByID(itemID string) (*SavedListing, error)
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

func (service *LibraryService) SaveListing(userID string, listingID string) (*SavedListing, error) {
	savedItem := &SavedListing{
		ID:        util.GenerateUUID(), // assume a UUID generator function
		UserID:    userID,
		ListingID: listingID,
	}
	if err := service.repository.Save(savedItem); err != nil {
		return nil, err
	}
	return savedItem, nil
}

func (service *LibraryService) ListUserLibrary(userID string) ([]*SavedListing, error) {
	return service.repository.GetByUserID(userID)
}

func (service *LibraryService) RemoveListing(itemID string) error {
	return service.repository.Delete(itemID)
}
