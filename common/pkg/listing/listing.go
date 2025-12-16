package listing

import (
	"common/pkg/hardware"
	"common/pkg/money"
	"common/pkg/util"
)

type PricingType int
type ListingAccessMode int

const (
	Private ListingAccessMode = iota
	Public
)

type Listing struct {
	ID                    string
	AuthorID              string
	Title                 string
	Description           string
	AccessMode            ListingAccessMode
	HardwareSpecification *hardware.HardwareSpecification
	Price                 money.Money
}

type ListingRepository interface {
	GetByID(listingID string) (*Listing, error)
	Save(listing *Listing) error
	Delete(listingID string) error
	ListByAuthor(authorID string) ([]*Listing, error)
}

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
		ID:          util.GenerateUUID(),
		AuthorID:    authorID,
		Title:       title,
		Description: description,
		AccessMode:  accessMode,
	}
	if err := s.repo.Save(listing); err != nil {
		return nil, err
	}
	return listing, nil
}

type ListingUpdate struct {
	Title                 *string
	Description           *string
	HardwareSpecification *hardware.HardwareSpecification
	Price                 *money.Money
	AccessMode            *ListingAccessMode
}

func (s *ListingService) UpdateListing(listingID string, update ListingUpdate) (*Listing, error) {
	listing, err := s.repo.GetByID(listingID)
	if err != nil {
		return nil, err
	}

	if update.Title != nil {
		listing.Title = *update.Title
	}
	if update.Description != nil {
		listing.Description = *update.Description
	}
	if update.HardwareSpecification != nil {
		listing.HardwareSpecification = update.HardwareSpecification
	}
	if update.Price != nil {
		listing.Price = *update.Price
	}
	if update.AccessMode != nil {
		listing.AccessMode = *update.AccessMode
	}

	if err := s.repo.Save(listing); err != nil {
		return nil, err
	}

	return listing, nil
}

// GetListing returns a read-only view of a listing including its pricings.
func (s *ListingService) GetListing(listingID string) (*Listing, error) {
	return s.repo.GetByID(listingID)
}

// ListByAuthor returns all listings authored by the given user ID as views.
func (s *ListingService) ListByAuthor(authorID string) ([]*Listing, error) {
	return s.repo.ListByAuthor(authorID)
}

// DeleteListing removes a listing by its ID.
func (s *ListingService) DeleteListing(listingID string) error {
	_, err := s.repo.GetByID(listingID)
	if err != nil {
		return err
	}
	return s.repo.Delete(listingID)
}
