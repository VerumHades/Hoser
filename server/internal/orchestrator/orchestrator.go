package orchestrator

type DeploymentState int
type DeploymentID int

const (
	StateRunning    DeploymentState = iota // Running without issues
	StateCrashed                           // Crashed, waiting for restart
	StateError                             // Crashed with an error such that reruning wont fix the issue
	StateRestarting                        // Restarting
	StateStopped                           // Stopped
)

type Deployable interface {
	Run() error
	Stop() error
}

type Deployment struct {
	ID         DeploymentID
	state      DeploymentState
	deployable *Deployable
}

type Orchestrator interface {
	Deploy() (*Deployable, error)
	Stop(id *Deployment) error
}
