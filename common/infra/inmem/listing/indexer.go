package listingmem

import (
	"strings"
	"sync"

	"common/pkg/listing"
)

// InMemoryListingSearchService is a simple in-memory implementation of ListingSearchService.
type InMemoryListingSearchService struct {
	mu       sync.RWMutex
	listings map[string]*listing.Listing
}

// NewInMemoryListingSearchService creates a new in-memory search service.
func NewInMemoryListingSearchService() *InMemoryListingSearchService {
	return &InMemoryListingSearchService{
		listings: make(map[string]*listing.Listing),
	}
}

// IndexListing stores the listing in memory.
func (s *InMemoryListingSearchService) IndexListing(listingItem *listing.Listing) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.listings[listingItem.ID] = listingItem
	return nil
}

// RemoveListing removes a listing by ID.
func (s *InMemoryListingSearchService) RemoveListing(listingID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.listings, listingID)
	return nil
}

// SearchListings returns listings whose title or description contains the query (case-insensitive).
func (s *InMemoryListingSearchService) SearchListings(query string) ([]*listing.Listing, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []*listing.Listing
	query = strings.ToLower(query)
	for _, l := range s.listings {
		if strings.Contains(strings.ToLower(l.Title), query) || strings.Contains(strings.ToLower(l.Description), query) {
			results = append(results, l)
		}
	}
	return results, nil
}
