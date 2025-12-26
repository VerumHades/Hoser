package user

import (
	"common/internal/shared"
	"errors"
)

// SavedListing represents a listing saved by a user in their library.
type SavedListing struct {
	id        shared.SavedListingID
	userID    shared.UserID
	listingID shared.ListingID
}

// NewSavedListing creates a new SavedListing with a generated unique ID.
func NewSavedListing(userID shared.UserID, listingID shared.ListingID) (*SavedListing, error) {
	id := shared.SavedListingID(shared.GenerateUUID())
	return NewSavedListingWithID(id, userID, listingID)
}

// NewSavedListingWithID creates a SavedListing with an existing ID.
// Used when reconstructing from a repository.
func NewSavedListingWithID(id shared.SavedListingID, userID shared.UserID, listingID shared.ListingID) (*SavedListing, error) {
	saved := &SavedListing{
		id:        id,
		userID:    userID,
		listingID: listingID,
	}
	if err := saved.Validate(); err != nil {
		return nil, err
	}
	return saved, nil
}

// ID returns the unique identifier of this saved listing.
func (s *SavedListing) ID() shared.SavedListingID {
	return s.id
}

// shared.UserID returns the owner of the library.
func (s *SavedListing) UserID() shared.UserID {
	return s.userID
}

// shared.ListingID returns the referenced listing ID.
func (s *SavedListing) ListingID() shared.ListingID {
	return s.listingID
}

// Validate ensures the saved listing is in a valid state.
func (s *SavedListing) Validate() error {
	if s.id == "" {
		return errors.New("saved listing ID cannot be empty")
	}
	if s.userID == "" {
		return errors.New("user ID cannot be empty")
	}
	if s.listingID == "" {
		return errors.New("listing ID cannot be empty")
	}
	return nil
}
