package rates

import (
	"errors"
	"time"

	"common/pkg/shared"
)

// HardwareResourceType defines which hardware resource the rate applies to.
type HardwareResourceType string

const (
	ResourceCPU  HardwareResourceType = "cpu"
	ResourceRAM  HardwareResourceType = "ram"
	ResourceDisk HardwareResourceType = "disk"
)

// HardwareCostRate represents a cost for a specific hardware resource starting from a given time.
type HardwareCostRate struct {
	id          shared.HardwareCostRateID
	resource    HardwareResourceType
	costInCents int64     // cost in minor units
	validFrom   time.Time // when this rate became effective
}

// NewHardwareCostRate creates a new hardware cost rate with a generated ID.
func NewHardwareCostRate(resource HardwareResourceType, costInCents int64, validFrom time.Time) (*HardwareCostRate, error) {
	id := shared.HardwareCostRateID(shared.GenerateUUID())
	return NewHardwareCostRateWithID(id, resource, costInCents, validFrom)
}

// NewHardwareCostRateWithID creates a rate with an existing ID (for repository hydration).
func NewHardwareCostRateWithID(id shared.HardwareCostRateID, resource HardwareResourceType, costInCents int64, validFrom time.Time) (*HardwareCostRate, error) {
	rate := &HardwareCostRate{
		id:          id,
		resource:    resource,
		costInCents: costInCents,
		validFrom:   validFrom,
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

// Resource returns which hardware resource this rate applies to.
func (h *HardwareCostRate) Resource() HardwareResourceType {
	return h.resource
}

// CostInCents returns the cost in minor units.
func (h *HardwareCostRate) CostInCents() int64 {
	return h.costInCents
}

// ValidFrom returns when this rate became effective.
func (h *HardwareCostRate) ValidFrom() time.Time {
	return h.validFrom
}

// Validate ensures the hardware cost rate is valid.
func (h *HardwareCostRate) Validate() error {
	if h.id == "" {
		return errors.New("hardware cost rate ID cannot be empty")
	}
	if h.costInCents <= 0 {
		return errors.New("hardware cost must be greater than zero")
	}
	switch h.resource {
	case ResourceCPU, ResourceRAM, ResourceDisk:
		// valid
	default:
		return errors.New("invalid hardware resource type")
	}
	if h.validFrom.IsZero() {
		return errors.New("validFrom cannot be zero")
	}
	return nil
}
