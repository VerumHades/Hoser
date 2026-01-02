package adapters

import (
	"common/pkg/domain/listing"
	userservices "common/pkg/services/user"
	"common/pkg/shared"
	"context"
)

// PlatformAccountsConfig holds all platform account IDs and platform cut percentage.
type PlatformAccountsConfig struct {
	ProfitAccountID         shared.AccountID
	HardwareRentAccountID   shared.AccountID
	HardwareRefundAccountID shared.AccountID
	PlatformCutPercentage   int
}

// AccountServiceAdapter implements both AccountQueryService and accountQueryService.
type AccountServiceAdapter struct {
	config             PlatformAccountsConfig
	listingRepository  listing.ListingQueryRepository
	userAccountService *userservices.UserAccountService // handles user ledger accounts
}

// NewAccountServiceAdapter constructs the adapter.
func NewAccountServiceAdapter(
	config PlatformAccountsConfig,
	listingRepository listing.ListingQueryRepository,
	userAccountSvc *userservices.UserAccountService,
) *AccountServiceAdapter {
	return &AccountServiceAdapter{
		config:             config,
		listingRepository:  listingRepository,
		userAccountService: userAccountSvc,
	}
}

// ----------------- AccountQueryService -----------------

func (a *AccountServiceAdapter) GetPlatformProfitAccountID(ctx context.Context) (shared.AccountID, error) {
	return a.config.ProfitAccountID, nil
}

func (a *AccountServiceAdapter) GetPlatformCutPercentage(ctx context.Context) (int, error) {
	return a.config.PlatformCutPercentage, nil
}

func (a *AccountServiceAdapter) GetListingOwnerAccountID(ctx context.Context, listingID shared.ListingID) (shared.AccountID, error) {
	listing, err := a.listingRepository.GetByID(ctx, listingID)
	if err != nil {
		return "", err
	}

	return a.GetUserAccountID(ctx, listing.AuthorID())
}

func (a *AccountServiceAdapter) GetUserAccountID(ctx context.Context, userID shared.UserID) (shared.AccountID, error) {
	account, err := a.userAccountService.GetUserAccount(ctx, userID)
	if err != nil {
		return "", err
	}

	return account.ID(), nil
}

// ----------------- accountQueryService -----------------

func (a *AccountServiceAdapter) GetPlatformHardwareRentAccountID(ctx context.Context) (shared.AccountID, error) {
	return a.config.HardwareRentAccountID, nil
}

func (a *AccountServiceAdapter) GetPlatformHardwareRefundAccountID(ctx context.Context) (shared.AccountID, error) {
	return a.config.HardwareRefundAccountID, nil
}
