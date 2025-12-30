package user

import (
	"common/pkg/domain/user"
	"common/pkg/shared"
	"context"
)

type UserAPIConfiguration struct {
	JWTSecret string `env:"JWT_SECRET" default:"SECRET"`
}

type userAuthentificationService interface {
	AuthenticateUser(context context.Context, username string, password string) (*user.User, error)
}

type userQueryService interface {
	GetByID(context context.Context, userID shared.UserID) (*user.User, error)
}

type userLibraryQueryService interface {
	// ExistsByUserAndListing checks if a user has saved a specific listing.
	ExistsByUserAndListing(
		ctx context.Context,
		userID shared.UserID,
		listingID shared.ListingID,
	) (bool, error)

	FetchNextBatchOfUserSavedListings(
		ctx context.Context,
		userID shared.UserID,
		request shared.BatchRequest[user.UserSavedListingViewCursor],
	) (items []*user.UserSavedListingView, nextCursor user.UserSavedListingViewCursor, err error)
}

type UserAPI struct {
	runningConfiguration UserAPIConfiguration

	userAuthentificationService userAuthentificationService
	userLibraryQueryService     userLibraryQueryService
	userQueryService            userQueryService
}
