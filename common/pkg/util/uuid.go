package util

import "github.com/google/uuid"

// GenerateUUID returns a new UUID as a string.
func GenerateUUID() string {
	return uuid.New().String()
}
