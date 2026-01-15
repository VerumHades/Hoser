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
	DocumentationMarkdown string
}

type ListingUpdateEvent struct {
	ID                    shared.ListingID
	Title                 *string
	Description           *string
	AccessMode            *ListingAccessMode
	HardwareSpecification *shared.HardwareSpecification
	PriceInMinorUnits     *int64
	ScreenshotKeys        *[]shared.ListingScreenshotID
	DocumentationMarkdown *string
}

type ListingDeleteEvent struct {
	ID shared.ListingID
}
