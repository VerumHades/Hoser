package handlers

import (
	"net/http"
	"time"

	"common/pkg/app"
	"common/pkg/hardware"
	"common/pkg/instance"

	"github.com/labstack/echo/v4"
)

type ApiInstance struct {
	ID                    string                          `json:"id"`
	ListingID             string                          `json:"listingId"`
	BillingID             string                          `json:"billingId"`
	State                 string                          `json:"state"`
	ContractState         string                          `json:"contractState"`
	HardwareSpecification *hardware.HardwareSpecification `json:"hardwareSpecification"`
	DesiredHardwareSpec   *hardware.HardwareSpecification `json:"desiredHardwareSpecification,omitempty"`
	Expiry                time.Time                       `json:"expiry"`
}

// LaunchInstanceRequest is the payload to create a new instance
type LaunchInstanceRequest struct {
	ListingID        string                          `json:"listingId"`
	BillingAccountID string                          `json:"billingAccountId"`
	HardwareSpec     *hardware.HardwareSpecification `json:"hardwareSpecification"`
}

// UpdateHardwareRequest is the payload to update hardware spec
type UpdateHardwareRequest struct {
	HardwareSpec *hardware.HardwareSpecification `json:"hardwareSpecification"`
}

/*
*
InstanceStateToAPI converts an internal instance state to its API representation.
*/
func InstanceStateToAPI(instanceState instance.InstanceState) string {
	switch instanceState {
	case instance.Running:
		return "RUNNING"
	case instance.Stopped:
		return "STOPPED"
	case instance.ErrorState:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

func ContractStateToAPI(contractState instance.ContractState) string {
	switch contractState {
	case instance.ContractInactive:
		return "INACTIVE"
	case instance.ContractActive:
		return "ACTIVE"
	default:
		return "UNKNOWN"
	}
}

func (a *App) MakeApiInstance(domainInstance *instance.Instance) ApiInstance {
	return ApiInstance{
		ID:                    domainInstance.ID,
		ListingID:             domainInstance.ListingID,
		BillingID:             domainInstance.BillingID,
		State:                 InstanceStateToAPI(domainInstance.State),
		ContractState:         ContractStateToAPI(domainInstance.ContractState),
		HardwareSpecification: domainInstance.HardwareSpecification,
		Expiry:                domainInstance.Expiry,
	}
}

// =================== HANDLERS ===================

// UserLaunchInstanceHandler launches a new instance for the authenticated user
func (a *App) UserLaunchInstanceHandler(c echo.Context) error {
	userID, err := a.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	var req LaunchInstanceRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON")
	}

	newInstance, err := a.InstanceEngineService.LaunchInstance(
		userID,
		req.ListingID,
		req.BillingAccountID,
		req.HardwareSpec,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to launch instance: "+err.Error())
	}

	newInstance, err = a.InstanceEngineService.RenewInstanceHardware(newInstance.ID, req.HardwareSpec, time.Hour*24*30)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to purchase hardware: "+err.Error())
	}

	return c.JSON(http.StatusOK, a.MakeApiInstance(newInstance))
}

// UserUpdateHardwareHandler updates hardware spec for an instance
func (a *App) UserUpdateHardwareHandler(c echo.Context) error {
	userID, err := a.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	instanceID := c.Param("id")
	var req UpdateHardwareRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON")
	}

	instanceObj, err := a.InstanceEngineService.GetInstance(instanceID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch instance")
	}
	if instanceObj == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Instance not found")
	}

	// Verify ownership by checking billing account owner
	billingAccount, err := a.BillingAccountService.GetAccount(instanceObj.BillingID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch billing account")
	}
	if billingAccount.OwnerID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "Not authorized to modify this instance")
	}

	updatedInstance, err := a.InstanceEngineService.RenewInstanceHardware(instanceID, req.HardwareSpec, time.Hour*24*30)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update hardware: "+err.Error())
	}

	return c.JSON(http.StatusOK, a.MakeApiInstance(updatedInstance))
}

// UserListInstancesByBillingHandler lists all instances for a billing account
func (a *App) UserListInstancesByBillingHandler(c echo.Context) error {
	userID, err := a.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	billingID := c.Param("billingId")
	billingAccount, err := a.BillingAccountService.GetAccount(billingID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch billing account")
	}
	if billingAccount.OwnerID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "Not authorized")
	}

	instances, err := a.InstanceEngineService.ListInstancesByBillingAccount(billingID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list instances")
	}

	apiInstances := make([]ApiInstance, len(instances))
	for i, inst := range instances {
		apiInstances[i] = a.MakeApiInstance(inst)
	}

	return c.JSON(http.StatusOK, apiInstances)
}

// UserListInstancesByOwnerHandler lists all instances for the authenticated user
func (a *App) UserListInstancesByOwnerHandler(c echo.Context) error {
	userID, err := a.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	instances, err := a.InstanceEngineService.ListInstancesByOwner(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list instances")
	}

	apiInstances := make([]ApiInstance, len(instances))
	for i, inst := range instances {
		apiInstances[i] = a.MakeApiInstance(inst)
	}

	return c.JSON(http.StatusOK, apiInstances)
}

// UserGetInstanceHandler retrieves a single instance by ID
func (a *App) UserGetInstanceHandler(c echo.Context) error {
	userID, err := a.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	instanceID := c.Param("id")
	instanceObj, err := a.InstanceEngineService.GetInstance(instanceID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch instance")
	}
	if instanceObj == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Instance not found")
	}

	billingAccount, err := a.BillingAccountService.GetAccount(instanceObj.BillingID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch billing account")
	}
	if billingAccount.OwnerID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "Not authorized")
	}

	return c.JSON(http.StatusOK, a.MakeApiInstance(instanceObj))
}

// UserGetInstanceStateHandler retrieves the deployment state of an instance.
func (a *App) UserGetInstanceStateHandler(c echo.Context) error {
	userID, err := a.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	instanceID := c.Param("id")
	instanceObj, err := a.InstanceEngineService.GetInstance(instanceID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch instance")
	}
	if instanceObj == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Instance not found")
	}

	billingAccount, err := a.BillingAccountService.GetAccount(instanceObj.BillingID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch billing account")
	}
	if billingAccount.OwnerID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "Not authorized")
	}

	// Fetch the instance state from the deployer
	instanceState, err := a.InstanceDeployer.GetInstanceState(instanceID)
	if err != nil {
		return c.NoContent(http.StatusNoContent)
	}

	// Map Go enum to string
	stateStr := "ErrorState"
	switch instanceState.State {
	case app.Running:
		stateStr = "Running"
	case app.Building:
		stateStr = "Building"
	case app.Stopped:
		stateStr = "Stopped"
	case app.ErrorState:
		stateStr = "ErrorState"
	}

	// Build the response matching the TypeScript interface
	response := map[string]interface{}{
		"instance_id":  instanceID,
		"state":        stateStr,
		"stateMessage": instanceState.StateMessage,
	}

	if instanceState.RuntimeError != nil {
		response["error"] = instanceState.RuntimeError.Error()
	}

	return c.JSON(http.StatusOK, response)
}
