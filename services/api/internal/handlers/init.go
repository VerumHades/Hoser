package handlers

import (
	"api/internal/database"
	"common/pkg/app"
	"common/pkg/auth"
	"common/pkg/configuration"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type App struct {
	UserAppService *app.UserAppService
	UserAuth       *auth.AuthenticationService

	RunningConfiguration *configuration.Configuration

	JWTSecret []byte
}

// =================== HELPER FUNCTIONS ===================

// GetUserFromContext retrieves the authenticated user from Echo context (JWT claims)
func (app *App) GetUserIDFromContext(c echo.Context) (string, error) {
	cookie, err := c.Cookie("jwt")
	if err != nil {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "Missing token")
	}

	token, err := jwt.Parse(cookie.Value, func(t *jwt.Token) (interface{}, error) {
		return []byte(app.JWTSecret), nil
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

type HardwareDTO struct {
	CPU  int   `json:"cpu"`
	RAM  int64 `json:"ram"`
	Disk int64 `json:"disk"`
}

type PublicListing struct {
	ID          string         `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Hardware    *HardwareDTO   `json:"hardware,omitempty"`
	Prices      []ListingPrice `json:"prices,omitempty"`
}

type CurrencyDTO struct {
	Name  string  `json:"name"`
	Short string  `json:"short"`
	Value float32 `json:"value"`
}

func (app *App) ConvertPricingListToAPI(pl database.PricingList) []ListingPrice {
	var priceEntries []ListingPrice
	if pl == nil {
		return priceEntries
	}

	list := pl.All() // assume All() returns []Pricing
	for i := range list {
		pricing := list[i]
		amount := pricing.Amount()

		priceEntries = append(priceEntries, ListingPrice{
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
