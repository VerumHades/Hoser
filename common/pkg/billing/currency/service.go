package currency

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
func (s *CurrencyService) GetCurrency(code string) (*CurrencyView, error) {
	currency, err := s.repo.GetByCode(code)
	if err != nil {
		return nil, err
	}
	return &CurrencyView{
		Code:   currency.Code,
		Name:   currency.Name,
		Symbol: currency.Symbol,
	}, nil
}

// ListCurrencies returns all currencies as views.
func (s *CurrencyService) ListCurrencies() ([]*CurrencyView, error) {
	currencies, err := s.repo.ListAll()
	if err != nil {
		return nil, err
	}

	views := make([]*CurrencyView, len(currencies))
	for i, c := range currencies {
		views[i] = &CurrencyView{
			Code:   c.Code,
			Name:   c.Name,
			Symbol: c.Symbol,
		}
	}
	return views, nil
}
