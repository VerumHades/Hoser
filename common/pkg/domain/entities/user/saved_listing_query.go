package user

import (
	"common/pkg/shared"
	"context"
	"time"
)

// / DTO representing a saved listing joined with its public listing
type UserSavedListingView struct {
	SavedListingID shared.SavedListingID
	ListingID      shared.ListingID
	Title          *string
	Description    *string
	CreatedAt      time.Time
}

// / Cursor for this join query
type UserSavedListingViewCursor struct {
	LastCreatedAt time.Time
}

// / Read-only repository for paginated user saved listings with public listing info
type UserSavedListingViewRepository interface {
	FetchNextBatchOfUserSavedListings(
		ctx context.Context,
		userID shared.UserID,
		request shared.BatchRequest[UserSavedListingViewCursor],
	) (items []*UserSavedListingView, nextCursor UserSavedListingViewCursor, err error)
}
