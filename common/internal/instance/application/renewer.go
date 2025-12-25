package app

import (
	"common/pkg/hardware"
	"common/pkg/instance"
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
	RefreshHardware(instanceID string, newSpec *hardware.HardwareSpecification) error
	GetInstanceState(instanceID string) (DeployedInstanceState, error)
	StartInstance(instanceID string) error
	StopInstance(instanceID string) error
}

type InstanceSubscriptionReconciler struct {
	instanceSubscriptionService *InstanceSubscriptionService
	instanceService             *instance.InstanceService
	deployer                    InstanceDeployer
	checkInterval               time.Duration
	stopChan                    chan struct{}
}

func NewInstanceSubscriptionReconciler(engineService *InstanceSubscriptionService, instanceService *instance.InstanceService, deployer InstanceDeployer, interval time.Duration) *InstanceSubscriptionReconciler {
	return &InstanceSubscriptionReconciler{
		instanceSubscriptionService: engineService,
		instanceService:             instanceService,
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

// handleExpiredInstances checks instances past expiry and attempts renewal.
func (r *InstanceSubscriptionReconciler) handleExpiredInstances(now time.Time) error {
	expiredInstances, err := r.instanceService.ListExpiredInstances(now)
	if err != nil {
		return err
	}

	for _, inst := range expiredInstances {
		// skip instances without pending updates
		if inst.ContractState != instance.ContractActive {
			continue
		}

		if !inst.RenewAutomatically {
			_, err := r.instanceService.SetContractState(inst.ID, instance.ContractInactive)
			if err != nil {
				return err
			}
			return nil
		}

		_, err := r.instanceSubscriptionService.RenewInstanceHardware(
			inst.ID,
			inst.HardwareSpecification,
			inst.RenewalDuration,
		)

		if err != nil {
			r.setInstanceState(inst.ID, instance.ContractInactive, instance.BillingFailed)
			continue
		}
	}

	return nil
}
func (r *InstanceSubscriptionReconciler) setInstanceState(instanceID string, constractState instance.ContractState, state instance.InstanceState) error {
	_, err := r.instanceService.SetContractState(instanceID, constractState)
	if err != nil {
		return err
	}

	_, err = r.instanceService.SetState(instanceID, state)
	if err != nil {
		return err
	}

	return nil
}

// handlePendingUpdates applies any scheduled hardware updates.
func (r *InstanceSubscriptionReconciler) handlePendingUpdates() error {
	pendingInstances, err := r.instanceService.ListByPendingHardwareUpdate()
	if err != nil {
		return err
	}

	for _, inst := range pendingInstances {
		if inst.ContractState != instance.ContractActive {
			continue
		}

		if err := r.deployer.RefreshHardware(inst.ID, inst.HardwareSpecification); err != nil {
			return err
		}
	}

	return nil
}

func (r *InstanceSubscriptionReconciler) startActiveInstances() error {
	activeNonRunning, err := r.instanceService.ListActiveNonRunning()
	if err != nil {
		return fmt.Errorf("fetching active non-running instances: %w", err)
	}

	for _, inst := range activeNonRunning {
		if err := r.deployer.StartInstance(inst.ID); err != nil {
			_, err = r.instanceService.SetState(inst.ID, instance.ErrorState)
			fmt.Printf("failed to start instance %s: %v\n", inst.ID, err)
			continue
		}
		_, err = r.instanceService.SetState(inst.ID, instance.Running)
		if err != nil {
			return err
		}
		fmt.Printf("started instance %s\n", inst.ID)
	}

	return nil
}

func (r *InstanceSubscriptionReconciler) stopInactiveInstances() error {
	inactiveRunning, err := r.instanceService.ListInactiveRunning()
	if err != nil {
		return fmt.Errorf("fetching inactive running instances: %w", err)
	}

	for _, inst := range inactiveRunning {
		if err := r.deployer.StopInstance(inst.ID); err != nil {
			_, err = r.instanceService.SetState(inst.ID, instance.ErrorState)
			fmt.Printf("failed to stop instance %s: %v\n", inst.ID, err)
			continue
		}
		_, err = r.instanceService.SetState(inst.ID, instance.Stopped)
		if err != nil {
			return err
		}
		fmt.Printf("stopped instance %s\n", inst.ID)
	}

	return nil
}
