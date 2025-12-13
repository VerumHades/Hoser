package rates

import (
	"common/pkg/money"
	"time"
)

// HardwareCostRateView is a read-only representation of a hardware cost rate.
type HardwareCostRateView struct {
	ID        string      // unique ID of the rate entry
	CPUCost   money.Money // cost per CPU per hour
	RAMCost   money.Money // cost per byte of RAM per hour
	DiskCost  money.Money // cost per byte of disk per hour
	ValidFrom time.Time   // when the rate starts
	ValidTo   *time.Time  // optional end time, nil if still active
}
