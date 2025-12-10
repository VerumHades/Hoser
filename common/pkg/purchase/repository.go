package purchase

// PurchaseRepository defines persistence operations for purchases.
type PurchaseRepository interface {
	Save(purchase *Purchase) error
	GetByID(id string) (*Purchase, error)
	ListByUser(userID string) ([]*Purchase, error)
	ListByListing(listingID string) ([]*Purchase, error)
	ExistsByUserAndListing(userID, listingID string) (bool, error)
	Delete(id string) error
}
