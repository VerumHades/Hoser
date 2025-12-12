package library

// SavedListingRepository defines persistence operations for user library items.
type SavedListingRepository interface {
	Save(item *SavedListing) error
	GetByUserID(userID string) ([]*SavedListing, error)
	Delete(itemID string) error
	GetByID(itemID string) (*SavedListing, error)
}
