package listing

type ListingRepository interface {
	GetByID(listingID string) (*Listing, error)
	Save(listing *Listing) error
	Delete(listingID string) error
	ListByAuthor(authorID string) ([]*Listing, error)

	AddPricing(listingID string, pricing *Pricing) error
	RemovePricing(listingID, pricingID string) error
	ListPricing(listingID string) ([]*Pricing, error)
}
