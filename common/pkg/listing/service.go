package listing

import (
	"common/pkg/hardware"
	"common/pkg/money"
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
func (s *ListingService) CreateListing(authorID, title, description string, accessMode ListingAccessMode) (*ListingView, error) {
	listing := &Listing{
		id:          util.GenerateUUID(),
		authorID:    authorID,
		title:       title,
		description: description,
		accessMode:  accessMode,
	}
	if err := s.repo.Save(listing); err != nil {
		return nil, err
	}
	return listing.ToView([]*Pricing{}), nil
}

type ListingUpdate struct {
	Title                 *string
	Description           *string
	HardwareSpecification *hardware.HardwareSpecification
}

func (s *ListingService) UpdateListing(listingID string, update ListingUpdate) (*ListingView, error) {
	listing, err := s.repo.GetByID(listingID)
	if err != nil {
		return nil, err
	}

	if update.Title != nil {
		listing.title = *update.Title
	}
	if update.Description != nil {
		listing.description = *update.Description
	}
	if update.HardwareSpecification != nil {
		listing.hardwareSpecification = update.HardwareSpecification
	}

	if err := s.repo.Save(listing); err != nil {
		return nil, err
	}

	// Fetch pricings for view
	pricings, err := s.repo.ListPricing(listingID)
	if err != nil {
		return nil, err
	}

	return listing.ToView(pricings), nil
}

// AddPricing adds a pricing entry to a listing via the repository.
func (s *ListingService) AddPricing(listingID string, typ PricingType, amount money.Money) (*Pricing, error) {
	pricing := &Pricing{
		id:     util.GenerateUUID(),
		typ:    typ,
		amount: amount,
	}
	if err := s.repo.AddPricing(listingID, pricing); err != nil {
		return nil, err
	}
	return pricing, nil
}

// RemovePricing removes a pricing entry from a listing via the repository.
func (s *ListingService) RemovePricing(listingID, pricingID string) error {
	return s.repo.RemovePricing(listingID, pricingID)
}

// GetListingView returns a read-only view of a listing including its pricings.
func (s *ListingService) GetListingView(listingID string) (*ListingView, error) {
	listing, err := s.repo.GetByID(listingID)
	if err != nil {
		return nil, err
	}

	pricings, err := s.repo.ListPricing(listingID)
	if err != nil {
		return nil, err
	}

	pricingViews := make([]*PricingView, len(pricings))
	for i, p := range pricings {
		pricingViews[i] = &PricingView{
			ID:     p.id,
			Type:   p.typ,
			Amount: p.amount,
		}
	}

	var hwView *hardware.HardwareSpecificationView
	if listing.hardwareSpecification != nil {
		hwView = listing.hardwareSpecification.ToView()
	}

	return &ListingView{
		ID:                    listing.id,
		AuthorID:              listing.authorID,
		Title:                 listing.title,
		Description:           listing.description,
		AccessMode:            listing.accessMode,
		Pricing:               pricingViews,
		HardwareSpecification: hwView,
	}, nil
}

// ListByAuthor returns all listings authored by the given user ID as views.
func (s *ListingService) ListByAuthor(authorID string) ([]*ListingView, error) {
	listings, err := s.repo.ListByAuthor(authorID)
	if err != nil {
		return nil, err
	}

	views := make([]*ListingView, len(listings))
	for i, l := range listings {
		pricings, _ := s.repo.ListPricing(l.id)
		views[i] = l.ToView(pricings)
	}

	return views, nil
}

// DeleteListing removes a listing by its ID.
func (s *ListingService) DeleteListing(listingID string) error {
	_, err := s.repo.GetByID(listingID)
	if err != nil {
		return err
	}
	return s.repo.Delete(listingID)
}
