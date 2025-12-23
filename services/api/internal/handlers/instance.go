package handlers

import (
	"net/http"
	"time"

	"common/pkg/hardware"
	"common/pkg/instance"

	"github.com/labstack/echo/v4"
)

// =================== TYPES ===================
type ApiInstance struct {
	ID                    string                          `json:"id"`
	ListingID             string                          `json:"listingId"`
	BillingID             string                          `json:"billingId"`
	State                 string                          `json:"state"`
	HardwareSpecification *hardware.HardwareSpecification `json:"hardwareSpecification"`
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

// =================== HELPERS ===================
func (app *App) MakeApiInstance(instance *instance.Instance) ApiInstance {
	return ApiInstance{
		ID:                    instance.ID,
		ListingID:             instance.ListingID,
		BillingID:             instance.BillingID,
		State:                 "",
		HardwareSpecification: instance.HardwareSpecification,
	}
}

// =================== HANDLERS ===================

// UserLaunchInstanceHandler launches a new instance for the authenticated user
func (app *App) UserLaunchInstanceHandler(c echo.Context) error {
	userID, err := app.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	var req LaunchInstanceRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON")
	}

	newInstance, err := app.InstanceEngineService.LaunchInstance(
		userID,
		req.ListingID,
		req.BillingAccountID,
		req.HardwareSpec,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to launch instance: "+err.Error())
	}

	newInstance, err = app.InstanceEngineService.RenewInstanceHardware(newInstance.ID, req.HardwareSpec, time.Hour*24*30)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to purchase hardware: "+err.Error())
	}

	return c.JSON(http.StatusOK, app.MakeApiInstance(newInstance))
}

// UserUpdateHardwareHandler updates hardware spec for an instance
func (app *App) UserUpdateHardwareHandler(c echo.Context) error {
	userID, err := app.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	instanceID := c.Param("id")
	var req UpdateHardwareRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON")
	}

	instanceObj, err := app.InstanceEngineService.GetInstance(instanceID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch instance")
	}
	if instanceObj == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Instance not found")
	}

	// Verify ownership by checking billing account owner
	billingAccount, err := app.BillingAccountService.GetAccount(instanceObj.BillingID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch billing account")
	}
	if billingAccount.OwnerID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "Not authorized to modify this instance")
	}

	updatedInstance, err := app.InstanceEngineService.RenewInstanceHardware(instanceID, req.HardwareSpec, time.Hour*24*30)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update hardware: "+err.Error())
	}

	return c.JSON(http.StatusOK, app.MakeApiInstance(updatedInstance))
}

// UserListInstancesByBillingHandler lists all instances for a billing account
func (app *App) UserListInstancesByBillingHandler(c echo.Context) error {
	userID, err := app.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	billingID := c.Param("billingId")
	billingAccount, err := app.BillingAccountService.GetAccount(billingID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch billing account")
	}
	if billingAccount.OwnerID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "Not authorized")
	}

	instances, err := app.InstanceEngineService.ListInstancesByBillingAccount(billingID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list instances")
	}

	apiInstances := make([]ApiInstance, len(instances))
	for i, inst := range instances {
		apiInstances[i] = app.MakeApiInstance(inst)
	}

	return c.JSON(http.StatusOK, apiInstances)
}

// UserListInstancesByOwnerHandler lists all instances for the authenticated user
func (app *App) UserListInstancesByOwnerHandler(c echo.Context) error {
	userID, err := app.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	instances, err := app.InstanceEngineService.ListInstancesByOwner(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list instances")
	}

	apiInstances := make([]ApiInstance, len(instances))
	for i, inst := range instances {
		apiInstances[i] = app.MakeApiInstance(inst)
	}

	return c.JSON(http.StatusOK, apiInstances)
}

// UserGetInstanceHandler retrieves a single instance by ID
func (app *App) UserGetInstanceHandler(c echo.Context) error {
	userID, err := app.GetUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	instanceID := c.Param("id")
	instanceObj, err := app.InstanceEngineService.GetInstance(instanceID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch instance")
	}
	if instanceObj == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Instance not found")
	}

	billingAccount, err := app.BillingAccountService.GetAccount(instanceObj.BillingID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch billing account")
	}
	if billingAccount.OwnerID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "Not authorized")
	}

	return c.JSON(http.StatusOK, app.MakeApiInstance(instanceObj))
}
