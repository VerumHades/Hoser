package database

import (
	"errors"
	"strings"
)

// DummyInteractor implements Interactor with in-memory mock data.
type DummyInteractor struct{}

// Returns a dummy user if username is not empty
func (d *DummyInteractor) GetUserByName(username string) (*User, error) {
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}

	users := map[string]*User{
		"alice": {
			ID:           1,
			Username:     "alice",
			PasswordHash: "$2y$10$LHg6JpCH8z5YHKsdYBAtiukBbWLLuR.Gyowp7KzX8R3rOm/hjHPJa",
			IsDeveloper:  true,
		},
		"bob": {
			ID:           2,
			Username:     "bob",
			PasswordHash: "$2y$10$LHg6JpCH8z5YHKsdYBAtiukBbWLLuR.Gyowp7KzX8R3rOm/hjHPJa",
			IsDeveloper:  false,
		},
		"carol": {
			ID:           3,
			Username:     "carol",
			PasswordHash: "$2y$10$LHg6JpCH8z5YHKsdYBAtiukBbWLLuR.Gyowp7KzX8R3rOm/hjHPJa",
			IsDeveloper:  true,
		},
	}

	user, ok := users[username]
	if !ok {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// Returns dummy listings (software packages a user owns or maintains)
func (d *DummyInteractor) GetUserListings(user *User) ([]*Listing, error) {
	if user == nil {
		return nil, errors.New("user cannot be nil")
	}

	listings := map[int][]*Listing{
		1: {
			{ID: 1, Title: "Node.js Dev Stack", Description: "Node.js 20 with MongoDB, Redis, and Nginx configured for local testing."},
			{ID: 2, Title: "Go Microservices Boilerplate", Description: "Multi-container setup with Go, PostgreSQL, and Kafka for event-driven apps."},
		},
		2: {
			{ID: 3, Title: "Python ML Environment", Description: "Preconfigured JupyterLab with TensorFlow, PyTorch, and scikit-learn in Docker."},
		},
		3: {
			{ID: 4, Title: "Rust Web API Template", Description: "Includes Rocket framework, Postgres, and Docker Compose setup."},
			{ID: 5, Title: "React + Vite Frontend", Description: "Modern React frontend dev environment with hot reload."},
		},
	}

	return listings[user.ID], nil
}

// Returns dummy rentals (software environments the user has currently deployed)
func (d *DummyInteractor) GetUserRentals(user *User) ([]*Rental, error) {
	if user == nil {
		return nil, errors.New("user cannot be nil")
	}

	rentals := map[int][]*Rental{
		1: {
			{ID: 11, Title: "Go API Testbed", Description: "Running Docker container: golang:1.22 on port 8080"},
			{ID: 12, Title: "Redis Cache Instance", Description: "Redis 7 container for caching layer testing"},
		},
		2: {
			{ID: 13, Title: "TensorFlow Training Node", Description: "GPU-enabled container with TF 2.15"},
		},
		3: {
			{ID: 14, Title: "Rust Playground", Description: "Isolated dev container with Cargo and Diesel preinstalled"},
			{ID: 15, Title: "Frontend Preview Server", Description: "Nginx + Node 20 container serving build artifacts"},
		},
	}

	return rentals[user.ID], nil
}

// Returns dummy public software rentals (shared or open environments)
func (d *DummyInteractor) QueryPublicRentals(options *RentalQueryOptions) ([]*Rental, error) {
	rentals := []*Rental{
		{ID: 101, Title: "Public PostgreSQL Test DB", Description: "Shared containerized Postgres 16 instance for integration testing."},
		{ID: 102, Title: "Docker-in-Docker Sandbox", Description: "Run isolated containers inside a shared Docker environment."},
		{ID: 103, Title: "Public Nginx Reverse Proxy", Description: "Shared Nginx setup for proxy and routing experiments."},
		{ID: 104, Title: "Open Source CI/CD Runner", Description: "GitHub Actions compatible runner in Docker, public usage enabled."},
	}

	// If no text query, return all rentals
	if options == nil || options.Text == "" {
		return rentals, nil
	}

	query := strings.ToLower(options.Text)
	var filtered []*Rental

	for _, r := range rentals {
		if strings.Contains(strings.ToLower(r.Title), query) || strings.Contains(strings.ToLower(r.Description), query) {
			filtered = append(filtered, r)
		}
	}

	return filtered, nil
}
