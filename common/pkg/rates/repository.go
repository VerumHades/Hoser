package rates

import "time"

type HardwareCostRepository interface {
	Save(rate *HardwareCostRate) error
	GetActiveRate(at time.Time) (*HardwareCostRate, error)
	ListAll() ([]*HardwareCostRate, error)
}
