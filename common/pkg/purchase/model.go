package purchase

import "common/pkg/currency"

// Purchase represents a purchase made by a user for a listing.
type Purchase struct {
	id        string
	userID    string
	listingID string
	pricePaid currency.Money
}

// ID returns the unique identifier of the purchase.
func (p *Purchase) ID() string {
	return p.id
}

// UserID returns the ID of the purchasing user.
func (p *Purchase) UserID() string {
	return p.userID
}

// ListingID returns the ID of the purchased listing.
func (p *Purchase) ListingID() string {
	return p.listingID
}

// PricePaid returns the amount paid for this purchase.
func (p *Purchase) PricePaid() currency.Money {
	return p.pricePaid
}

// SetPricePaid updates the purchase price.
func (p *Purchase) SetPricePaid(price currency.Money) {
	p.pricePaid = price
}
