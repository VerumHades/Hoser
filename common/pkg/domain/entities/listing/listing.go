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

// --- Listing constructors with creation event ---

func NewListing(authorID shared.UserID, title, description string, accessMode ListingAccessMode, hardwareSpec *shared.HardwareSpecification, price int64) (*Listing, DomainEvent, error) {
	id := shared.ListingID(shared.GenerateUUID())
	return NewListingWithID(id, authorID, title, description, accessMode, hardwareSpec, price, time.Now())
}

func NewListingWithID(id shared.ListingID, authorID shared.UserID, title, description string, accessMode ListingAccessMode, hardwareSpec *shared.HardwareSpecification, price int64, createdAt time.Time) (*Listing, DomainEvent, error) {
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
		return nil, nil, err
	}

	payload := map[string]any{
		"title":                 listing.title,
		"description":           listing.description,
		"accessMode":            listing.accessMode,
		"hardwareSpecification": listing.hardwareSpecification,
		"priceInMinorUnits":     listing.priceInMinorUnits,
	}

	event := NewListingCreated(listing.id, payload, createdAt)
	return listing, event, nil
}

// --- Listing getters ---

func (l *Listing) ID() shared.ListingID          { return l.id }
func (l *Listing) AuthorID() shared.UserID       { return l.authorID }
func (l *Listing) Title() string                 { return l.title }
func (l *Listing) Description() string           { return l.description }
func (l *Listing) AccessMode() ListingAccessMode { return l.accessMode }
func (l *Listing) HardwareSpecification() *shared.HardwareSpecification {
	return l.hardwareSpecification
}
func (l *Listing) PriceInMinorUnits() int64 { return l.priceInMinorUnits }
func (l *Listing) CreatedAt() time.Time     { return l.createdAt }

// --- Mutating setters with update events ---

func (l *Listing) SetTitle(title string) (DomainEvent, error) {
	if title == "" {
		return nil, fmt.Errorf("title cannot be empty")
	}
	if title == l.title {
		return nil, nil
	}
	l.title = title
	payload := map[string]any{"title": title}
	return NewListingUpdated(l.id, payload, time.Now()), nil
}

func (l *Listing) SetDescription(description string) (DomainEvent, error) {
	if description == l.description {
		return nil, nil
	}
	l.description = description
	payload := map[string]any{"description": description}
	l.description = description
	return NewListingUpdated(l.id, payload, time.Now()), nil
}

func (l *Listing) SetHardware(hardware *shared.HardwareSpecification) (DomainEvent, error) {
	if hardware == l.hardwareSpecification {
		return nil, nil
	}
	l.hardwareSpecification = hardware
	payload := map[string]any{"hardwareSpecification": hardware}
	return NewListingUpdated(l.id, payload, time.Now()), nil
}

func (l *Listing) SetPriceInMinorUnits(price int64) (DomainEvent, error) {
	if price < 0 {
		return nil, fmt.Errorf("price cannot be negative")
	}
	if price == l.priceInMinorUnits {
		return nil, nil
	}
	l.priceInMinorUnits = price
	payload := map[string]any{"priceInMinorUnits": price}
	return NewListingUpdated(l.id, payload, time.Now()), nil
}

func (l *Listing) SetAccessMode(mode ListingAccessMode) (DomainEvent, error) {
	if mode == l.accessMode {
		return nil, nil
	}
	l.accessMode = mode
	payload := map[string]any{"accessMode": mode}
	return NewListingUpdated(l.id, payload, time.Now()), nil
}

// --- Validation ---

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
