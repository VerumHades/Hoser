package database

import (
	"errors"
	"strconv"
	"strings"
)

type DummyCurrency struct {
	Name  string
	Short string
	Value float32
}

func (c *DummyCurrency) GetAsNumber() float32 { return c.Value }
func (c *DummyCurrency) GetName() string      { return c.Name }
func (c *DummyCurrency) GetShort() string     { return c.Short }

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

	SinglePurchasePrice      Currency
	MonthlySubscriptionPrice Currency
	MonthlyHardwarePrice     Currency
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
func (l *DummyListing) GetSinglePurchasePrice() Currency {
	return l.SinglePurchasePrice
}
func (l *DummyListing) GetMonthlySubscriptionPrice() Currency {
	return l.MonthlySubscriptionPrice
}
func (l *DummyListing) GetMonthlyHardwarePrice() Currency {
	return l.MonthlyHardwarePrice
}

func (l *DummyListing) SetSinglePurchasePrice(price Currency) error {
	l.SinglePurchasePrice = price
	return nil
}
func (l *DummyListing) SetMonthlySubscriptionPrice(price Currency) error {
	l.MonthlySubscriptionPrice = price
	return nil
}
func (l *DummyListing) SetMonthlyHardwarePrice(price Currency) error {
	l.MonthlyHardwarePrice = price
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
	&DummyListing{
		ID:                       0,
		Title:                    "Node.js Dev Stack",
		Description:              "Node.js 20 with MongoDB...",
		AccessMode:               Public,
		SinglePurchasePrice:      &DummyCurrency{Name: "Credits", Short: "CR", Value: 10},
		MonthlySubscriptionPrice: &DummyCurrency{Name: "Credits", Short: "CR", Value: 2},
		MonthlyHardwarePrice:     &DummyCurrency{Name: "Credits", Short: "CR", Value: 1}, // Added
	},

	&DummyListing{
		ID:                       1,
		Title:                    "Go Microservices Boilerplate",
		Description:              "Go + Kafka setup",
		AccessMode:               Public,
		MonthlySubscriptionPrice: &DummyCurrency{Name: "Credits", Short: "CR", Value: 3},
		MonthlyHardwarePrice:     &DummyCurrency{Name: "Credits", Short: "CR", Value: 1}, // Added
	},

	&DummyListing{
		ID:                   2,
		Title:                "Python ML Environment",
		Description:          "TensorFlow + PyTorch",
		AccessMode:           Private,
		MonthlyHardwarePrice: &DummyCurrency{Name: "Credits", Short: "CR", Value: 7}, // unchanged
	},

	// NEW LISTINGS

	&DummyListing{
		ID:                       3,
		Title:                    "Unity Game Build Cloud",
		Description:              "Automated Unity build CI runners.",
		AccessMode:               Public,
		SinglePurchasePrice:      &DummyCurrency{Name: "Credits", Short: "CR", Value: 5},
		MonthlySubscriptionPrice: &DummyCurrency{Name: "Credits", Short: "CR", Value: 2},
		MonthlyHardwarePrice:     &DummyCurrency{Name: "Credits", Short: "CR", Value: 1},
	},

	&DummyListing{
		ID:          4,
		Title:       "Rust Embedded Toolchain",
		Description: "Cross-compile firmware for ARM MCUs.",
		AccessMode:  Private,
		// No pricing at all → valid, optional
	},

	&DummyListing{
		ID:                       5,
		Title:                    "Shared GPU ML Lab Node",
		Description:              "Multi-user GPU workspace environment.",
		AccessMode:               Public,
		MonthlySubscriptionPrice: &DummyCurrency{Name: "Credits", Short: "CR", Value: 6},
		MonthlyHardwarePrice:     &DummyCurrency{Name: "Credits", Short: "CR", Value: 3}, // required & present
	},
}

var userListingMap = map[int][]int{
	1: {0, 1, 3}, // Alice owns listings 0, 1, 3
	2: {2},       // Bob owns listing 2
	3: {4},       // Carol owns listing 4
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

func (d *DummyInteractor) GetPublicListing(uuid string) (Listing, error) {
	index, err := uuidToIndex(uuid)
	if err != nil {
		return nil, errors.New("invalid UUID")
	}

	// Out of bounds or deleted listing
	if index < 0 || index >= len(globalListings) || globalListings[index] == nil {
		return nil, errors.New("listing not found")
	}

	listing := globalListings[index]

	// Ensure listing is public
	if listing.GetAccessMode() != Public {
		return nil, errors.New("listing is not public")
	}

	return listing, nil
}

func (d *DummyInteractor) QueryPublicListings(options *RentalQueryOptions) ([]Listing, error) {
	var results []Listing

	for _, l := range globalListings {
		// Skip deleted or nil entries
		if l == nil {
			continue
		}

		// Only include public listings
		if l.GetAccessMode() != Public {
			continue
		}

		// If no search text, include all public listings
		if options == nil || options.Text == "" {
			results = append(results, l)
			continue
		}

		// Text filtering (title or description)
		q := strings.ToLower(options.Text)
		if strings.Contains(strings.ToLower(l.GetTitle()), q) ||
			strings.Contains(strings.ToLower(l.GetDescription()), q) {
			results = append(results, l)
		}
	}

	return results, nil
}

func uuidToIndex(uuid string) (int, error) {
	n, err := strconv.Atoi(uuid)
	if err != nil {
		return 0, errors.New("conversion failed")
	}
	return n, nil
}
