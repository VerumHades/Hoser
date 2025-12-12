package library

// SavedListingView exposes only the IDs for a saved listing.
type SavedListingView struct {
	ID        string
	UserID    string
	ListingID string
}

// ToView converts a SavedListing model into a SavedListingView.
func (item *SavedListing) ToView() *SavedListingView {
	return &SavedListingView{
		ID:        item.ID,
		UserID:    item.UserID,
		ListingID: item.ListingID,
	}
}
