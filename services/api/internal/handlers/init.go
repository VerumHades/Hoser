package handlers

import (
	"net/http"
	"server/internal/configuration"
	"server/internal/database"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type App struct {
	DatabaseInteractor   database.Store
	RunningConfiguration *configuration.Configuration
	JWTSecret            []byte
}

// =================== APP INIT ===================

func NewApp(db database.Store, cfg *configuration.Configuration, jwtSecret []byte) *App {
	return &App{
		DatabaseInteractor:   db,
		RunningConfiguration: cfg,
		JWTSecret:            jwtSecret,
	}
}

// =================== HELPER FUNCTIONS ===================

// GetUserFromContext retrieves the authenticated user from Echo context (JWT claims)
func (app *App) GetUserFromContext(c echo.Context) (database.User, error) {
	cookie, err := c.Cookie("jwt")
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Missing token")
	}

	token, err := jwt.Parse(cookie.Value, func(t *jwt.Token) (interface{}, error) {
		return []byte(app.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, echo.NewHTTPError(http.StatusUnauthorized)
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, echo.NewHTTPError(http.StatusUnauthorized)
	}

	user, err := app.DatabaseInteractor.GetUserByID(userID)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusUnauthorized)
	}

	return user, nil
}

type HardwareDTO struct {
	CPU  int   `json:"cpu"`
	RAM  int64 `json:"ram"`
	Disk int64 `json:"disk"`
}

type PublicListing struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Hardware    *HardwareDTO      `json:"hardware,omitempty"`
	Prices      []ApiPricingEntry `json:"prices,omitempty"`
}

type CurrencyDTO struct {
	Name  string  `json:"name"`
	Short string  `json:"short"`
	Value float32 `json:"value"`
}

func (app *App) ConvertPricingListToAPI(pl database.PricingList) []ApiPricingEntry {
	var priceEntries []ApiPricingEntry
	if pl == nil {
		return priceEntries
	}

	list := pl.All() // assume All() returns []Pricing
	for i := range list {
		pricing := list[i]
		amount := pricing.Amount()

		priceEntries = append(priceEntries, ApiPricingEntry{
			ID:   pricing.UUID(),
			Type: int(pricing.Type()),
			Currency: CurrencyRequest{
				Name:  amount.Name(),
				Short: amount.Short(),
				Value: amount.Value(),
			},
		})
	}

	return priceEntries
}

func (app *App) ConvertListingToPublic(l database.Listing) *PublicListing {
	return &PublicListing{
		ID:          l.UUID(),
		Title:       l.Title(),
		Description: l.Description(),
		Hardware: &HardwareDTO{
			CPU:  l.HardwareRequirements().CPUCount(),
			RAM:  l.HardwareRequirements().RAMBytes(),
			Disk: l.HardwareRequirements().DiskBytes(),
		},
		Prices: app.ConvertPricingListToAPI(l.Pricing()),
	}
}

func (app *App) PublicListingsHandler(c echo.Context) error {
	options := database.ListingQueryOptions{}
	if q := c.QueryParam("q"); q != "" {
		options.Text = q
	}

	listings, err := app.DatabaseInteractor.QueryPublicListings(&options)
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

// =================== TEST APP ===================

func NewTestApp() *App {
	db := &database.DummyStore{}
	cfg := &configuration.Configuration{
		AllowedOrigins: []string{"http://localhost"},
	}
	jwtSecret := []byte("test-secret-key")

	return NewApp(db, cfg, jwtSecret)
}
