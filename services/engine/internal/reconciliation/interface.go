package reconciliation

type Controller interface {
	CreateInstance(specification InstanceCreationSpecification) (string, error)
	ExtendInstanceLifetime(instanceIdentifier string, newExpiration int) error
	DeleteInstance(instanceIdentifier string) error

	ListInstancesByOwner(ownerIdentifier string) ([]InstanceSummary, error)
	GetInstanceForOwner(instanceIdentifier string, ownerIdentifier string) (InstanceDetail, error)
	UpdateInstanceForOwner(instanceIdentifier string, ownerIdentifier string, update InstanceUpdateSpecification) error
}

type InstanceCreationSpecification struct {
	OwnerIdentifier string
	Image           string
	Hardware        HardwareSpecification
	Expiration      int
	ExpirationHook  string
}

type HardwareSpecification struct {
	CPU    int // cores
	Memory int // bytes
	Disk   int // bytes
}

type InstanceUpdateSpecification struct {
	Metadata      map[string]string
	Configuration map[string]any
}

type InstanceSummary struct {
	Identifier string
	Status     InstanceStatus
}

type InstanceDetail struct {
	Identifier string
	Status     InstanceStatus
	Metadata   map[string]string
}

type InstanceStatus string

const (
	InstanceStatusCreating InstanceStatus = "creating"
	InstanceStatusRunning  InstanceStatus = "running"
	InstanceStatusStopped  InstanceStatus = "stopped"
	InstanceStatusDeleting InstanceStatus = "deleting"
	InstanceStatusError    InstanceStatus = "error"
)
