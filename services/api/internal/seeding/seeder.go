package seeding

import (
	"common/pkg/domain/entities/listing"
	"common/pkg/shared"
	"context"
	"fmt"
	"math/rand"
)

/**
 * ListingRepository defines the persistence contract for creating listings.
 */
type ListingRepository interface {
	Create(ctx context.Context, listing *listing.Listing) error
}

/**
 * ListingIndexer defines the contract for indexing listings for search.
 */
type ListingIndexer interface {
	Index(ctx context.Context, listing *listing.Listing) error
}

/**
 * ListingSeeder handles the generation and indexing of mock listing data.
 */
type ListingSeeder struct {
	repository ListingRepository
	indexer    ListingIndexer
	imagePool  []string
}

/**
 * NewListingSeeder initializes a new seeder with required dependencies and image assets.
 */
func NewListingSeeder(
	repository ListingRepository,
	indexer ListingIndexer,
	imagePool []string,
) *ListingSeeder {
	return &ListingSeeder{
		repository: repository,
		indexer:    indexer,
		imagePool:  imagePool,
	}
}

/**
 * Seed generates and persists a specified number of random listings.
 */
func (seeder *ListingSeeder) Seed(ctx context.Context, count int, ownerID string) error {
	for i := 0; i < count; i++ {
		generatedListing, err := seeder.generateRandomListing(ownerID, i)
		if err != nil {
			return err
		}

		err = seeder.persistAndIndex(ctx, generatedListing)
		if err != nil {
			return err
		}
	}
	return nil
}

/**
 * persistAndIndex handles the dual-write to the database and search index.
 */
func (seeder *ListingSeeder) persistAndIndex(ctx context.Context, listing *listing.Listing) error {
	err := seeder.repository.Create(ctx, listing)
	if err != nil {
		return fmt.Errorf("failed to persist listing: %w", err)
	}

	err = seeder.indexer.Index(ctx, listing)
	if err != nil {
		return fmt.Errorf("failed to index listing: %w", err)
	}

	return nil
}

/**
 * generateRandomListing constructs a single listing with randomized realistic data.
 */
func (seeder *ListingSeeder) generateRandomListing(ownerID string, index int) (*listing.Listing, error) {
	title := seeder.getRandomTitle(index)
	description := seeder.getRandomDescription()
	imageHandle := seeder.pickRandomImage()

	newListing, _, err := listing.NewListing(
		shared.UserID(ownerID),
		title,
		description,
		listing.Public,
		&shared.HardwareSpecification{},
		seeder.getRandomPrice(),
		"# Hello World!",
	)

	if err == nil {
		newListing.Mutate().SetScreenshotKeys([]shared.ListingScreenshotID{shared.ListingScreenshotID(imageHandle)}).Apply()
	}

	return newListing, err
}

func (seeder *ListingSeeder) getRandomTitle(index int) string {
	prefixes := []string{"Premium", "Enterprise", "Open Source", "Lite", "Advanced"}
	subjects := []string{"Data Processor", "Image Classifier", "Signal Analyzer", "Neural Bridge"}
	return fmt.Sprintf("%s %s #%d", prefixes[rand.Intn(len(prefixes))], subjects[rand.Intn(len(subjects))], index)
}

func (seeder *ListingSeeder) getRandomDescription() string {
	descriptions := []string{
		"A highly optimized solution for large scale computing.",
		"Designed for low-latency environments and real-time processing.",
		"An experimental module for high-throughput data streams.",
	}
	return descriptions[rand.Intn(len(descriptions))]
}

func (seeder *ListingSeeder) pickRandomImage() string {
	if len(seeder.imagePool) == 0 {
		return ""
	}
	return seeder.imagePool[rand.Intn(len(seeder.imagePool))]
}

func (seeder *ListingSeeder) getRandomPrice() int64 {
	return int64((rand.Intn(90) + 10) * 100)
}
