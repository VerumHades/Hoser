package listing

import (
	"errors"
	"fmt"
	"time"

	"common/pkg/shared"
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
	priceInMinorUnits     int64
	createdAt             time.Time
}

// NewListing creates a new listing with a generated ID and current timestamp.
func NewListing(authorID shared.UserID, title, description string, accessMode ListingAccessMode, hardwareSpec *shared.HardwareSpecification, price int64) (*Listing, error) {
	id := shared.ListingID(shared.GenerateUUID())
	return NewListingWithID(id, authorID, title, description, accessMode, hardwareSpec, price, time.Now())
}

// NewListingWithID creates a listing with an existing ID and specified creation time.
func NewListingWithID(id shared.ListingID, authorID shared.UserID, title, description string, accessMode ListingAccessMode, hardwareSpec *shared.HardwareSpecification, price int64, createdAt time.Time) (*Listing, error) {
	listing := &Listing{
		id:                    id,
		authorID:              authorID,
		title:                 title,
		description:           description,
		accessMode:            accessMode,
		hardwareSpecification: hardwareSpec,
		priceInMinorUnits:     price,
		createdAt:             createdAt,
	}
	if err := listing.Validate(); err != nil {
		return nil, err
	}
	return listing, nil
}

// CreatedAt returns the listing creation timestamp.
func (l *Listing) CreatedAt() time.Time {
	return l.createdAt
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
func (l *Listing) PriceInMinorUnits() int64 {
	return l.priceInMinorUnits
}

func (l *Listing) SetTitle(title string) error {
	if title == "" {
		return fmt.Errorf("title cannot be empty")
	}
	l.title = title
	return nil
}

func (l *Listing) SetDescription(description string) error {
	l.description = description
	return nil
}

func (l *Listing) SetHardware(hardware *shared.HardwareSpecification) error {
	l.hardwareSpecification = hardware
	return nil
}

func (l *Listing) SetPriceInMinorUnits(price int64) error {
	if price < 0 {
		return fmt.Errorf("price cannot be negative")
	}
	l.priceInMinorUnits = price
	return nil
}

func (l *Listing) SetAccessMode(mode ListingAccessMode) error {
	l.accessMode = mode
	return nil
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
		return errors.New("hardware specification cannot be nil")
	}
	if l.priceInMinorUnits < 0 {
		return errors.New("price cannot be negative")
	}
	if l.createdAt.IsZero() {
		return errors.New("createdAt cannot be zero")
	}
	return nil
}
