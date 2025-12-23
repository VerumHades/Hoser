package app

import (
	"common/pkg/listing"
	"fmt"
)

// PublicListingService handles listing access rules and ownership checks.
type PublicListingService struct {
	listingFacade *ListingFacadeService
}

// NewPublicListingService creates a new PublicListingService.
func NewPublicListingService(listingFacade *ListingFacadeService) *PublicListingService {
	return &PublicListingService{
		listingFacade: listingFacade,
	}
}

// ListByAuthor returns all listings authored by a user.
func (s *PublicListingService) ListByAuthor(authorID string) ([]*listing.Listing, error) {
	return s.listingFacade.ListByAuthor(authorID)
}

// SearchPublicListings returns all public listings matching the query.
func (s *PublicListingService) SearchPublicListings(query string) ([]*listing.Listing, error) {
	listings, err := s.listingFacade.SearchListings(query)
	if err != nil {
		return nil, err
	}

	var publicListings []*listing.Listing
	for _, listingView := range listings {
		if listingView.AccessMode == listing.Public {
			publicListings = append(publicListings, listingView)
		}
	}

	return publicListings, nil
}

// GetOwnedListing returns a listing only if the user is the owner.
func (s *PublicListingService) GetOwnedListing(userID string, listingID string) (*listing.Listing, error) {
	listingView, err := s.listingFacade.GetListing(listingID)
	if err != nil {
		return nil, err
	}

	if listingView.AuthorID != userID {
		return nil, fmt.Errorf("user %s is not the owner of listing %s", userID, listingID)
	}

	return listingView, nil
}

// GetPublicListing returns a listing only if it is public.
func (s *PublicListingService) GetPublicListing(listingID string) (*listing.Listing, error) {
	listingView, err := s.listingFacade.GetListing(listingID)
	if err != nil {
		return nil, err
	}

	if listingView.AccessMode != listing.Public {
		return nil, fmt.Errorf("listing %s is not public", listingID)
	}

	return listingView, nil
}

// CreateListing creates a new listing.
func (s *PublicListingService) CreateListing(
	authorID string,
	title string,
	description string,
	accessMode listing.ListingAccessMode,
) (*listing.Listing, error) {
	return s.listingFacade.CreateListing(authorID, title, description, accessMode)
}

// UpdateListing updates an owned listing.
func (s *PublicListingService) UpdateListing(
	userID string,
	listingID string,
	update listing.ListingUpdate,
) (*listing.Listing, error) {
	if _, err := s.GetOwnedListing(userID, listingID); err != nil {
		return nil, err
	}

	return s.listingFacade.UpdateListing(listingID, update)
}

// DeleteListing deletes an owned listing.
func (s *PublicListingService) DeleteListing(userID string, listingID string) error {
	if _, err := s.GetOwnedListing(userID, listingID); err != nil {
		return err
	}

	return s.listingFacade.DeleteListing(listingID)
}
