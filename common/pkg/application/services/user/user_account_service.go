package userservices

import (
	"context"
	"errors"

	"common/pkg/shared"
)

// UserAccountService provides simple operations for user ledger accounts.
type UserAccountService struct {
	accountQueryRepo   accounting.AccountQueryRepository
	accountCommandRepo accounting.AccountCommandRepository
}

// NewUserAccountService constructs the service.
func NewUserAccountService(
	accountQueryRepo accounting.AccountQueryRepository,
	accountCommandRepo accounting.AccountCommandRepository,
) *UserAccountService {
	return &UserAccountService{
		accountQueryRepo:   accountQueryRepo,
		accountCommandRepo: accountCommandRepo,
	}
}

// GetUserAccount retrieves the first ledger account for a given shared.
func (s *UserAccountService) GetUserAccount(ctx context.Context, userID shared.UserID) (*accounting.Account, error) {
	if userID == "" {
		return nil, errors.New("userID cannot be empty")
	}

	account, err := s.accountQueryRepo.GetFirstByOwner(ctx, accounting.AccountOwnerUser, string(userID))
	if err != nil {
		return nil, err
	}
	return account, nil
}

// CreateUserAccount creates a new ledger account for a user with the given account ID and type.
func (s *UserAccountService) CreateUserAccount(
	ctx context.Context,
	accountID shared.AccountID,
	userID shared.UserID,
	accountType accounting.AccountType,

) (*accounting.Account, error) {
	if accountID == "" || userID == "" {
		return nil, errors.New("accountID and userID cannot be empty")
	}

	account, err := accounting.NewAccount(accountID, accountType, accounting.AccountOwnerUser, string(userID))
	if err != nil {
		return nil, err
	}

	return s.accountCommandRepo.Create(ctx, account)
}

// UserAccountExists checks whether a ledger account exists for the shared.
func (s *UserAccountService) UserAccountExists(ctx context.Context, userID shared.UserID) (bool, error) {
	if userID == "" {
		return false, errors.New("userID cannot be empty")
	}

	accounts, err := s.accountQueryRepo.GetByOwner(ctx, accounting.AccountOwnerUser, string(userID))
	if err != nil {
		return false, err
	}
	return len(accounts) > 0, nil
}
