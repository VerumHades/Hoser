package listing

import (
	"errors"

	"common/internal/domain/money"
	"common/internal/shared"
)

// ListingAccessMode defines whether a listing is private or public.
type ListingAccessMode int

const (
	Private ListingAccessMode = iota
	Public
)

// PricingType represents different pricing strategies (expandable).
type PricingType int

// Listing represents a deployable or rentable listing.
type Listing struct {
	id                    shared.ListingID
	authorID              shared.UserID
	title                 string
	description           string
	accessMode            ListingAccessMode
	hardwareSpecification *shared.HardwareSpecification
	price                 money.Money
}

// NewListing creates a new listing with a generated ID.
func NewListing(authorID shared.UserID, title, description string, accessMode ListingAccessMode, hardwareSpec *shared.HardwareSpecification, price money.Money) (*Listing, error) {
	id := shared.ListingID(shared.GenerateUUID())
	return NewListingWithID(id, authorID, title, description, accessMode, hardwareSpec, price)
}

// NewListingWithID creates a listing with an existing ID (for repository hydration).
func NewListingWithID(id shared.ListingID, authorID shared.UserID, title, description string, accessMode ListingAccessMode, hardwareSpec *shared.HardwareSpecification, price money.Money) (*Listing, error) {
	listing := &Listing{
		id:                    id,
		authorID:              authorID,
		title:                 title,
		description:           description,
		accessMode:            accessMode,
		hardwareSpecification: hardwareSpec,
		price:                 price,
	}
	if err := listing.Validate(); err != nil {
		return nil, err
	}
	return listing, nil
}

// ID returns the listing's unique identifier.
func (l *Listing) ID() shared.ListingID {
	return l.id
}

// AuthorID returns the author's user ID.
func (l *Listing) AuthorID() shared.UserID {
	return l.authorID
}

// Title returns the listing title.
func (l *Listing) Title() string {
	return l.title
}

// Description returns the listing description.
func (l *Listing) Description() string {
	return l.description
}

// AccessMode returns whether the listing is public or private.
func (l *Listing) AccessMode() ListingAccessMode {
	return l.accessMode
}

// HardwareSpecification returns the listing's shared specification.
func (l *Listing) HardwareSpecification() *shared.HardwareSpecification {
	return l.hardwareSpecification
}

// Price returns the listing's price.
func (l *Listing) Price() money.Money {
	return l.price
}

// SetTitle updates the listing's title.
func (l *Listing) SetTitle(title string) error {
	if title == "" {
		return errors.New("title cannot be empty")
	}
	l.title = title
	return nil
}

// SetDescription updates the listing description.
func (l *Listing) SetDescription(description string) {
	l.description = description
}

// SetAccessMode updates the listing's access mode.
func (l *Listing) SetAccessMode(mode ListingAccessMode) {
	l.accessMode = mode
}

// SetPrice updates the listing's price.
func (l *Listing) SetPrice(price money.Money) {
	l.price = price
}

// Validate ensures the listing is in a valid state.
func (l *Listing) Validate() error {
	if l.id == "" {
		return errors.New("listing ID cannot be empty")
	}
	if l.authorID == "" {
		return errors.New("author ID cannot be empty")
	}
	if l.title == "" {
		return errors.New("title cannot be empty")
	}
	if l.hardwareSpecification == nil {
		return errors.New("shared specification cannot be nil")
	}
	if l.price.IsZero() {
		return errors.New("price cannot be zero")
	}
	return nil
}
