package database

import (
	"errors"
	"strconv"
	"strings"
)

// -----------------------
// Types
// -----------------------

type DummyUser struct {
	ID           int
	Username     string
	PasswordHash string
	IsDev        bool
}

func (u *DummyUser) GetUsername() string     { return u.Username }
func (u *DummyUser) GetPasswordHash() string { return u.PasswordHash }
func (u *DummyUser) IsDeveloper() bool       { return u.IsDev }
func (u *DummyUser) GetListings() ([]Listing, error) {
	indices, exists := userListingMap[u.ID]
	if !exists {
		return nil, nil
	}
	var listings []Listing
	for _, idx := range indices {
		listings = append(listings, globalListings[idx])
	}
	return listings, nil
}
func (u *DummyUser) GetRentals() ([]Rental, error) { return dummyRentals[u.ID], nil }

func (u *DummyUser) CreateListing(title string, description string) (Listing, error) {
	listing := &DummyListing{
		ID:          len(globalListings), // new ID = next index
		Title:       title,
		Description: description,
	}
	globalListings = append(globalListings, listing)
	userListingMap[u.ID] = append(userListingMap[u.ID], listing.ID)
	return listing, nil
}

func (u *DummyUser) GetListing(uuid string) (Listing, error) {
	index, err := uuidToIndex(uuid)
	if err != nil {
		return nil, err
	}

	if index < 0 || index >= len(globalListings) {
		return nil, errors.New("invalid id")
	}

	// Ensure the user owns this listing
	indices := userListingMap[u.ID]
	for _, i := range indices {
		if i == index {
			return globalListings[i], nil
		}
	}
	return nil, errors.New("listing not owned by user")
}

func (u *DummyUser) DeleteListing(uuid string) (Listing, error) {
	index, err := uuidToIndex(uuid)
	if err != nil {
		return nil, err
	}

	// Ensure index is valid
	if index < 0 || index >= len(globalListings) {
		return nil, errors.New("invalid listing ID")
	}

	// Ensure the user owns this listing
	ownedIndices, exists := userListingMap[u.ID]
	if !exists {
		return nil, errors.New("user has no listings")
	}

	owned := false
	newOwnedIndices := []int{}
	for _, i := range ownedIndices {
		if i == index {
			owned = true
		} else {
			newOwnedIndices = append(newOwnedIndices, i)
		}
	}

	if !owned {
		return nil, errors.New("listing not owned by user")
	}

	// Remove from user map
	userListingMap[u.ID] = newOwnedIndices

	// Optionally, remove from globalListings by setting it to nil to preserve indices
	deletedListing := globalListings[index]
	globalListings[index] = nil

	return deletedListing, nil
}

type DummyListing struct {
	ID          int
	Title       string
	Description string
}

func (l *DummyListing) GetTitle() string            { return l.Title }
func (l *DummyListing) GetDescription() string      { return l.Description }
func (l *DummyListing) SetTitle(title string) error { l.Title = title; return nil }
func (l *DummyListing) SetDescription(description string) error {
	l.Description = description
	return nil
}
func (l *DummyListing) GetUUID() string { return strconv.Itoa(l.ID) }

type DummyRental struct {
	ID          int
	Title       string
	Description string
}

func (r *DummyRental) GetTitle() string       { return r.Title }
func (r *DummyRental) GetDescription() string { return r.Description }

// -----------------------
// Dummy data
// -----------------------

var dummyUsers = map[string]*DummyUser{
	"alice": {ID: 1, Username: "alice", PasswordHash: "$2y$10$yyty1BZZiACa4rMHC/ksfOd3fUnwxS5skZAo3Fo6iTCx0bxfZIcMS", IsDev: true},
	"bob":   {ID: 2, Username: "bob", PasswordHash: "$2y$10$yyty1BZZiACa4rMHC/ksfOd3fUnwxS5skZAo3Fo6iTCx0bxfZIcMS", IsDev: false},
	"carol": {ID: 3, Username: "carol", PasswordHash: "$2y$10$yyty1BZZiACa4rMHC/ksfOd3fUnwxS5skZAo3Fo6iTCx0bxfZIcMS", IsDev: true},
}

// Global listing storage
var globalListings = []Listing{
	&DummyListing{ID: 0, Title: "Node.js Dev Stack", Description: "Node.js 20 with MongoDB..."},
	&DummyListing{ID: 1, Title: "Go Microservices Boilerplate", Description: "Go + Kafka setup"},
	&DummyListing{ID: 2, Title: "Python ML Environment", Description: "TensorFlow + PyTorch"},
}

// Map from user ID → indices of listings in globalListings
var userListingMap = map[int][]int{
	1: {0, 1}, // alice owns listings 0 and 1
	2: {2},    // bob owns listing 2
}

// Rentals (unchanged)
var dummyRentals = map[int][]Rental{
	1: {&DummyRental{ID: 11, Title: "Go API Testbed", Description: "Running Docker container"}},
	2: {&DummyRental{ID: 13, Title: "TensorFlow Training Node", Description: "GPU container"}},
}

// -----------------------
// Interactor
// -----------------------

type DummyInteractor struct{}

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

func uuidToIndex(uuid string) (int, error) {
	n, err := strconv.Atoi(uuid)
	if err != nil {
		return 0, errors.New("conversion failed")
	}
	return n, nil
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
