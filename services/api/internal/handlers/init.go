package handlers

import (
	"common/pkg/app"
	"common/pkg/auth"
	"common/pkg/billing/account"
	"common/pkg/billing/payments/payment"
	"common/pkg/hardware"
	"common/pkg/listing"
	"common/pkg/money"
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
	UserAppService                 *app.UserAppService
	UserAuth                       *auth.AuthenticationService
	ListingService                 *app.PublicListingService
	BillingAccountService          *account.BillingAccountService
	PaymentService                 *payment.PaymentService
	InstanceEngineService          *app.InstanceSubscriptionService
	HardwareCostCalculationService *app.HardwareCostCalculationService
	InstanceDeployer               app.InstanceDeployer

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

type PublicListing struct {
	ID          string                          `json:"id"`
	Title       string                          `json:"title"`
	Description string                          `json:"description"`
	Hardware    *hardware.HardwareSpecification `json:"hardware,omitempty"`
	Price       *money.Money                    `json:"prices,omitempty"`
}

// ConvertListingToPublic converts a domain listing into a public-facing API listing.
func (app *App) ConvertListingToPublic(listingEntity *listing.Listing) *PublicListing {
	return &PublicListing{
		ID:          listingEntity.ID,
		Title:       listingEntity.Title,
		Description: listingEntity.Description,
		Hardware:    listingEntity.HardwareSpecification,
		Price:       &listingEntity.Price,
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
