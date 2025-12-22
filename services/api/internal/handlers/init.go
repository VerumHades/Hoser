package handlers

import (
	"common/pkg/app"
	"common/pkg/auth"
	"common/pkg/billing/account"
	"common/pkg/billing/payments/payment"
	"common/pkg/hardware"
	"common/pkg/listing"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type Configuration struct {
	Port              string   `env:"PORT" default:"8080"`
	Address           string   `env:"ADDRESS" default:"localhost"`
	AllowedOrigins    []string `env:"ALLOWED_ORIGINS" separator:","`
	JWTSecret         string   `env:"JWT_SECRET" default:"SECRET"`
	InterserverSecret string   `env:"INTERSERVER_SECRET" default:"SECRET"`

	MongoDatabaseName     string `env:"MONGO_DATABASE" default:"apiDatabase"`
	MongoDatabaseUserName string `env:"MONGO_USER" default:"apiUser"`
	MongoDatabasePassword string `env:"MONGO_PASSWORD" default:"apiUserPassword"`
	MongoDatabaseHost     string `env:"MONGO_HOST" default:"localhost"`
	MongoDatabasePort     string `env:"MONGO_PORT" default:"27017"`
}

type App struct {
	UserAppService        *app.UserAppService
	UserAuth              *auth.AuthenticationService
	ListingService        *app.PublicListingService
	BillingAccountService *account.BillingAccountService
	PaymentService        *payment.PaymentService
	InstanceEngineService *app.InstanceEngineService

	RunningConfiguration *Configuration
}

// =================== HELPER FUNCTIONS ===================

// GetUserFromContext retrieves the authenticated user from Echo context (JWT claims)
func (app *App) GetUserIDFromContext(c echo.Context) (string, error) {
	cookie, err := c.Cookie("jwt")
	if err != nil {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "Missing token")
	}

	token, err := jwt.Parse(cookie.Value, func(t *jwt.Token) (interface{}, error) {
		return []byte(app.RunningConfiguration.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", echo.NewHTTPError(http.StatusUnauthorized)
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", echo.NewHTTPError(http.StatusUnauthorized)
	}

	return userID, nil
}

type HardwareSpecificationDTO struct {
	CPU  *int   `json:"cpu,omitempty"`
	RAM  *int64 `json:"ramBytes,omitempty"`
	Disk *int64 `json:"diskBytes,omitempty"`
}

func (app *App) HardwareSpecificationToDTO(spec *hardware.HardwareSpecification) *HardwareSpecificationDTO {
	if spec == nil {
		return nil
	}
	return &HardwareSpecificationDTO{
		CPU:  &spec.CPUCount,
		RAM:  &spec.RAMBytes,
		Disk: &spec.DiskBytes,
	}
}

func (application *App) HardwareDTOToSpec(dto *HardwareSpecificationDTO) *hardware.HardwareSpecification {
	if dto == nil {
		return nil
	}

	hardwareSpecification := &hardware.HardwareSpecification{}

	if dto.CPU != nil {
		hardwareSpecification.CPUCount = *dto.CPU
	}

	if dto.RAM != nil {
		hardwareSpecification.RAMBytes = *dto.RAM
	}

	if dto.Disk != nil {
		hardwareSpecification.DiskBytes = *dto.Disk
	}

	return hardwareSpecification
}

type CurrencyRequest struct {
	Value float32 `json:"value"`
	Name  string  `json:"name"`
	Short string  `json:"short"`
}

type PublicListing struct {
	ID          string                    `json:"id"`
	Title       string                    `json:"title"`
	Description string                    `json:"description"`
	Hardware    *HardwareSpecificationDTO `json:"hardware,omitempty"`
	Price       *CurrencyRequest          `json:"prices,omitempty"`
}

// ConvertListingToPublic converts a domain listing into a public-facing API listing.
func (app *App) ConvertListingToPublic(listingEntity *listing.Listing) *PublicListing {
	return &PublicListing{
		ID:          listingEntity.ID,
		Title:       listingEntity.Title,
		Description: listingEntity.Description,
		Hardware:    app.HardwareSpecificationToDTO(listingEntity.HardwareSpecification),
		Price: &CurrencyRequest{
			Value: float32(listingEntity.Price.Amount),
			Short: listingEntity.Price.CurrencyCode,
		},
	}
}

type ListingQueryOptions struct {
	Text string
}

func (app *App) PublicListingsHandler(c echo.Context) error {
	options := ListingQueryOptions{}
	if q := c.QueryParam("q"); q != "" {
		options.Text = q
	}

	listings, err := app.ListingService.SearchPublicListings(options.Text)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	// Convert to API response DTOs
	response := make([]*PublicListing, 0, len(listings))
	for _, l := range listings {
		response = append(response, app.ConvertListingToPublic(l))
	}

	return c.JSON(http.StatusOK, response)
}

func (app *App) PublicListingHandler(c echo.Context) error {
	listingID := c.Param("id")
	if listingID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing listing ID"})
	}

	listingEntity, err := app.ListingService.GetPublicListing(listingID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "listing not found"})
	}

	publicListing := app.ConvertListingToPublic(listingEntity)
	return c.JSON(http.StatusOK, publicListing)
}
