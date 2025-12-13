package rates

import (
	"common/pkg/money"
	"time"
)

type HardwareCostRate struct {
	ID        string
	CPUCost   money.Money
	RAMCost   money.Money
	DiskCost  money.Money
	ValidFrom time.Time
	ValidTo   *time.Time // nil = currently valid indefinitely
}
