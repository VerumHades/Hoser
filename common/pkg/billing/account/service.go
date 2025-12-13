package account

import (
	"common/pkg/billing/payments"
)

type BillingAccountService struct {
	accountRepository BillingAccountRepository
}

func (s *BillingAccountService) CreateAccount(
	ownerID string,
	paymentProvider payments.PaymentProvider,
	providerAccountID string,
) (*BillingAccountView, error) {
	account := NewBillingAccount(
		ownerID,
		paymentProvider,
		providerAccountID,
	)

	err := s.accountRepository.Save(account)
	if err != nil {
		return nil, err
	}

	return account.ToView(), nil
}

func (s *BillingAccountService) ListByOwner(ownerID string) ([]*BillingAccountView, error) {
	accounts, err := s.accountRepository.ListByOwner(ownerID)
	if err != nil {
		return nil, err
	}

	views := make([]*BillingAccountView, len(accounts))
	for i, acct := range accounts {
		views[i] = acct.ToView()
	}

	return views, nil
}

func (s *BillingAccountService) GetAccountView(accountID string) (*BillingAccountView, error) {
	account, err := s.accountRepository.GetByID(accountID)
	if err != nil {
		return nil, err
	}

	return account.ToView(), nil
}

func (s *BillingAccountService) SuspendAccount(accountID string) error {
	account, err := s.accountRepository.GetByID(accountID)
	if err != nil {
		return err
	}
	account.status = BillingAccountStatusSuspended
	return s.accountRepository.Save(account)
}

func (s *BillingAccountService) CloseAccount(accountID string) error {
	account, err := s.accountRepository.GetByID(accountID)
	if err != nil {
		return err
	}
	account.status = BillingAccountStatusClosed
	return s.accountRepository.Save(account)
}
