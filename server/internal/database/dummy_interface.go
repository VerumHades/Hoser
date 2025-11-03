package database

import (
	"errors"
	"strings"
)

// Concrete struct types implementing the interfaces

type DummyUser struct {
	ID           int
	Username     string
	PasswordHash string
	IsDev        bool
}

func (u *DummyUser) GetUsername() string             { return u.Username }
func (u *DummyUser) GetPasswordHash() string         { return u.PasswordHash }
func (u *DummyUser) IsDeveloper() bool               { return u.IsDev }
func (u *DummyUser) GetListings() ([]Listing, error) { return dummyListings[u.ID], nil }
func (u *DummyUser) GetRentals() ([]Rental, error)   { return dummyRentals[u.ID], nil }

type DummyListing struct {
	ID          int
	Title       string
	Description string
}

func (u *DummyUser) CreateListing(title string, description string) (Listing, error) {
	listing := &DummyListing{
		ID:          len(dummyListings[u.ID]) + 1,
		Title:       title,
		Description: description,
	}
	dummyListings[u.ID] = append(dummyListings[u.ID], listing)
	return dummyListings[u.ID][len(dummyListings[u.ID])-1], nil
}
func (l *DummyListing) GetTitle() string       { return l.Title }
func (l *DummyListing) GetDescription() string { return l.Description }

type DummyRental struct {
	ID          int
	Title       string
	Description string
}

func (r *DummyRental) GetTitle() string       { return r.Title }
func (r *DummyRental) GetDescription() string { return r.Description }

// DummyInteractor implements Interactor with mock data
type DummyInteractor struct{}

var dummyUsers = map[string]*DummyUser{
	"alice": {ID: 1, Username: "alice", PasswordHash: "$2y$10$LHg6JpCH8z5YHKsdYBAtiukBbWLLuR.Gyowp7KzX8R3rOm/hjHPJa", IsDev: true},
	"bob":   {ID: 2, Username: "bob", PasswordHash: "$2y$10$LHg6JpCH8z5YHKsdYBAtiukBbWLLuR.Gyowp7KzX8R3rOm/hjHPJa", IsDev: false},
	"carol": {ID: 3, Username: "carol", PasswordHash: "$2y$10$LHg6JpCH8z5YHKsdYBAtiukBbWLLuR.Gyowp7KzX8R3rOm/hjHPJa", IsDev: true},
}

var dummyListings = map[int][]Listing{
	1: {
		&DummyListing{ID: 1, Title: "Node.js Dev Stack", Description: "Node.js 20 with MongoDB..."},
		&DummyListing{ID: 2, Title: "Go Microservices Boilerplate", Description: "Go + Kafka setup"},
	},
	2: {
		&DummyListing{ID: 3, Title: "Python ML Environment", Description: "TensorFlow + PyTorch"},
	},
}

var dummyRentals = map[int][]Rental{
	1: {
		&DummyRental{ID: 11, Title: "Go API Testbed", Description: "Running Docker container"},
	},
	2: {
		&DummyRental{ID: 13, Title: "TensorFlow Training Node", Description: "GPU container"},
	},
}

// Implementation of Interactor interface

func (d *DummyInteractor) GetUserByName(username string) (User, error) {
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}
	u, ok := dummyUsers[username]
	if !ok {
		return nil, errors.New("user not found")
	}
	return u, nil
}

func (d *DummyInteractor) QueryPublicRentals(options *RentalQueryOptions) ([]Rental, error) {
	rentals := []Rental{
		&DummyRental{ID: 101, Title: "Public PostgreSQL", Description: "Shared test DB"},
		&DummyRental{ID: 102, Title: "Docker-in-Docker Sandbox", Description: "Run containers inside containers"},
	}

	if options == nil || options.Text == "" {
		return rentals, nil
	}

	query := strings.ToLower(options.Text)
	var filtered []Rental
	for _, r := range rentals {
		if strings.Contains(strings.ToLower(r.GetTitle()), query) ||
			strings.Contains(strings.ToLower(r.GetDescription()), query) {
			filtered = append(filtered, r)
		}
	}
	return filtered, nil
}
