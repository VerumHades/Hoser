package listing

import (
	"common/pkg/billing/currency"
	"common/pkg/util"
)

// ListingService provides higher-level business logic around Listings.
type ListingService struct {
	repo ListingRepository
}

// NewListingService creates a new ListingService instance.
func NewListingService(repo ListingRepository) *ListingService {
	return &ListingService{repo: repo}
}

// CreateListing creates and stores a new listing.
func (s *ListingService) CreateListing(authorID, title, description string, accessMode ListingAccessMode) (*Listing, error) {
	listing := &Listing{
		id:          util.GenerateUUID(), // assume you have a helper for UUIDs
		authorID:    authorID,
		title:       title,
		description: description,
		accessMode:  accessMode,
		pricing:     []*Pricing{},
	}
	if err := s.repo.Save(listing); err != nil {
		return nil, err
	}
	return listing, nil
}

// UpdateTitle changes a listing's title and persists it.
func (s *ListingService) UpdateTitle(listingID, newTitle string) error {
	listing, err := s.repo.GetByID(listingID)
	if err != nil {
		return err
	}
	listing.title = newTitle
	return s.repo.Save(listing)
}

// AddPricing adds a pricing entry to a listing.
func (s *ListingService) AddPricing(listingID string, typ PricingType, amount currency.Money) (*Pricing, error) {
	listing, err := s.repo.GetByID(listingID)
	if err != nil {
		return nil, err
	}
	pricing := &Pricing{
		id:     util.GenerateUUID(),
		typ:    typ,
		amount: amount,
	}
	listing.pricing = append(listing.pricing, pricing)
	if err := s.repo.Save(listing); err != nil {
		return nil, err
	}
	return pricing, nil
}

// GetListingView returns a read-only view of a listing.
func (s *ListingService) GetListingView(listingID string) (*ListingView, error) {
	listing, err := s.repo.GetByID(listingID)
	if err != nil {
		return nil, err
	}

	pricingViews := make([]*PricingView, len(listing.pricing))
	for i, p := range listing.pricing {
		pricingViews[i] = &PricingView{
			ID:     p.id,
			Type:   p.typ,
			Amount: p.amount,
		}
	}

	return listing.ToView(), nil
}

// ListByAuthor returns all listings authored by the given user ID as views.
func (s *ListingService) ListByAuthor(authorID string) ([]*ListingView, error) {
	listings, err := s.repo.ListByAuthor(authorID)
	if err != nil {
		return nil, err
	}

	views := make([]*ListingView, len(listings))
	for i, l := range listings {
		views[i] = l.ToView()
	}
	return views, nil
}

// DeleteListing removes a listing by its ID.
func (s *ListingService) DeleteListing(listingID string) error {
	// Retrieve the listing first to ensure it exists
	_, err := s.repo.GetByID(listingID)
	if err != nil {
		return err
	}

	// Delete the listing
	return s.repo.Delete(listingID)
}
