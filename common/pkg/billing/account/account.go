package account

import (
	"time"

	"common/pkg/billing/payments"
	"common/pkg/util"
)

type BillingAccountStatus string

const (
	BillingAccountStatusActive    BillingAccountStatus = "active"
	BillingAccountStatusSuspended BillingAccountStatus = "suspended"
	BillingAccountStatusClosed    BillingAccountStatus = "closed"
)

type BillingAccount struct {
	ID                string
	OwnerID           string
	Status            BillingAccountStatus
	PaymentProvider   payments.PaymentProvider
	ProviderAccountID string
	CreatedAt         time.Time
}

func NewBillingAccount(
	ownerID string,
	paymentProvider payments.PaymentProvider,
	providerAccountID string,
) *BillingAccount {
	return &BillingAccount{
		ID:                util.GenerateUUID(),
		OwnerID:           ownerID,
		Status:            BillingAccountStatusActive,
		PaymentProvider:   paymentProvider,
		ProviderAccountID: providerAccountID,
		CreatedAt:         time.Now(),
	}
}

func (b *BillingAccount) IsChargeable() bool {
	return b.Status == BillingAccountStatusActive
}

func (b *BillingAccount) Suspend() {
	b.Status = BillingAccountStatusSuspended
}

func (b *BillingAccount) Close() {
	b.Status = BillingAccountStatusClosed
}

type BillingAccountRepository interface {
	Save(account *BillingAccount) error

	GetByID(accountID string) (*BillingAccount, error)
	ListByOwner(ownerID string) ([]*BillingAccount, error)
}

type BillingAccountService struct {
	accountRepository BillingAccountRepository
}

func NewBillingAccountService(repo BillingAccountRepository) *BillingAccountService {
	return &BillingAccountService{
		accountRepository: repo,
	}
}

func (s *BillingAccountService) CreateAccount(
	ownerID string,
	paymentProvider payments.PaymentProvider,
	providerAccountID string,
) (*BillingAccount, error) {
	account := NewBillingAccount(
		ownerID,
		paymentProvider,
		providerAccountID,
	)

	err := s.accountRepository.Save(account)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (s *BillingAccountService) ListByOwner(ownerID string) ([]*BillingAccount, error) {
	return s.accountRepository.ListByOwner(ownerID)
}

func (s *BillingAccountService) GetAccount(accountID string) (*BillingAccount, error) {
	return s.accountRepository.GetByID(accountID)
}

func (s *BillingAccountService) SuspendAccount(accountID string) error {
	account, err := s.accountRepository.GetByID(accountID)
	if err != nil {
		return err
	}
	account.Status = BillingAccountStatusSuspended
	return s.accountRepository.Save(account)
}

func (s *BillingAccountService) CloseAccount(accountID string) error {
	account, err := s.accountRepository.GetByID(accountID)
	if err != nil {
		return err
	}
	account.Status = BillingAccountStatusClosed
	return s.accountRepository.Save(account)
}
