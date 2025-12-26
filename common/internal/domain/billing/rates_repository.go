package billing

import (
	"time"
)

type HardwareCostRepository interface {
	Save(rate *HardwareCostRate) (*HardwareCostRate, error)

	GetActiveRate(at time.Time) (*HardwareCostRate, error)

	FetchNextBatchOrderedByEffectiveDate(
		lastSeenEffectiveDate time.Time,
		maximumBatchSize int,
	) ([]*HardwareCostRate, error)
}
