package library

import "common/pkg/util"

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

// SaveListing saves a listing to a user's library.
func (service *LibraryService) SaveListing(userID string, listingID string) (*SavedListingView, error) {
	savedItem := &SavedListing{
		ID:        util.GenerateUUID(), // assume a UUID generator function
		UserID:    userID,
		ListingID: listingID,
	}
	if err := service.repository.Save(savedItem); err != nil {
		return nil, err
	}
	return savedItem.ToView(), nil
}

// ListUserLibrary returns all listings saved by a user.
func (service *LibraryService) ListUserLibrary(userID string) ([]*SavedListing, error) {
	return service.repository.GetByUserID(userID)
}

// RemoveListing removes a saved listing from a user's library.
func (service *LibraryService) RemoveListing(itemID string) error {
	return service.repository.Delete(itemID)
}
