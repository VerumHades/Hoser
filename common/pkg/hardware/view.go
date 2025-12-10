package hardware

// HardwareSpecificationView is read-only for external use.
type HardwareSpecificationView struct {
	CPUCount  int
	RAMBytes  int64
	DiskBytes int64
}

// ToView converts an internal HardwareSpecification to a read-only view.
func (h *HardwareSpecification) ToView() *HardwareSpecificationView {
	return &HardwareSpecificationView{
		CPUCount:  h.cpuCount,
		RAMBytes:  h.ramBytes,
		DiskBytes: h.diskBytes,
	}
}
