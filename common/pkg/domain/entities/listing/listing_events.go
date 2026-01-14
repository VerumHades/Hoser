package listing

import (
	"common/pkg/shared"
)

// ListingCreated is emitted when a new listing is created.
type ListingCreatedEvent struct {
	ID                    shared.ListingID
	Title                 string
	Description           string
	AccessMode            ListingAccessMode
	HardwareSpecification shared.HardwareSpecification
	PriceInMinorUnits     int64
	ScreenshotKeys        []shared.ListingScreenshotID
}

type ListingUpdateEvent struct {
	ID                    shared.ListingID
	Title                 *string
	Description           *string
	AccessMode            *ListingAccessMode
	HardwareSpecification *shared.HardwareSpecification
	PriceInMinorUnits     *int64
	ScreenshotKeys        []shared.ListingScreenshotID
}

type ListingDeleteEvent struct {
	ID shared.ListingID
}
