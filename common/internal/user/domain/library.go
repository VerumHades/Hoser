package domain

import (
	"common/internal/shared"
	"errors"
)

// SavedListingID represents the unique identifier for a saved listing.
type SavedListingID string

// UserID is reused from the User entity as a value object
// ListingID represents the referenced listing identifier.
type ListingID string

// SavedListing represents a listing saved by a user in their library.
type SavedListing struct {
	id        SavedListingID
	userID    UserID
	listingID ListingID
}

// NewSavedListing creates a new SavedListing with a generated unique ID.
func NewSavedListing(userID UserID, listingID ListingID) (*SavedListing, error) {
	id := SavedListingID(shared.GenerateUUID())
	return NewSavedListingWithID(id, userID, listingID)
}

// NewSavedListingWithID creates a SavedListing with an existing ID.
// Used when reconstructing from a repository.
func NewSavedListingWithID(id SavedListingID, userID UserID, listingID ListingID) (*SavedListing, error) {
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
func (s *SavedListing) ID() SavedListingID {
	return s.id
}

// UserID returns the owner of the library.
func (s *SavedListing) UserID() UserID {
	return s.userID
}

// ListingID returns the referenced listing ID.
func (s *SavedListing) ListingID() ListingID {
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
