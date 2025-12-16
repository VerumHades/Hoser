package app

import "common/pkg/listing"

type ListingSearchService interface {
	IndexListing(listing *listing.Listing) error
	RemoveListing(listingID string) error
	SearchListings(query string) ([]*listing.Listing, error)
}

type ListingFacadeService struct {
	listingService *listing.ListingService
	searchService  ListingSearchService
}

func NewListingFacadeService(ls *listing.ListingService, ss ListingSearchService) *ListingFacadeService {
	return &ListingFacadeService{
		listingService: ls,
		searchService:  ss,
	}
}

// CreateListing creates a listing and indexes it
func (f *ListingFacadeService) CreateListing(authorID, title, description string, accessMode listing.ListingAccessMode) (*listing.Listing, error) {
	listing, err := f.listingService.CreateListing(authorID, title, description, accessMode)
	if err != nil {
		return nil, err
	}

	if err := f.searchService.IndexListing(listing); err != nil {
		return nil, err
	}

	return listing, nil
}

// SearchListings returns listings matching the search query.
// No access control or visibility rules are applied.
func (f *ListingFacadeService) SearchListings(query string) ([]*listing.Listing, error) {
	return f.searchService.SearchListings(query)
}

// ListByAuthor returns all listings authored by the given user ID
func (f *ListingFacadeService) ListByAuthor(authorID string) ([]*listing.Listing, error) {
	return f.listingService.ListByAuthor(authorID)
}

// UpdateListing updates a listing's fields and reindexes it.
func (f *ListingFacadeService) UpdateListing(listingID string, update listing.ListingUpdate) (*listing.Listing, error) {
	// Update listing via service
	listingView, err := f.listingService.UpdateListing(listingID, update)
	if err != nil {
		return nil, err
	}

	err = f.searchService.IndexListing(listingView)
	if err != nil {
		return nil, err
	}
	// Reindex using the view (only public/exported fields)
	return listingView, nil
}

func (s *ListingFacadeService) GetListing(listingID string) (*listing.Listing, error) {
	return s.listingService.GetListing(listingID)
}

// DeleteListing deletes a listing and removes it from the search index
func (f *ListingFacadeService) DeleteListing(listingID string) error {
	// Retrieve the public view for search index removal
	listingView, err := f.listingService.GetListing(listingID)
	if err != nil {
		return err
	}

	// Delete the listing via service
	if err := f.listingService.DeleteListing(listingID); err != nil {
		return err
	}

	// Remove from search index using the view's ID
	return f.searchService.RemoveListing(listingView.ID)
}
