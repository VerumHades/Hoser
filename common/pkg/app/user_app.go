package app

import (
	"common/pkg/listing"
	"common/pkg/user"
)

// UserAppService handles cross-domain operations involving users.
type UserAppService struct {
	userService *user.UserService
	listingRepo listing.ListingRepository
}

func NewUserAppService(userService *user.UserService, listingRepo listing.ListingRepository) *UserAppService {
	return &UserAppService{
		userService: userService,
		listingRepo: listingRepo,
	}
}

func (s *UserAppService) ListUserListings(userID string) ([]*listing.ListingView, error) {
	listings, err := s.listingRepo.ListByAuthor(userID)
	if err != nil {
		return nil, err
	}

	views := make([]*listing.ListingView, len(listings))
	for i, l := range listings {
		views[i] = l.ToView()
	}
	return views, nil
}
