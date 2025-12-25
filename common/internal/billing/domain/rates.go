package domain

import (
	"errors"
	"time"

	"common/internal/shared"
)

// HardwareCostRate represents the cost of CPU, RAM, and Disk over a period.
type HardwareCostRate struct {
	id        shared.HardwareCostRateID
	cpuCost   shared.Money
	ramCost   shared.Money
	diskCost  shared.Money
	validFrom time.Time
	validTo   *time.Time // nil = currently valid indefinitely
}

// NewHardwareCostRate creates a new hardware cost rate with a generated ID.
func NewHardwareCostRate(cpuCost, ramCost, diskCost shared.Money, validFrom time.Time, validTo *time.Time) (*HardwareCostRate, error) {
	id := shared.HardwareCostRateID(shared.GenerateUUID())
	return NewHardwareCostRateWithID(id, cpuCost, ramCost, diskCost, validFrom, validTo)
}

// NewHardwareCostRateWithID creates a hardware cost rate with an existing ID (for repository hydration).
func NewHardwareCostRateWithID(id shared.HardwareCostRateID, cpuCost, ramCost, diskCost shared.Money, validFrom time.Time, validTo *time.Time) (*HardwareCostRate, error) {
	rate := &HardwareCostRate{
		id:        id,
		cpuCost:   cpuCost,
		ramCost:   ramCost,
		diskCost:  diskCost,
		validFrom: validFrom,
		validTo:   validTo,
	}
	if err := rate.Validate(); err != nil {
		return nil, err
	}
	return rate, nil
}

// ID returns the unique identifier.
func (h *HardwareCostRate) ID() shared.HardwareCostRateID {
	return h.id
}

// CPUCost returns the CPU cost.
func (h *HardwareCostRate) CPUCost() shared.Money {
	return h.cpuCost
}

// RAMCost returns the RAM cost.
func (h *HardwareCostRate) RAMCost() shared.Money {
	return h.ramCost
}

// DiskCost returns the disk cost.
func (h *HardwareCostRate) DiskCost() shared.Money {
	return h.diskCost
}

// ValidFrom returns the start of the validity period.
func (h *HardwareCostRate) ValidFrom() time.Time {
	return h.validFrom
}

// ValidTo returns the end of the validity period.
func (h *HardwareCostRate) ValidTo() *time.Time {
	return h.validTo
}

// SetValidTo updates the end of the validity period.
func (h *HardwareCostRate) SetValidTo(validTo time.Time) {
	h.validTo = &validTo
}

// Validate ensures the hardware cost rate is in a valid state.
func (h *HardwareCostRate) Validate() error {
	if h.id == "" {
		return errors.New("hardware cost rate ID cannot be empty")
	}
	if h.cpuCost.IsZero() && h.ramCost.IsZero() && h.diskCost.IsZero() {
		return errors.New("at least one hardware cost must be non-zero")
	}
	if h.validFrom.IsZero() {
		return errors.New("validFrom cannot be zero")
	}
	if h.validTo != nil && !h.validTo.After(h.validFrom) {
		return errors.New("validTo must be after validFrom")
	}
	return nil
}
