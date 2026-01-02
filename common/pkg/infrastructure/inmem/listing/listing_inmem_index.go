package inmemlisting

import (
	"context"
	"errors"
	"sort"
	"strings"
	"sync"

	"common/pkg/domain/listing"
	"common/pkg/shared"
)

// InMemoryListingSearchIndex is an in-memory implementation of listing.ListingSearchIndex.
type InMemoryListingSearchIndex struct {
	mu       sync.RWMutex
	listings map[shared.ListingID]*listing.Listing
}

// NewInMemoryListingSearchIndex constructs a new in-memory search index.
func NewInMemoryListingSearchIndex() *InMemoryListingSearchIndex {
	return &InMemoryListingSearchIndex{
		listings: make(map[shared.ListingID]*listing.Listing),
	}
}

// Index adds a listing to the in-memory index.
func (idx *InMemoryListingSearchIndex) Index(ctx context.Context, l *listing.Listing) error {
	if l == nil {
		return errors.New("listing cannot be nil")
	}

	idx.mu.Lock()
	defer idx.mu.Unlock()
	idx.listings[l.ID()] = l
	return nil
}

// Update replaces an existing listing in the index.
func (idx *InMemoryListingSearchIndex) Update(ctx context.Context, l *listing.Listing) error {
	if l == nil {
		return errors.New("listing cannot be nil")
	}

	idx.mu.Lock()
	defer idx.mu.Unlock()

	if _, exists := idx.listings[l.ID()]; !exists {
		return errors.New("listing not found")
	}

	idx.listings[l.ID()] = l
	return nil
}

// Remove deletes a listing from the index.
func (idx *InMemoryListingSearchIndex) Remove(ctx context.Context, listingID shared.ListingID) error {
	idx.mu.Lock()
	defer idx.mu.Unlock()
	delete(idx.listings, listingID)
	return nil
}

// SearchNextBatch returns a paginated batch of listings matching the query.
// Currently, it performs a simple substring search on listing titles.
func (idx *InMemoryListingSearchIndex) SearchNextBatch(
	ctx context.Context,
	query string,
	request shared.BatchRequest[listing.ListingSearchCursor],
) (listings []*listing.Listing, nextCursor listing.ListingSearchCursor, err error) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	allListings := make([]*listing.Listing, 0, len(idx.listings))
	for _, l := range idx.listings {
		if query == "" || containsIgnoreCase(l.Title(), query) {
			allListings = append(allListings, l)
		}
	}

	// Sort by CreatedAt
	sort.Slice(allListings, func(i, j int) bool {
		return allListings[i].CreatedAt().Before(allListings[j].CreatedAt())
	})

	startIndex := 0
	for i, l := range allListings {
		if l.CreatedAt().After(request.Cursor.LastCreatedAt) {
			startIndex = i
			break
		}
	}

	endIndex := startIndex + request.MaxBatchSize
	if endIndex > len(allListings) {
		endIndex = len(allListings)
	}

	batch := allListings[startIndex:endIndex]

	var newCursor listing.ListingSearchCursor
	if len(batch) > 0 {
		newCursor.LastCreatedAt = batch[len(batch)-1].CreatedAt()
	}

	return batch, newCursor, nil
}

// containsIgnoreCase checks if str contains substr, ignoring case.
func containsIgnoreCase(str, substr string) bool {
	return len(substr) == 0 || (len(str) >= len(substr) &&
		strings.Contains(strings.ToLower(str), strings.ToLower(substr)))
}
