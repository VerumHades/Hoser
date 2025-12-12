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
