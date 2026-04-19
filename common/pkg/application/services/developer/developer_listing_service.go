package developer

import (
	"common/pkg/application/unitofwork"
	"common/pkg/domain/entities/events"
	"common/pkg/domain/entities/listing"
	"common/pkg/domain/repositories"
	"common/pkg/shared"
	"common/pkg/util"
	"context"
	"time"
)

// DeveloperListingService exposes operations for managing listings and GitHub setups.
type DeveloperListingService struct {
	transactionalEventPublisher unitofwork.TransactionalEventPublisher
	storageProvider             StorageProvider

	listingRepository      repositories.ListingCommandRepository
	listingQueryRepository repositories.ListingQueryRepository
	githubSetupRepository  repositories.GitHubSetupCommandRepository
	listingSearchIndex     listing.ListingSearchIndex
}

// NewDeveloperListingService constructs a new DeveloperListingService.
func NewDeveloperListingService(
	transactionalEventPublisher unitofwork.TransactionalEventPublisher,
	storageProvider StorageProvider,
	listingRepository repositories.ListingCommandRepository,
	listingQueryRepository repositories.ListingQueryRepository,
	githubSetupRepository repositories.GitHubSetupCommandRepository,
	listingSearchIndex listing.ListingSearchIndex,
) *DeveloperListingService {
	return &DeveloperListingService{
		transactionalEventPublisher: transactionalEventPublisher,
		storageProvider:             storageProvider,
		listingRepository:           listingRepository,
		listingQueryRepository:      listingQueryRepository,
		githubSetupRepository:       githubSetupRepository,
		listingSearchIndex:          listingSearchIndex,
	}
}

type CreateListingRequest struct {
	AuthorID          shared.UserID
	Title             string
	Description       string
	AccessMode        listing.ListingAccessMode
	Hardware          *shared.HardwareSpecification
	PriceInMinorUnits int64
}

func (s *DeveloperListingService) CreateListing(ctx context.Context, req CreateListingRequest) (l *listing.Listing, err error) {
	listing, event, err := listing.NewListing(req.AuthorID, req.Title, req.Description, req.AccessMode, req.Hardware, req.PriceInMinorUnits, "")

	return listing, s.transactionalEventPublisher.PublishWithTransaction(ctx, func(txCtx context.Context) ([]events.DomainEvent, error) {

		if err != nil {
			return nil, err
		}
		if err := s.listingRepository.Create(txCtx, listing); err != nil {
			return nil, err
		}
		return []events.DomainEvent{event}, nil // or listing.Mutate().Apply() depending on your domain design
	})
}

type ListingMutationFunction func(ctx context.Context, mutator *listing.ListingMutationBuilder) error

// UpdateListing updates a listing and reindexes it.
func (s *DeveloperListingService) UpdateListing(ctx context.Context, listingID shared.ListingID, updateFunction ListingMutationFunction) (*listing.Listing, error) {
	existingListing, err := util.GetExistingEntity(ctx, listingID, s.listingQueryRepository)
	if err != nil {
		return nil, err
	}

	err = s.transactionalEventPublisher.PublishWithTransaction(ctx, func(txContext context.Context) ([]events.DomainEvent, error) {
		mutator := existingListing.Mutate()
		err := updateFunction(txContext, mutator)

		if err != nil {
			return nil, err
		}

		event, err := mutator.Apply()

		if err != nil {
			return nil, err
		}

		if err := s.listingRepository.Update(txContext, existingListing); err != nil {
			return nil, err
		}

		return []events.DomainEvent{event}, nil
	})

	if err != nil {
		return nil, err
	}

	return existingListing, nil
}

// DeleteListing removes a listing and deletes it from the search index.
func (s *DeveloperListingService) DeleteListing(ctx context.Context, listingID shared.ListingID) error {
	return s.transactionalEventPublisher.PublishWithTransaction(ctx, func(txCtx context.Context) ([]events.DomainEvent, error) {
		err := s.listingRepository.Delete(txCtx, listingID)

		if err != nil {
			return nil, err
		}
		event := events.NewDomainEventEnvelope(listing.ListingDeleteEvent{ID: listingID})
		return []events.DomainEvent{event}, nil // or listing.Mutate().Apply() depending on your domain design
	})
}

func (s *DeveloperListingService) GetOwnedListing(ctx context.Context, listingID shared.ListingID, userID shared.UserID) (*listing.Listing, error) {
	return s.listingQueryRepository.GetByIDAndAuthor(ctx, listingID, userID)
}

func (s *DeveloperListingService) FetchNextBatchByAuthor(
	ctx context.Context,
	authorID shared.UserID,
	query listing.SearchQuery,
	request shared.BatchRequest[listing.ListingSearchCursor],
) (listings []*listing.Listing, nextCursor listing.ListingSearchCursor, err error) {
	query.AuthorID = &authorID
	return s.listingSearchIndex.SearchNextBatch(ctx, query, request)
}

type StorageProvider interface {
	GenerateUploadLink(ctx context.Context, key string, maxSize int64, expires time.Duration) (string, error)
	DeleteObject(ctx context.Context, key string) error
}

func (s *DeveloperListingService) CreateScreenshotUploadUrl(
	ctx context.Context,
	listingID shared.ListingID,
	userID shared.UserID,
) (string, error) {
	existingListing, err := s.listingQueryRepository.GetByIDAndAuthor(ctx, listingID, userID)
	if err != nil {
		return "", err
	}

	const MaxFileSize = 10 * 1024 * 1024
	const LinkExpiry = 15 * time.Minute
	const MaxScreenshots = 5

	if len(existingListing.ScreenshotKeys()) >= MaxScreenshots {
		return "", shared.ErrLimitReached
	}

	screenshotId := shared.ListingScreenshotID(shared.GenerateUUID())
	fileKey := shared.GenerateFileKey(listingID, screenshotId)

	uploadURL, err := s.storageProvider.GenerateUploadLink(ctx, fileKey, MaxFileSize, LinkExpiry)
	if err != nil {
		return "", err
	}

	keys := existingListing.ScreenshotKeys()
	keys = append(keys, screenshotId)

	existingListing.Mutate().SetScreenshotKeys(keys).Apply()
	err = s.listingRepository.Update(ctx, existingListing)

	return uploadURL, err
}

func (s *DeveloperListingService) DeleteScreenshot(
	ctx context.Context,
	listingID shared.ListingID,
	userID shared.UserID,
	screenshotId shared.ListingScreenshotID,
) error {
	// 1. Fetch the listing to ensure ownership and that the screenshot exists
	existingListing, err := s.listingQueryRepository.GetByIDAndAuthor(ctx, listingID, userID)
	if err != nil {
		return err
	}

	fileKey := shared.GenerateFileKey(listingID, screenshotId)

	err = s.transactionalEventPublisher.PublishWithTransaction(ctx, func(txCtx context.Context) ([]events.DomainEvent, error) {
		mutator := existingListing.Mutate()

		currentKeys := existingListing.ScreenshotKeys()
		newKeys := make([]shared.ListingScreenshotID, 0, len(currentKeys))
		for _, k := range currentKeys {
			if k != screenshotId {
				newKeys = append(newKeys, k)
			}
		}

		event, err := mutator.SetScreenshotKeys(newKeys).Apply()
		if err != nil {
			return nil, err
		}

		if err := s.listingRepository.Update(txCtx, existingListing); err != nil {
			return nil, err
		}

		return []events.DomainEvent{event}, nil
	})

	if err != nil {
		return err
	}

	return s.storageProvider.DeleteObject(ctx, fileKey)
}
