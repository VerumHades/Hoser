package currency

// CurrencyRepository defines operations for persisting and retrieving currencies.
type CurrencyRepository interface {
	// Save stores or updates a currency.
	Save(currency *Currency) error

	// GetByCode retrieves a currency by its code, e.g., "USD".
	GetByCode(code string) (*Currency, error)

	// ListAll returns all stored currencies.
	ListAll() ([]*Currency, error)
}
