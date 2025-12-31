package localdockerdeployer

import (
	"common/pkg/app"
	"common/pkg/hardware"
	"common/pkg/instance"
	"errors"
	"fmt"
	"os/exec"
	"sync"

	"github.com/google/uuid"
)

// PublicGitHubDockerDeployer implements InstanceDeployer and fetches GitHub setups for build+deploy.
type PublicGitHubDockerDeployer struct {
	mu                   sync.Mutex
	publicListingService *app.PublicListingService
	instanceService      *instance.InstanceService

	instanceIDToContainer map[string]string
	instanceStates        map[string]app.DeployedInstanceState
}

// NewPublicGitHubDockerDeployer creates a new deployer.
func NewPublicGitHubDockerDeployer(
	publicListingService *app.PublicListingService,
	instanceService *instance.InstanceService,
) *PublicGitHubDockerDeployer {
	return &PublicGitHubDockerDeployer{
		publicListingService:  publicListingService,
		instanceService:       instanceService,
		instanceIDToContainer: make(map[string]string),
		instanceStates:        make(map[string]app.DeployedInstanceState),
	}
}

// RefreshHardware is a no-op for local Docker deployer.
func (d *PublicGitHubDockerDeployer) RefreshHardware(instanceID string, newSpec *hardware.HardwareSpecification) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	state, ok := d.instanceStates[instanceID]
	if !ok {
		return errors.New("instance not found")
	}

	state.RuntimeError = nil
	state.StateMessage = "Hardware refreshed"
	d.instanceStates[instanceID] = state
	return nil
}

// GetInstanceState returns current state.
func (d *PublicGitHubDockerDeployer) GetInstanceState(instanceID string) (app.DeployedInstanceState, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	state, ok := d.instanceStates[instanceID]
	if !ok {
		return app.DeployedInstanceState{}, errors.New("instance not found")
	}
	return state, nil
}

// StartInstance fetches the GitHub setup from the listing, builds the image, and deploys it.
// StartInstance triggers an asynchronous build and deployment.
func (d *PublicGitHubDockerDeployer) StartInstance(instanceID string) error {
	d.mu.Lock()
	d.instanceStates[instanceID] = app.DeployedInstanceState{
		State:        app.Building,
		StateMessage: "Queued for build",
	}
	d.mu.Unlock()

	go d.buildAndRunInstance(instanceID)

	return nil
}

type boundedBuffer struct {
	maxBytes int
	buffer   []byte
}

/*
NewBoundedBuffer creates a bounded buffer that retains only the last maxBytes.
*/
func NewBoundedBuffer(maxBytes int) *boundedBuffer {
	return &boundedBuffer{
		maxBytes: maxBytes,
	}
}

/*
Write appends data and truncates old data if necessary.
*/
func (b *boundedBuffer) Write(p []byte) (int, error) {
	b.buffer = append(b.buffer, p...)
	if len(b.buffer) > b.maxBytes {
		b.buffer = b.buffer[len(b.buffer)-b.maxBytes:]
	}
	return len(p), nil
}

/*
String returns buffered output as string.
*/
func (b *boundedBuffer) String() string {
	return string(b.buffer)
}

// buildAndRunInstance performs the git clone, docker build, and container run.
func (d *PublicGitHubDockerDeployer) buildAndRunInstance(instanceID string) {
	updateState := func(state app.DeployedInstanceStateEnum, message string, runtimeError error) {
		d.mu.Lock()
		defer d.mu.Unlock()
		d.instanceStates[instanceID] = app.DeployedInstanceState{
			State:        state,
			StateMessage: message,
			RuntimeError: runtimeError,
		}
	}

	inst, err := d.instanceService.GetInstance(instanceID)
	if err != nil {
		updateState(app.ErrorState, "Failed to fetch instance", err)
		return
	}

	updateState(app.Building, "Fetching GitHub setup", nil)
	setup, err := d.publicListingService.GetSetupForListing(inst.ListingID)
	if err != nil {
		updateState(app.ErrorState, "Failed to fetch GitHub setup", err)
		return
	}

	tempDir := fmt.Sprintf("/tmp/%s", uuid.NewString())
	updateState(app.Building, "Cloning repository", nil)

	if err := exec.Command("git", "clone", setup.RepoURL, tempDir).Run(); err != nil {
		updateState(app.ErrorState, "Failed to clone repository", err)
		return
	}

	imageTag := fmt.Sprintf("%s:%s", inst.ListingID, uuid.NewString())
	updateState(app.Building, "Building Docker image", nil)
	outputBuffer := NewBoundedBuffer(32 * 1024)

	buildCommand := exec.Command("docker", "build", "-t", imageTag, tempDir)
	buildCommand.Stdout = outputBuffer
	buildCommand.Stderr = outputBuffer

	if err := buildCommand.Run(); err != nil {
		updateState(
			app.ErrorState,
			"Docker build failed",
			fmt.Errorf("docker build failed:\n%s", outputBuffer.String()),
		)
		return
	}

	containerName := "instance-" + instanceID
	updateState(app.Building, "Starting container", nil)

	if err := exec.Command("docker", "run", "-d", "--name", containerName, imageTag).Run(); err != nil {
		updateState(app.ErrorState, "Failed to start Docker container", err)
		return
	}

	d.mu.Lock()
	d.instanceIDToContainer[instanceID] = containerName
	d.mu.Unlock()

	updateState(app.Running, "Instance running", nil)
}

// StopInstance stops the Docker container.
func (d *PublicGitHubDockerDeployer) StopInstance(instanceID string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	containerName, ok := d.instanceIDToContainer[instanceID]
	if !ok {
		return errors.New("container not found")
	}

	cmd := exec.Command("docker", "stop", containerName)
	if err := cmd.Run(); err != nil {
		d.instanceStates[instanceID] = app.DeployedInstanceState{State: app.ErrorState, StateMessage: "Failed to stop container", RuntimeError: err}
		return err
	}

	d.instanceStates[instanceID] = app.DeployedInstanceState{State: app.Stopped, StateMessage: "Container stopped", RuntimeError: nil}
	return nil
}
