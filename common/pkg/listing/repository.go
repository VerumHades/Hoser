package listing

// ListingRepository defines persistence operations for Listings.
type ListingRepository interface {
	// Save creates or updates a listing.
	Save(listing *Listing) error

	// GetByID retrieves a listing by its UUID.
	GetByID(id string) (*Listing, error)

	// Delete removes a listing by its UUID.
	Delete(id string) error

	// ListByAuthor returns all listings by a given author ID.
	ListByAuthor(authorID string) ([]*Listing, error)
}
