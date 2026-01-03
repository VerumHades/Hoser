package listing

import (
	"common/pkg/domain/entities/events"
	"common/pkg/shared"
)

// ListingMutationBuilder allows building and applying multiple changes to a Listing
type ListingMutationBuilder struct {
	listing    *Listing
	payload    ListingUpdateEvent
	hasChanges bool
}

// Mutate starts a builder for modifying the Listing
func (l *Listing) Mutate() *ListingMutationBuilder {
	return &ListingMutationBuilder{
		listing: l,
	}
}

// SetTitle schedules a title change
func (b *ListingMutationBuilder) SetTitle(title string) *ListingMutationBuilder {
	if title != "" && title != b.listing.title {
		b.payload.Title = &title
		b.hasChanges = true
	}
	return b
}

// SetDescription schedules a description change
func (b *ListingMutationBuilder) SetDescription(description string) *ListingMutationBuilder {
	if description != b.listing.description {
		b.payload.Description = &description
		b.hasChanges = true
	}
	return b
}

// SetHardware schedules a hardware specification change
func (b *ListingMutationBuilder) SetHardware(hardware *shared.HardwareSpecification) *ListingMutationBuilder {
	if hardware != b.listing.hardwareSpecification {
		b.payload.HardwareSpecification = hardware
		b.hasChanges = true
	}
	return b
}

// SetPriceInMinorUnits schedules a price change
func (b *ListingMutationBuilder) SetPriceInMinorUnits(price int64) *ListingMutationBuilder {
	if price >= 0 && price != b.listing.priceInMinorUnits {
		b.payload.PriceInMinorUnits = &price
		b.hasChanges = true
	}
	return b
}

// SetAccessMode schedules an access mode change
func (b *ListingMutationBuilder) SetAccessMode(mode ListingAccessMode) *ListingMutationBuilder {
	if mode != b.listing.accessMode {
		b.payload.AccessMode = &mode
		b.hasChanges = true
	}
	return b
}

// Apply applies all changes and returns a single DomainEventEnvelope
func (b *ListingMutationBuilder) Apply() (events.DomainEventEnvelope[ListingUpdateEvent], error) {
	if !b.hasChanges {
		return events.DomainEventEnvelope[ListingUpdateEvent]{}, nil
	}

	if b.payload.Title != nil {
		b.listing.title = *b.payload.Title
	}
	if b.payload.Description != nil {
		b.listing.description = *b.payload.Description
	}
	if b.payload.AccessMode != nil {
		b.listing.accessMode = *b.payload.AccessMode
	}
	if b.payload.PriceInMinorUnits != nil {
		b.listing.priceInMinorUnits = *b.payload.PriceInMinorUnits
	}
	if b.payload.HardwareSpecification != nil {
		b.listing.hardwareSpecification = b.payload.HardwareSpecification
	}

	b.payload.ID = b.listing.id

	envelope := events.NewDomainEventEnvelope(b.payload)
	return envelope, nil
}
