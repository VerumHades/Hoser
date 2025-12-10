package purchase

import "common/pkg/currency"

// PurchaseView is a read-only representation of a purchase.
type PurchaseView struct {
	ID        string
	UserID    string
	ListingID string
	PricePaid currency.Money
}

// ToView converts a Purchase model to a read-only view.
func (p *Purchase) ToView() *PurchaseView {
	return &PurchaseView{
		ID:        p.id,
		UserID:    p.userID,
		ListingID: p.listingID,
		PricePaid: p.pricePaid,
	}
}
