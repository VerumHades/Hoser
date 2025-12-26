package instance

import (
	"common/internal/domain/instance"
	"common/internal/shared"
	"fmt"
	"time"
)

type DeployedInstanceStateEnum int

const (
	Running DeployedInstanceStateEnum = iota
	Building
	Stopped
	ErrorState
)

type DeployedInstanceState struct {
	State        DeployedInstanceStateEnum
	StateMessage string
	RuntimeError error
}

type InstanceDeployer interface {
	RefreshHardware(instanceID string, newSpec *shared.HardwareSpecification) error
	GetInstanceState(instanceID string) (DeployedInstanceState, error)
	StartInstance(instanceID string) error
	StopInstance(instanceID string) error
}

type InstanceSubscriptionReconciler struct {
	instanceSubscriptionService *InstanceSubscriptionService
	instanceRepository          *instance.InstanceRepository
	deployer                    InstanceDeployer
	checkInterval               time.Duration
	stopChan                    chan struct{}
}

func NewInstanceSubscriptionReconciler(engineService *InstanceSubscriptionService, instanceService *instance.InstanceRepository, deployer InstanceDeployer, interval time.Duration) *InstanceSubscriptionReconciler {
	return &InstanceSubscriptionReconciler{
		instanceSubscriptionService: engineService,
		instanceRepository:          instanceService,
		deployer:                    deployer,
		checkInterval:               interval,
		stopChan:                    make(chan struct{}),
	}
}

func (r *InstanceSubscriptionReconciler) Start() {
	go func() {
		ticker := time.NewTicker(r.checkInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := r.reconcile(); err != nil {
					fmt.Println("reconciliation error:", err)
				}
			case <-r.stopChan:
				return
			}
		}
	}()
}

func (r *InstanceSubscriptionReconciler) Stop() {
	close(r.stopChan)
}

func (r *InstanceSubscriptionReconciler) reconcile() error {
	now := time.Now()

	if err := r.handleExpiredInstances(now); err != nil {
		return fmt.Errorf("handling expired instances: %w", err)
	}

	if err := r.handlePendingUpdates(); err != nil {
		return fmt.Errorf("handling pending updates: %w", err)
	}

	if err := r.startActiveInstances(); err != nil {
		return fmt.Errorf("handling starts updates: %w", err)
	}

	if err := r.stopInactiveInstances(); err != nil {
		return fmt.Errorf("handling stops updates: %w", err)
	}

	return nil
}

// handleExpiredInstances checks instances past expiry and attempts renewal in batches.
func (r *InstanceSubscriptionReconciler) handleExpiredInstances(now time.Time) error {
	const batchSize = 50
	var lastSeenInstanceID shared.InstanceID

	for {
		expiredInstances, err := r.instanceRepository.FetchNextBatchExpiredBefore(now, lastSeenInstanceID, batchSize)
		if err != nil {
			return err
		}
		if len(expiredInstances) == 0 {
			break
		}

		for _, inst := range expiredInstances {
			if inst.ContractState != instance.ContractActive {
				continue
			}

			if !inst.RenewAutomatically {
				if _, err := r.instanceRepository.SetContractState(inst.ID, instance.ContractInactive); err != nil {
					return err
				}
				continue
			}

			_, err := r.instanceSubscriptionService.RenewInstanceHardware(inst.ID, inst.HardwareSpecification)
			if err != nil {
				r.setInstanceState(inst.ID, instance.ContractInactive, instance.BillingFailed)
				continue
			}
		}

		lastSeenInstanceID = expiredInstances[len(expiredInstances)-1].ID
	}

	return nil
}

// handlePendingUpdates applies any scheduled hardware updates in batches.
func (r *InstanceSubscriptionReconciler) handlePendingUpdates() error {
	const batchSize = 50
	var lastSeenInstanceID shared.InstanceID

	for {
		pendingInstances, err := r.instanceRepository.FetchNextBatchPendingHardwareUpdate(lastSeenInstanceID, batchSize)
		if err != nil {
			return err
		}
		if len(pendingInstances) == 0 {
			break
		}

		for _, inst := range pendingInstances {
			if inst.ContractState != instance.ContractActive {
				continue
			}
			if err := r.deployer.RefreshHardware(inst.ID, inst.HardwareSpecification); err != nil {
				return err
			}
		}

		lastSeenInstanceID = pendingInstances[len(pendingInstances)-1].ID
	}

	return nil
}

// startActiveInstances starts all active but non-running instances in batches.
func (r *InstanceSubscriptionReconciler) startActiveInstances() error {
	const batchSize = 50
	var lastSeenInstanceID shared.InstanceID

	for {
		activeNonRunning, err := r.instanceRepository.FetchNextBatchByState(instance.Running, instance.ContractActive, lastSeenInstanceID, batchSize)
		if err != nil {
			return fmt.Errorf("fetching active non-running instances: %w", err)
		}
		if len(activeNonRunning) == 0 {
			break
		}

		for _, inst := range activeNonRunning {
			if err := r.deployer.StartInstance(inst.ID); err != nil {
				_, _ = r.instanceRepository.SetState(inst.ID, instance.ErrorState)
				fmt.Printf("failed to start instance %s: %v\n", inst.ID, err)
				continue
			}
			_, err = r.instanceRepository.SetState(inst.ID, instance.Running)
			if err != nil {
				return err
			}
			fmt.Printf("started instance %s\n", inst.ID)
		}

		lastSeenInstanceID = activeNonRunning[len(activeNonRunning)-1].ID
	}

	return nil
}

// stopInactiveInstances stops all inactive but running instances in batches.
func (r *InstanceSubscriptionReconciler) stopInactiveInstances() error {
	const batchSize = 50
	var lastSeenInstanceID shared.InstanceID

	for {
		inactiveRunning, err := r.instanceRepository.FetchNextBatchByState(instance.Stopped, instance.ContractInactive, lastSeenInstanceID, batchSize)
		if err != nil {
			return fmt.Errorf("fetching inactive running instances: %w", err)
		}
		if len(inactiveRunning) == 0 {
			break
		}

		for _, inst := range inactiveRunning {
			if err := r.deployer.StopInstance(inst.ID); err != nil {
				_, _ = r.instanceRepository.SetState(inst.ID, instance.ErrorState)
				fmt.Printf("failed to stop instance %s: %v\n", inst.ID, err)
				continue
			}
			_, err = r.instanceRepository.SetState(inst.ID, instance.Stopped)
			if err != nil {
				return err
			}
			fmt.Printf("stopped instance %s\n", inst.ID)
		}

		lastSeenInstanceID = inactiveRunning[len(inactiveRunning)-1].ID
	}

	return nil
}
