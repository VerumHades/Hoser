package adapters

import (
	"context"

	"common/pkg/domain/entities/user"
	"common/pkg/domain/repositories"
	"common/pkg/shared"
)

// UserSavedListingViewAdapter provides a unified interface combining
// batch fetching of joined saved listings and existence checks.
type UserSavedListingViewAdapter struct {
	viewRepo  user.UserSavedListingViewRepository
	savedRepo repositories.SavedListingQueryRepository
}

// NewUserSavedListingViewAdapter constructs the adapter.
func NewUserSavedListingViewAdapter(
	viewRepo user.UserSavedListingViewRepository,
	savedRepo repositories.SavedListingQueryRepository,
) *UserSavedListingViewAdapter {
	return &UserSavedListingViewAdapter{
		viewRepo:  viewRepo,
		savedRepo: savedRepo,
	}
}

// FetchNextBatch fetches the next batch of saved listings joined with public listing info.
func (a *UserSavedListingViewAdapter) FetchNextBatchOfUserSavedListings(
	ctx context.Context,
	userID shared.UserID,
	request shared.BatchRequest[user.UserSavedListingViewCursor],
) ([]*user.UserSavedListingView, user.UserSavedListingViewCursor, error) {
	return a.viewRepo.FetchNextBatchOfUserSavedListings(ctx, userID, request)
}

// ExistsByUserAndListing checks if a user has saved a given listing.
func (a *UserSavedListingViewAdapter) ExistsByUserAndListing(
	ctx context.Context,
	userID shared.UserID,
	listingID shared.ListingID,
) (bool, error) {
	return a.savedRepo.ExistsByUserAndListing(ctx, userID, listingID)
}
