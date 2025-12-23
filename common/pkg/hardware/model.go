package hardware

type HardwareSpecification struct {
	CPUCount  int   `json:"cpu"`
	RAMBytes  int64 `json:"ramBytes"`
	DiskBytes int64 `json:"diskBytes"`
}
