package domain

// UserRepository defines persistence operations for users.
type UserRepository interface {
	Save(user *User) error
	GetByID(id string) (*User, error)
	GetByUsername(username string) (*User, error)
	Delete(id string) error
}

// SavedListingRepository defines persistence operations for user library items.
type SavedListingRepository interface {
	Save(item *SavedListing) error
	GetByUserID(userID string) ([]*SavedListing, error)
	GetByID(itemID string) (*SavedListing, error)
	Delete(itemID string) error
	ExistsByUserAndListing(userID string, listingID string) (bool, error)
}
