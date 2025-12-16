package currency

type Currency struct {
	Code   string // "USD"
	Name   string // "US Dollar"
	Symbol string // optional "$"
}

// CurrencyRepository defines operations for persisting and retrieving currencies.
type CurrencyRepository interface {
	Save(currency *Currency) error
	GetByCode(code string) (*Currency, error)
	ListAll() ([]*Currency, error)
}

// CurrencyService provides higher-level operations around currencies.
type CurrencyService struct {
	repo CurrencyRepository
}

// NewCurrencyService creates a new service.
func NewCurrencyService(repo CurrencyRepository) *CurrencyService {
	return &CurrencyService{repo: repo}
}

// CreateCurrency validates and stores a new currency.
func (s *CurrencyService) CreateCurrency(code, name, symbol string) (*Currency, error) {
	currency := &Currency{
		Code:   code,
		Name:   name,
		Symbol: symbol,
	}
	err := s.repo.Save(currency)
	if err != nil {
		return nil, err
	}
	return currency, nil
}

// GetCurrency retrieves a currency by code.
func (s *CurrencyService) GetCurrency(code string) (*Currency, error) {
	return s.repo.GetByCode(code)
}

// ListCurrencies returns all currencies as views.
func (s *CurrencyService) ListCurrencies() ([]*Currency, error) {
	return s.repo.ListAll()
}
