package events

import (
	"time"
)

// DomainEvent is the interface all listing events implement
type DomainEvent interface {
	EventType() EventType
	AggregateID() string
	OccurredAt() time.Time
}

// EventType represents the type of domain event.
type EventType string

const (
	EventTypeListingCreated EventType = "listing_created"
	EventTypeListingUpdated EventType = "listing_updated"
)
