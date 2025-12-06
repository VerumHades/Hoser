package configuration

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Configuration struct {
	Port            string
	Address         string
	ClientDirectory string
	AllowedOrigins  []string
}

func Load() Configuration {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on system environment.")
	}

	return Configuration{
		Port:            getEnviromentalValueOrDefault("PORT", "8080"),
		Address:         getEnviromentalValueOrDefault("ADDRESS", "localhost"),
		ClientDirectory: getEnviromentalValueOrDefault("CLIENT_DIRECTORY", "../client/dist"),
		AllowedOrigins:  strings.Split(getEnviromentalValueOrDefault("ALLOWED_ORIGINS", ""), ","),
	}
}

func getEnviromentalValueOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		fmt.Printf("%s -> %s\n", key, value)
		return value
	}
	return fallback
}
