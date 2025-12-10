package reconciliation

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/google/uuid"
)

type DockerController struct{}

func (controller DockerController) CreateInstance(specification InstanceCreationSpecification) (string, error) {
	instanceIdentifier := uuid.NewString()

	createError := controller.createContainer(
		instanceIdentifier,
		specification.Image,
	)

	if createError != nil {
		return "", createError
	}

	return instanceIdentifier, nil
}

func (controller DockerController) ListInstancesByOwner(ownerIdentifier string) ([]InstanceSummary, error) {
	cmd := exec.Command("docker", "ps", "-a", "--format", "{{.Names}}")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")

	var instances []InstanceSummary

	for _, name := range lines {
		if name == "" {
			continue
		}

		if !strings.HasPrefix(name, ownerIdentifier+"-") {
			continue
		}

		status, statusError := controller.inspectStatus(name)
		if statusError != nil {
			continue
		}

		instances = append(instances, InstanceSummary{
			Identifier: name,
			Status:     status,
		})
	}

	return instances, nil
}

func (controller DockerController) ExtendInstanceLifetime(instanceIdentifier string, newExpiration int) error {
	// Retrieve existing instance metadata...
	instance, err := controller.GetInstanceForOwner(instanceIdentifier, "owner-placeholder")
	if err != nil {
		return err
	}

	if newExpiration <= int(time.Now().Unix()) {
		return fmt.Errorf("expiration must be in the future")
	}

	// Update the expiration in your storage (metadata, database, etc.)
	instance.Metadata["expiration"] = fmt.Sprintf("%d", newExpiration)
	return nil
}

func (controller DockerController) GetInstanceForOwner(instanceIdentifier string, ownerIdentifier string) (InstanceDetail, error) {
	if !strings.HasPrefix(instanceIdentifier, ownerIdentifier+"-") {
		return InstanceDetail{}, fmt.Errorf("instance not found")
	}

	status, err := controller.inspectStatus(instanceIdentifier)
	if err != nil {
		return InstanceDetail{}, err
	}

	return InstanceDetail{
		Identifier: instanceIdentifier,
		Status:     status,
	}, nil
}

func (controller DockerController) UpdateInstanceForOwner(
	instanceIdentifier string,
	ownerIdentifier string,
	update InstanceUpdateSpecification,
) error {
	if !strings.HasPrefix(instanceIdentifier, ownerIdentifier+"-") {
		return fmt.Errorf("instance not found")
	}

	status, statusError := controller.inspectStatus(instanceIdentifier)
	if statusError != nil {
		return statusError
	}

	if status == InstanceStatusRunning {
		stopError := controller.stopContainer(instanceIdentifier)
		if stopError != nil {
			return stopError
		}
	}

	return controller.startContainer(instanceIdentifier)
}

func (controller DockerController) DeleteInstance(instanceIdentifier string) error {
	cmd := exec.Command("docker", "rm", "-f", instanceIdentifier)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to delete container: %s", string(output))
	}

	return nil
}

func (controller DockerController) createContainer(instanceIdentifier string, image string) error {
	cmd := exec.Command(
		"docker",
		"run",
		"-d",
		"--name",
		instanceIdentifier,
		image,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create container: %s", string(output))
	}

	return nil
}

func (controller DockerController) inspectStatus(containerName string) (InstanceStatus, error) {
	cmd := exec.Command(
		"docker",
		"inspect",
		"-f",
		"{{.State.Running}}",
		containerName,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return InstanceStatusError, fmt.Errorf("container not found")
	}

	running := strings.TrimSpace(string(output)) == "true"

	if running {
		return InstanceStatusRunning, nil
	}

	return InstanceStatusStopped, nil
}

func (controller DockerController) stopContainer(containerName string) error {
	cmd := exec.Command("docker", "stop", containerName)
	_, err := cmd.CombinedOutput()
	return err
}

func (controller DockerController) startContainer(containerName string) error {
	cmd := exec.Command("docker", "start", containerName)
	_, err := cmd.CombinedOutput()
	return err
}
