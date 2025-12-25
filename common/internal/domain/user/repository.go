package user

import "common/internal/shared"

// UserRepository defines persistence operations for users.
type UserRepository interface {
	Save(user *User) error
	GetByID(id shared.UserID) (*User, error)
	GetByUsername(username string) (*User, error)
	Delete(id shared.UserID) error
}

// SavedListingRepository defines persistence operations for user library items.
type SavedListingRepository interface {
	Save(item *SavedListing) error
	GetByUserID(userID shared.UserID) ([]*SavedListing, error)
	GetByID(itemID shared.SavedListingID) (*SavedListing, error)
	Delete(itemID shared.SavedListingID) error
	ExistsByUserAndListing(userID shared.UserID, listingID shared.ListingID) (bool, error)
}
