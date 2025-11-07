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
		if idx >= 0 && idx < len(globalListings) && globalListings[idx] != nil {
			listings = append(listings, globalListings[idx])
		}
	}
	return listings, nil
}

func (u *DummyUser) GetRentals() ([]Rental, error) { return dummyRentals[u.ID], nil }

func (u *DummyUser) CreateListing(title string, description string) (Listing, error) {
	listing := &DummyListing{
		ID:          len(globalListings),
		Title:       title,
		Description: description,
		AccessMode:  Private, // default new listings to private
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
	if index < 0 || index >= len(globalListings) || globalListings[index] == nil {
		return nil, errors.New("invalid id")
	}
	for _, owned := range userListingMap[u.ID] {
		if owned == index {
			return globalListings[index], nil
		}
	}
	return nil, errors.New("listing not owned by user")
}

func (u *DummyUser) DeleteListing(uuid string) (Listing, error) {
	index, err := uuidToIndex(uuid)
	if err != nil {
		return nil, err
	}
	if index < 0 || index >= len(globalListings) || globalListings[index] == nil {
		return nil, errors.New("invalid listing ID")
	}
	ownedIndices := userListingMap[u.ID]
	newOwned := []int{}
	found := false
	for _, i := range ownedIndices {
		if i == index {
			found = true
		} else {
			newOwned = append(newOwned, i)
		}
	}
	if !found {
		return nil, errors.New("listing not owned by user")
	}
	userListingMap[u.ID] = newOwned
	deleted := globalListings[index]
	globalListings[index] = nil
	return deleted, nil
}

func (u *DummyUser) Rent(listingUUID string) error {
	// Convert UUID to index
	index, err := uuidToIndex(listingUUID)
	if err != nil {
		return errors.New("invalid listing UUID")
	}

	// Check bounds
	if index < 0 || index >= len(globalListings) || globalListings[index] == nil {
		return errors.New("listing not found")
	}

	listing := globalListings[index]

	// Only allow renting public listings
	if listing.GetAccessMode() != Public {
		return errors.New("listing is not public")
	}

	// Create a new rental
	rental := &DummyRental{
		ID:                len(dummyRentals[u.ID]) + 100, // unique dummy ID
		SourceListingUUID: listing.GetUUID(),
		Title:             listing.GetTitle(),
		Description:       listing.GetDescription(),
	}

	dummyRentals[u.ID] = append(dummyRentals[u.ID], rental)

	return nil
}

// -----------------------
// Listing
// -----------------------

type DummyListing struct {
	ID          int
	Title       string
	Description string
	AccessMode  ListingAccessMode
}

func (l *DummyListing) GetUUID() string                  { return strconv.Itoa(l.ID) }
func (l *DummyListing) GetTitle() string                 { return l.Title }
func (l *DummyListing) GetDescription() string           { return l.Description }
func (l *DummyListing) GetAccessMode() ListingAccessMode { return l.AccessMode }

func (l *DummyListing) SetTitle(title string) error {
	l.Title = title
	return nil
}
func (l *DummyListing) SetDescription(description string) error {
	l.Description = description
	return nil
}
func (l *DummyListing) SetAccessMode(mode ListingAccessMode) error {
	l.AccessMode = mode
	return nil
}

// -----------------------
// Rentals
// -----------------------

type DummyRental struct {
	ID                int
	SourceListingUUID string
	Title             string
	Description       string
}

func (r *DummyRental) GetUUID() string {
	return strconv.Itoa(r.ID)
}

func (r *DummyRental) GetSourceListingUUID() string {
	return r.SourceListingUUID
}

func (r *DummyRental) GetTitle() string {
	return r.Title
}

func (r *DummyRental) GetDescription() string {
	return r.Description
}

func (r *DummyRental) SetTitle(title string) error {
	r.Title = title
	return nil
}

func (r *DummyRental) SetDescription(description string) error {
	r.Description = description
	return nil
}

// -----------------------
// Dummy Data
// -----------------------

var dummyUsers = map[string]*DummyUser{
	"alice": {ID: 1, Username: "alice", PasswordHash: "$2y$10$yyty1BZZiACa4rMHC/ksfOd3fUnwxS5skZAo3Fo6iTCx0bxfZIcMS", IsDev: true},
	"bob":   {ID: 2, Username: "bob", PasswordHash: "$2y$10$yyty1BZZiACa4rMHC/ksfOd3fUnwxS5skZAo3Fo6iTCx0bxfZIcMS", IsDev: false},
	"carol": {ID: 3, Username: "carol", PasswordHash: "$2y$10$yyty1BZZiACa4rMHC/ksfOd3fUnwxS5skZAo3Fo6iTCx0bxfZIcMS", IsDev: true},
}

var globalListings = []Listing{
	&DummyListing{ID: 0, Title: "Node.js Dev Stack", Description: "Node.js 20 with MongoDB...", AccessMode: Public},
	&DummyListing{ID: 1, Title: "Go Microservices Boilerplate", Description: "Go + Kafka setup", AccessMode: Public},
	&DummyListing{ID: 2, Title: "Python ML Environment", Description: "TensorFlow + PyTorch", AccessMode: Private},
}

var userListingMap = map[int][]int{
	1: {0, 1},
	2: {2},
}

var dummyRentals = map[int][]Rental{
	1: {
		&DummyRental{
			ID:                11,
			SourceListingUUID: "0", // Rental originated from listing 0
			Title:             "Go API Testbed",
			Description:       "Running Docker container",
		},
	},
	2: {
		&DummyRental{
			ID:                13,
			SourceListingUUID: "2", // Rental from listing 2
			Title:             "TensorFlow Training Node",
			Description:       "GPU container",
		},
	},
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

func (d *DummyInteractor) QueryPublicRentals(options *RentalQueryOptions) ([]Rental, error) {
	rentals := []Rental{
		&DummyRental{ID: 101, Title: "Public PostgreSQL", Description: "Shared test DB"},
		&DummyRental{ID: 102, Title: "Docker-in-Docker Sandbox", Description: "Run containers inside containers"},
	}

	if options == nil || options.Text == "" {
		return rentals, nil
	}

	query := strings.ToLower(options.Text)
	var result []Rental
	for _, r := range rentals {
		if strings.Contains(strings.ToLower(r.GetTitle()), query) ||
			strings.Contains(strings.ToLower(r.GetDescription()), query) {
			result = append(result, r)
		}
	}
	return result, nil
}

func (d *DummyInteractor) GetPublicListing(uuid string) (Listing, error) {
	index, err := uuidToIndex(uuid)
	if err != nil {
		return nil, errors.New("invalid UUID")
	}

	if index < 0 || index >= len(globalListings) || globalListings[index] == nil {
		return nil, errors.New("listing not found")
	}

	listing := globalListings[index]

	// Check if listing is public
	if listing.GetAccessMode() != Public {
		return nil, errors.New("listing is not public")
	}

	return listing, nil
}

func uuidToIndex(uuid string) (int, error) {
	n, err := strconv.Atoi(uuid)
	if err != nil {
		return 0, errors.New("conversion failed")
	}
	return n, nil
}
