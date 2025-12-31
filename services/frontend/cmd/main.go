package main

import (
	"common/pkg/infrastructure/configuration"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	config, err := configuration.Load[configuration.Configuration]()
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	// ---------------------------
	// CORS middleware
	// ---------------------------
	if len(config.AllowedOrigins) > 0 && config.AllowedOrigins[0] != "" {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: config.AllowedOrigins,
			AllowMethods: []string{echo.GET, echo.HEAD, echo.OPTIONS},
		}))
	}

	// ---------------------------
	// Serve frontend SPA
	// ---------------------------
	fsRoot := config.ClientDirectory
	e.GET("/*", func(c echo.Context) error {
		path := filepath.Join(fsRoot, c.Request().URL.Path)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			// Fallback to index.html
			return c.File(filepath.Join(fsRoot, "index.html"))
		}
		return c.File(path)
	})

	address := fmt.Sprintf("%s:%s", config.Address, config.Port)
	fmt.Printf("Frontend server running on http://%s\n", address)
	e.Logger.Fatal(e.Start(address))
}
