package listing

import (
	"errors"
	"time"

	"common/pkg/domain/entities/events"
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

	screenshotKeys        []shared.ListingScreenshotID
	documentationMarkdown string
}

// --- Listing constructors with creation event ---

func NewListing(authorID shared.UserID, title, description string, accessMode ListingAccessMode, hardwareSpec *shared.HardwareSpecification, price int64, documentationMarkdown string) (*Listing, events.DomainEvent, error) {
	id := shared.ListingID(shared.GenerateUUID())
	return NewListingWithID(id, authorID, title, description, accessMode, hardwareSpec, price, time.Now(), nil, documentationMarkdown)
}

func NewListingWithID(
	id shared.ListingID,
	authorID shared.UserID,
	title, description string,
	accessMode ListingAccessMode,
	hardwareSpec *shared.HardwareSpecification,
	price int64,
	createdAt time.Time,
	screenshotKeys []shared.ListingScreenshotID,
	documentationMarkdown string,
) (*Listing, events.DomainEventEnvelope[ListingCreatedEvent], error) {

	listing := &Listing{
		id:                    id,
		authorID:              authorID,
		title:                 title,
		description:           description,
		accessMode:            accessMode,
		hardwareSpecification: hardwareSpec,
		priceInMinorUnits:     price,
		createdAt:             createdAt,
		screenshotKeys:        screenshotKeys,
		documentationMarkdown: documentationMarkdown,
	}

	if err := listing.Validate(); err != nil {
		return nil, events.DomainEventEnvelope[ListingCreatedEvent]{}, err
	}

	// Build typed payload instead of map[string]any
	payload := ListingCreatedEvent{
		ID:                    listing.id,
		Title:                 listing.title,
		Description:           listing.description,
		AccessMode:            listing.accessMode,
		HardwareSpecification: *listing.hardwareSpecification,
		PriceInMinorUnits:     listing.priceInMinorUnits,
		ScreenshotKeys:        screenshotKeys,
		DocumentationMarkdown: documentationMarkdown,
	}

	// Wrap in typed domain event envelope
	eventEnvelope := events.NewDomainEventEnvelope(payload)
	return listing, eventEnvelope, nil
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
func (l *Listing) PriceInMinorUnits() int64                     { return l.priceInMinorUnits }
func (l *Listing) CreatedAt() time.Time                         { return l.createdAt }
func (l *Listing) ScreenshotKeys() []shared.ListingScreenshotID { return l.screenshotKeys }
func (l *Listing) DocumentationMarkdown() string                { return l.documentationMarkdown }

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
	if l.priceInMinorUnits < 100 {
		return errors.New("price must be at least one whole unit (>= 100)")
	}
	if l.createdAt.IsZero() {
		return errors.New("createdAt cannot be zero")
	}
	return nil
}
