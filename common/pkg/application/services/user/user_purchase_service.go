package userservices

import (
	"context"
	"fmt"

	"common/pkg/application/unitofwork"
	"common/pkg/domain/entities/accounting"
	"common/pkg/domain/entities/events"
	"common/pkg/domain/repositories"
	"common/pkg/shared"
)

type AccountQueryService interface {
	GetPlatformProfitAccountID(ctx context.Context) (shared.AccountID, error)
	GetPlatformCutPercentage(ctx context.Context) (int, error)

	GetListingOwnerAccountID(ctx context.Context, listingID shared.ListingID) (shared.AccountID, error)
	GetUserAccountID(ctx context.Context, userID shared.UserID) (shared.AccountID, error)
}

// UserPurchaseService handles all user listing purchase logic.
type UserPurchaseService struct {
	accountQueryService         AccountQueryService
	transactionalEventPublisher *unitofwork.TransactionalEventPublisher

	listingRepository repositories.ListingQueryRepository

	transactionRepository      repositories.LedgerTransactionCommandRepository
	transactionQueryRepository repositories.LedgerTransactionQueryRepository
	settlementQueryRepository  repositories.SettlementQueryRepository
}

func NewUserPurchaseService(
	accountQueryService AccountQueryService,
	transactionalEventPublisher *unitofwork.TransactionalEventPublisher,
	listingRepository repositories.ListingQueryRepository,
	transactionRepository repositories.LedgerTransactionCommandRepository,
	transactionQueryRepository repositories.LedgerTransactionQueryRepository,
	settlementQueryRepository repositories.SettlementQueryRepository,
) *UserPurchaseService {
	return &UserPurchaseService{
		accountQueryService:         accountQueryService,
		transactionalEventPublisher: transactionalEventPublisher,
		listingRepository:           listingRepository,
		transactionRepository:       transactionRepository,
		transactionQueryRepository:  transactionQueryRepository,
		settlementQueryRepository:   settlementQueryRepository,
	}
}

func (s *UserPurchaseService) GetLastestUserListingTransaction(ctx context.Context, userID shared.UserID, listingID shared.ListingID) (*accounting.LedgerTransaction, error) {
	accountID, err := s.accountQueryService.GetUserAccountID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user has no ledger account")
	}
	return s.transactionQueryRepository.GetLatestByReferenceAndAccount(ctx, accountID, string(listingID))
}

// HasUserBoughtListing checks if a user already owns a listing.
func (s *UserPurchaseService) DoesUserOwnListing(ctx context.Context, userID shared.UserID, listingID shared.ListingID) (bool, error) {
	latest, err := s.GetLastestUserListingTransaction(ctx, userID, listingID)
	if err == shared.ErrNotFound {
		return false, nil
	} else if err != nil {
		return false, err
	}

	if latest.ReferenceType() != accounting.ReferenceTypePurchase {
		return false, nil
	}

	settlement, err := s.settlementQueryRepository.GetLastByLedgerTransactionID(ctx, latest.ID())
	if err == shared.ErrNotFound {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return settlement.Status() == accounting.SettlementStatusCompleted, nil
}

// HasUserBoughtListing checks if a user already owns a listing.
func (s *UserPurchaseService) IsPurchaseProcessing(ctx context.Context, userID shared.UserID, listingID shared.ListingID) (bool, error) {
	latest, err := s.GetLastestUserListingTransaction(ctx, userID, listingID)
	if err == shared.ErrNotFound {
		return false, nil
	} else if err != nil {
		return false, err
	}

	if latest.ReferenceType() != accounting.ReferenceTypePurchase {
		return false, nil
	}

	settlement, err := s.settlementQueryRepository.GetLastByLedgerTransactionID(ctx, latest.ID())
	if err == shared.ErrNotFound {
		return true, nil
	} else if err != nil {
		return false, err
	}

	return settlement.Status() != accounting.SettlementStatusAbbandoned && settlement.Status() != accounting.SettlementStatusCompleted, nil
}

// PurchaseListing ensures a user owns a listing and executes the payment if necessary.
func (s *UserPurchaseService) PurchaseListing(
	ctx context.Context,
	userID shared.UserID,
	listingID shared.ListingID,
) error {
	listing, err := s.listingRepository.GetByID(ctx, listingID)
	if err != nil {
		return fmt.Errorf("failed to check listing existence: %w", err)
	}
	if listing == nil {
		return fmt.Errorf("listing %s does not exist", listingID)
	}

	accountID, err := s.accountQueryService.GetUserAccountID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user has no ledger account")
	}

	if latest, err := s.transactionQueryRepository.GetLatestByReferenceAndAccount(ctx, accountID, string(listingID)); err != nil && err != shared.ErrNotFound {
		return err
	} else if latest != nil && latest.ReferenceType() == accounting.ReferenceTypePurchase {
		return fmt.Errorf("an existing identical purchase is already processing")
	}

	if hasBought, err := s.DoesUserOwnListing(ctx, userID, listingID); err != nil {
		return fmt.Errorf("failed to check ownership: %w", err)
	} else if hasBought {
		return fmt.Errorf("user %s has already purchased listing %s", userID, listingID)
	}

	platformCutPercentage, err := s.accountQueryService.GetPlatformCutPercentage(ctx)

	purchaseAmount := listing.PriceInMinorUnits()
	platformCut := (purchaseAmount * int64(platformCutPercentage)) / 100
	sellerAmount := purchaseAmount - platformCut

	platformAccountID, err := s.accountQueryService.GetPlatformProfitAccountID(ctx)
	if err != nil {
		return err
	}
	sellerAccountID, err := s.accountQueryService.GetListingOwnerAccountID(ctx, listingID)
	if err != nil {
		return err
	}

	entries, err := accounting.NewLedgerEntries(
		accounting.LedgerEntryShorthand(accountID, -purchaseAmount),
		accounting.LedgerEntryShorthand(platformAccountID, platformCut),
		accounting.LedgerEntryShorthand(sellerAccountID, sellerAmount),
	)

	if err != nil {
		return err
	}

	return s.transactionalEventPublisher.PublishWithTransaction(ctx, func(txContext context.Context) ([]events.DomainEvent, error) {
		transaction, event, err := accounting.NewLedgerTransaction(accounting.ReferenceTypePurchase, string(listingID), entries)
		if err != nil {
			return []events.DomainEvent{}, err
		}
		err = s.transactionRepository.Create(ctx, transaction)
		if err != nil {
			return []events.DomainEvent{}, err
		}

		return []events.DomainEvent{event}, nil
	})
}
