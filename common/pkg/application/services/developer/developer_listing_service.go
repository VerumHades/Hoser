package developer

import (
	"common/pkg/application/unitofwork"
	"common/pkg/domain/entities/events"
	"common/pkg/domain/entities/listing"
	"common/pkg/domain/repositories"
	"common/pkg/shared"
	"common/pkg/util"
	"context"
)

// DeveloperListingService exposes operations for managing listings and GitHub setups.
type DeveloperListingService struct {
	transactionalEventPublisher unitofwork.TransactionalEventPublisher

	listingRepository      repositories.ListingCommandRepository
	listingQueryRepository repositories.ListingQueryRepository
	githubSetupRepository  repositories.GitHubSetupCommandRepository
}

// NewDeveloperListingService constructs a new DeveloperListingService.
func NewDeveloperListingService(
	transactionalEventPublisher unitofwork.TransactionalEventPublisher,
	listingRepository repositories.ListingCommandRepository,
	listingQueryRepository repositories.ListingQueryRepository,
	githubSetupRepository repositories.GitHubSetupCommandRepository,
) *DeveloperListingService {
	return &DeveloperListingService{
		transactionalEventPublisher: transactionalEventPublisher,
		listingRepository:           listingRepository,
		listingQueryRepository:      listingQueryRepository,
		githubSetupRepository:       githubSetupRepository,
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
	listing, event, err := listing.NewListing(req.AuthorID, req.Title, req.Description, req.AccessMode, req.Hardware, req.PriceInMinorUnits)

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
	request shared.BatchRequest[repositories.ListingCursor],
) (listings []*listing.Listing, nextCursor repositories.ListingCursor, err error) {
	return s.listingQueryRepository.FetchNextBatchByAuthor(ctx, authorID, request)
}
