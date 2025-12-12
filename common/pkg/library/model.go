package library

// SavedListing represents a listing saved by a user in their library.
type SavedListing struct {
	ID        string // unique ID for this saved item
	UserID    string // the owner of the library
	ListingID string // the referenced listing
}
