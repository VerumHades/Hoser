package userservices

import (
	"context"
	"fmt"

	"common/pkg/domain/listing"
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
	accountQueryService AccountQueryService

	listingRepository listing.ListingQueryRepository

	transactionRepository      accounting.LedgerTransactionCommandRepository
	transactionQueryRepository accounting.LedgerTransactionQueryRepository
	settlementQueryRepository  accounting.SettlementQueryRepository
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
	if err != nil {
		return false, err
	}
	if latest == nil {
		return false, nil
	}

	if latest.ReferenceType() != accounting.ReferenceTypePurchase {
		return false, nil
	}

	settlement, err := s.settlementQueryRepository.GetLastByLedgerTransactionID(ctx, latest.ID())
	if err != nil {
		return false, err
	}

	return settlement.Status() == accounting.SettlementStatusCompleted, nil
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

	if latest, err := s.transactionQueryRepository.GetLatestByReferenceAndAccount(ctx, accountID, string(listingID)); err != nil {
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

	transaction, err := accounting.NewLedgerTransaction(accounting.ReferenceTypePurchase, string(listingID), entries)
	if err != nil {
		return err
	}
	err = s.transactionRepository.Create(ctx, nil, transaction)
	if err != nil {
		return err
	}

	return nil
}
