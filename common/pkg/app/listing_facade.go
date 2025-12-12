package app

import "common/pkg/listing"

type ListingSearchService interface {
	IndexListing(listing *listing.ListingView) error
	RemoveListing(listingID string) error
	SearchListings(query string) ([]*listing.ListingView, error)
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
func (f *ListingFacadeService) CreateListing(authorID, title, description string, accessMode listing.ListingAccessMode) (*listing.ListingView, error) {
	listing, err := f.listingService.CreateListing(authorID, title, description, accessMode)
	if err != nil {
		return nil, err
	}

	if err := f.searchService.IndexListing(listing.ToView()); err != nil {
		return nil, err
	}

	return listing.ToView(), nil
}

// ListByAuthor returns all listings authored by the given user ID
func (f *ListingFacadeService) ListByAuthor(authorID string) ([]*listing.ListingView, error) {
	return f.listingService.ListByAuthor(authorID)
}

// UpdateListing updates a listing's title and reindexes it
func (f *ListingFacadeService) UpdateListing(listingID, newTitle string) error {
	// Update title via service
	if err := f.listingService.UpdateTitle(listingID, newTitle); err != nil {
		return err
	}

	// Retrieve the public view of the updated listing
	listingView, err := f.listingService.GetListingView(listingID)
	if err != nil {
		return err
	}

	// Reindex using the view (only public/exported fields)
	return f.searchService.IndexListing(listingView)
}

// DeleteListing deletes a listing and removes it from the search index
func (f *ListingFacadeService) DeleteListing(listingID string) error {
	// Retrieve the public view for search index removal
	listingView, err := f.listingService.GetListingView(listingID)
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
