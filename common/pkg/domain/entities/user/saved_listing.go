package user

import (
	"common/pkg/shared"
	"errors"
	"time"
)

// SavedListing represents a listing saved by a user in their library.
type SavedListing struct {
	id        shared.SavedListingID
	userID    shared.UserID
	listingID shared.ListingID
	createdAt time.Time
}

// NewSavedListing creates a new SavedListing with a generated unique ID and current timestamp.
func NewSavedListing(userID shared.UserID, listingID shared.ListingID) (*SavedListing, error) {
	id := shared.SavedListingID(shared.GenerateUUID())
	return NewSavedListingWithID(id, userID, listingID, time.Now())
}

// NewSavedListingWithID creates a SavedListing with an existing ID and creation time (for repository hydration).
func NewSavedListingWithID(id shared.SavedListingID, userID shared.UserID, listingID shared.ListingID, createdAt time.Time) (*SavedListing, error) {
	saved := &SavedListing{
		id:        id,
		userID:    userID,
		listingID: listingID,
		createdAt: createdAt,
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

// UserID returns the owner of the library.
func (s *SavedListing) UserID() shared.UserID {
	return s.userID
}

// ListingID returns the referenced listing ID.
func (s *SavedListing) ListingID() shared.ListingID {
	return s.listingID
}

// CreatedAt returns the timestamp when this saved listing was created.
func (s *SavedListing) CreatedAt() time.Time {
	return s.createdAt
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
	if s.createdAt.IsZero() {
		return errors.New("createdAt cannot be zero")
	}
	return nil
}
