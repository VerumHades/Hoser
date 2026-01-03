package listing

import (
	"encoding/json"
	"time"

	"common/pkg/domain/entities/events"
	"common/pkg/shared"
)

// DomainEvent interface for all listing events
type DomainEvent interface {
	EventType() events.EventType
	AggregateID() shared.ListingID
	OccurredAt() time.Time
	Payload() []byte
}

// ListingCreated is emitted when a new listing is created.
type ListingCreated struct {
	listingID  shared.ListingID
	payload    []byte
	occurredAt time.Time
}

func NewListingCreated(listingID shared.ListingID, payload interface{}, occurredAt time.Time) ListingCreated {
	data, _ := json.Marshal(payload) // payload can be partial state
	return ListingCreated{
		listingID:  listingID,
		payload:    data,
		occurredAt: occurredAt.UTC(),
	}
}

func (e ListingCreated) EventType() events.EventType   { return events.EventTypeListingCreated }
func (e ListingCreated) AggregateID() shared.ListingID { return e.listingID }
func (e ListingCreated) OccurredAt() time.Time         { return e.occurredAt }
func (e ListingCreated) Payload() []byte               { return e.payload }

// ListingUpdated is emitted when a listing is updated.
type ListingUpdated struct {
	listingID  shared.ListingID
	payload    []byte
	occurredAt time.Time
}

func NewListingUpdated(listingID shared.ListingID, payload interface{}, occurredAt time.Time) ListingUpdated {
	data, _ := json.Marshal(payload) // payload can include only changed fields
	return ListingUpdated{
		listingID:  listingID,
		payload:    data,
		occurredAt: occurredAt.UTC(),
	}
}

func (e ListingUpdated) EventType() events.EventType   { return events.EventTypeListingUpdated }
func (e ListingUpdated) AggregateID() shared.ListingID { return e.listingID }
func (e ListingUpdated) OccurredAt() time.Time         { return e.occurredAt }
func (e ListingUpdated) Payload() []byte               { return e.payload }
