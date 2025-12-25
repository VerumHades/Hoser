package shared

type HardwareSpecification struct {
	CPUCount  int   `json:"cpu"`
	RAMBytes  int64 `json:"ramBytes"`
	DiskBytes int64 `json:"diskBytes"`
}

// Equals returns true if all fields of two hardware specifications are identical.
func (h *HardwareSpecification) Equals(other *HardwareSpecification) bool {
	if h == nil || other == nil {
		return false
	}
	return h.CPUCount == other.CPUCount &&
		h.RAMBytes == other.RAMBytes &&
		h.DiskBytes == other.DiskBytes
}

// IsUpgradeOf returns true if the current hardware specification is greater than or equal
// to the other specification in all resources (CPU, RAM, Disk). Useful for upgrade checks.
func (h *HardwareSpecification) IsUpgradeOf(other *HardwareSpecification) bool {
	if h == nil || other == nil {
		return false
	}
	return h.CPUCount >= other.CPUCount &&
		h.RAMBytes >= other.RAMBytes &&
		h.DiskBytes >= other.DiskBytes
}

// IsDowngradeOf returns true if the current hardware specification is less than or equal
// to the other specification in all resources.
func (h *HardwareSpecification) IsDowngradeOf(other *HardwareSpecification) bool {
	if h == nil || other == nil {
		return false
	}
	return h.CPUCount <= other.CPUCount &&
		h.RAMBytes <= other.RAMBytes &&
		h.DiskBytes <= other.DiskBytes
}
