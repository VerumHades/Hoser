package instance

import "common/pkg/hardware"

// InstanceState represents the current state of an instance.
type InstanceState int

const (
	Building InstanceState = iota
	Running
	Stopped
	NotBilled
)

// Instance represents a hardware rental instance derived from a listing.
type Instance struct {
	id                    string
	listingID             string
	billingID             string
	state                 InstanceState
	hardwareSpecification *hardware.HardwareSpecification
}
