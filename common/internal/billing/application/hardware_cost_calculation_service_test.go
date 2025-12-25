package application

import (
	"errors"
	"testing"
	"time"

	"common/internal/billing/domain"
	"common/internal/shared"

	"github.com/stretchr/testify/assert"
)

/* Stub repository to return predefined rates or errors */
type mockHardwareCostRepository struct {
	rate *domain.HardwareCostRate
	err  error
}

func (m *mockHardwareCostRepository) GetActiveRate(at time.Time) (*domain.HardwareCostRate, error) {
	return m.rate, m.err
}

func (m *mockHardwareCostRepository) Save(rate *domain.HardwareCostRate) error {
	return nil
}

func (m *mockHardwareCostRepository) ListAll() ([]*domain.HardwareCostRate, error) {
	return []*domain.HardwareCostRate{}, nil
}

/* Stub currency conversion service */
type mockCurrencyConversionService struct{}

func (m *mockCurrencyConversionService) Convert(money shared.Money, targetCurrency string) (shared.Money, error) {
	if targetCurrency == "ERR" {
		return shared.Money{}, errors.New("conversion error")
	}
	// Simulate conversion: multiply amount by 2 for testing
	return shared.Money{
		Amount:       money.Amount * 2,
		CurrencyCode: targetCurrency,
	}, nil
}

func TestCalculateCost_EdgeCases(t *testing.T) {
	now := time.Now()
	next := now.Add(time.Hour)

	// Setup a controlled hardware cost rate
	rate, _ := domain.NewHardwareCostRate(
		shared.Money{Amount: 1.0, CurrencyCode: "USD"}, // CPU cost per hour
		shared.Money{Amount: 0.5, CurrencyCode: "USD"}, // RAM cost per GB-hour
		shared.Money{Amount: 0.1, CurrencyCode: "USD"}, // Disk cost per GB-hour
		now, &next,
	)

	repo := &mockHardwareCostRepository{rate: rate}
	conversion := &mockCurrencyConversionService{}
	service := NewHardwareCostCalculationService(repo, conversion)

	tests := []struct {
		name           string
		spec           *shared.HardwareSpecification
		duration       time.Duration
		targetCurrency string
		expectedAmount float64
		expectError    bool
	}{
		{
			name:           "normal usage",
			spec:           &shared.HardwareSpecification{CPUCount: 2, RAMBytes: 4, DiskBytes: 10},
			duration:       time.Hour,
			targetCurrency: "USD",
			expectedAmount: 2*1.0 + 4*0.5 + 10*0.1, // 2 + 2 + 1 = 5
			expectError:    false,
		},
		{
			name:           "zero duration",
			spec:           &shared.HardwareSpecification{CPUCount: 2, RAMBytes: 4, DiskBytes: 10},
			duration:       0,
			targetCurrency: "USD",
			expectedAmount: 0.0,
			expectError:    false,
		},
		{
			name:           "currency conversion doubles amount",
			spec:           &shared.HardwareSpecification{CPUCount: 1, RAMBytes: 2, DiskBytes: 5},
			duration:       time.Hour,
			targetCurrency: "EUR",
			expectedAmount: (1*1.0 + 2*0.5 + 5*0.1) * 2, // simulated conversion multiplies by 2
			expectError:    false,
		},
		{
			name:           "conversion failure triggers error",
			spec:           &shared.HardwareSpecification{CPUCount: 1, RAMBytes: 1, DiskBytes: 1},
			duration:       time.Hour,
			targetCurrency: "ERR",
			expectError:    true,
		},
		{
			name:           "large hardware specification",
			spec:           &shared.HardwareSpecification{CPUCount: 1000, RAMBytes: 2000, DiskBytes: 5000},
			duration:       time.Hour,
			targetCurrency: "USD",
			expectedAmount: 1000*1.0 + 2000*0.5 + 5000*0.1, // 1000 + 1000 + 500 = 2500
			expectError:    false,
		},
		{
			name:           "negative hardware specification triggers zero or error",
			spec:           &shared.HardwareSpecification{CPUCount: -1, RAMBytes: -2, DiskBytes: -3},
			duration:       time.Hour,
			targetCurrency: "USD",
			expectedAmount: 0.0, // assume service clamps negative values to zero
			expectError:    false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			total, err := service.CalculateCost(test.spec, test.duration, now, test.targetCurrency)
			if test.expectError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, test.targetCurrency, total.CurrencyCode)
			assert.InDelta(t, test.expectedAmount, total.Amount, 0.0001)
		})
	}
}

func TestCalculateCost_RepositoryUnavailable(t *testing.T) {
	repo := &mockHardwareCostRepository{err: errors.New("repo failure")}
	conversion := &mockCurrencyConversionService{}
	service := NewHardwareCostCalculationService(repo, conversion)

	spec := &shared.HardwareSpecification{CPUCount: 1, RAMBytes: 1, DiskBytes: 1}
	total, err := service.CalculateCost(spec, time.Hour, time.Now(), "USD")
	assert.Nil(t, total)
	assert.Error(t, err)
}
