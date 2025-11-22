package database

import (
	"errors"
	"strconv"
	"strings"
)

type StubHardwareSpecification struct {
	CPUCountVal  int
	RAMBytesVal  int64
	DiskBytesVal int64

	SetCPUCountErr  error
	SetRAMBytesErr  error
	SetDiskBytesErr error
}

func (s *StubHardwareSpecification) CPUCount() int {
	return s.CPUCountVal
}

func (s *StubHardwareSpecification) RAMBytes() int64 {
	return s.RAMBytesVal
}

func (s *StubHardwareSpecification) DiskBytes() int64 {
	return s.DiskBytesVal
}

func (s *StubHardwareSpecification) SetCPUCount(count int) error {
	if s.SetCPUCountErr != nil {
		return s.SetCPUCountErr
	}
	s.CPUCountVal = count
	return nil
}

func (s *StubHardwareSpecification) SetRAMBytes(bytes int64) error {
	if s.SetRAMBytesErr != nil {
		return s.SetRAMBytesErr
	}
	s.RAMBytesVal = bytes
	return nil
}

func (s *StubHardwareSpecification) SetDiskBytes(bytes int64) error {
	if s.SetDiskBytesErr != nil {
		return s.SetDiskBytesErr
	}
	s.DiskBytesVal = bytes
	return nil
}

//
// =========================
// Currency
// =========================
//

type DummyCurrency struct {
	name  string
	short string
	value float32
}

func (c *DummyCurrency) Name() string      { return c.name }
func (c *DummyCurrency) Short() string     { return c.short }
func (c *DummyCurrency) AsNumber() float32 { return c.value }

func (c *DummyCurrency) SetName(name string)    { c.name = name }
func (c *DummyCurrency) SetShort(short string)  { c.short = short }
func (c *DummyCurrency) SetValue(value float32) { c.value = value }

//
// =========================
// User
// =========================
//

type DummyUser struct {
	IDVal       int
	UsernameStr string
	PasswordStr string
	Dev         bool
}

func (u *DummyUser) ID() string           { return strconv.Itoa(u.IDVal) }
func (u *DummyUser) Username() string     { return u.UsernameStr }
func (u *DummyUser) PasswordHash() string { return u.PasswordStr }
func (u *DummyUser) IsDeveloper() bool    { return u.Dev }
func (u *DummyUser) Listings() ([]Listing, error) {
	indices := userListingMap[u.IDVal]
	var listings []Listing
	for _, idx := range indices {
		if idx >= 0 && idx < len(globalListings) && globalListings[idx] != nil {
			listings = append(listings, globalListings[idx])
		}
	}
	return listings, nil
}

func (u *DummyUser) Rentals() ([]Rental, error) {
	r, ok := dummyRentals[u.IDVal]
	if !ok {
		return []Rental{}, nil
	}
	return r, nil
}

func (u *DummyUser) CreateListing(title string, description string) (Listing, error) {
	listing := &DummyListing{
		IDVal:    len(globalListings),
		TitleStr: title,
		DescStr:  description,
		Access:   Private,
	}
	globalListings = append(globalListings, listing)
	userListingMap[u.IDVal] = append(userListingMap[u.IDVal], listing.IDVal)
	return listing, nil
}

func (u *DummyUser) Listing(uuid string) (Listing, error) {
	index, err := uuidToIndex(uuid)
	if err != nil {
		return nil, err
	}
	if index < 0 || index >= len(globalListings) || globalListings[index] == nil {
		return nil, errors.New("invalid id")
	}
	for _, owned := range userListingMap[u.IDVal] {
		if owned == index {
			return globalListings[index], nil
		}
	}
	return nil, errors.New("listing not owned by user")
}

func (u *DummyUser) DeleteListing(uuid string) error {
	index, err := uuidToIndex(uuid)
	if err != nil {
		return err
	}
	if index < 0 || index >= len(globalListings) || globalListings[index] == nil {
		return errors.New("invalid listing ID")
	}

	ownedIndices := userListingMap[u.IDVal]
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
		return errors.New("listing not owned by user")
	}
	userListingMap[u.IDVal] = newOwned
	globalListings[index] = nil
	return nil
}

func (u *DummyUser) Rent(listingUUID string) error {
	index, err := uuidToIndex(listingUUID)
	if err != nil {
		return errors.New("invalid listing UUID")
	}
	if index < 0 || index >= len(globalListings) || globalListings[index] == nil {
		return errors.New("listing not found")
	}

	listing := globalListings[index]

	if listing.AccessMode() != Public {
		return errors.New("listing is not public")
	}

	if dummyRentals[u.IDVal] == nil {
		dummyRentals[u.IDVal] = []Rental{}
	}

	rental := &DummyRental{
		ID:                 len(dummyRentals[u.IDVal]) + 100,
		_SourceListingUUID: listing.UUID(),
		Title:              listing.Title(),
		Description:        listing.Description(),

		Hardware:  &DummyHardwareSpec{},       // optional dummy hardware
		StateSpec: &DummyStateSpecification{}, // new state spec
	}

	dummyRentals[u.IDVal] = append(dummyRentals[u.IDVal], rental)
	return nil
}

type DummyHardwareSpec struct {
	CPU  int
	RAM  int64
	Disk int64
}

func (h *DummyHardwareSpec) CPUCount() int    { return h.CPU }
func (h *DummyHardwareSpec) RAMBytes() int64  { return h.RAM }
func (h *DummyHardwareSpec) DiskBytes() int64 { return h.Disk }

func (h *DummyHardwareSpec) SetCPUCount(c int) error    { h.CPU = c; return nil }
func (h *DummyHardwareSpec) SetRAMBytes(b int64) error  { h.RAM = b; return nil }
func (h *DummyHardwareSpec) SetDiskBytes(b int64) error { h.Disk = b; return nil }

func (h *DummyHardwareSpec) MonthlyPrice() Currency {
	return &DummyCurrency{
		name:  "Credits",
		short: "CR",
		value: float32(h.CPU)*1 +
			float32(h.RAM)/(1024*1024*1024)*0.5 +
			float32(h.Disk)/(1024*1024*1024)*0.1, // arbitrary example
	}
}

//
// =========================
// Listing
// =========================
//

type DummyListing struct {
	IDVal    int
	TitleStr string
	DescStr  string
	Access   ListingAccessMode

	SinglePrice     Currency
	MonthlySubPrice Currency
	MonthlyHWPrice  Currency
	HardwareSpec    *DummyHardwareSpec
}

func (l *DummyListing) UUID() string                  { return strconv.Itoa(l.IDVal) }
func (l *DummyListing) Title() string                 { return l.TitleStr }
func (l *DummyListing) Description() string           { return l.DescStr }
func (l *DummyListing) AccessMode() ListingAccessMode { return l.Access }

func (l *DummyListing) SetTitle(title string) error {
	l.TitleStr = title
	return nil
}
func (l *DummyListing) SetDescription(description string) error {
	l.DescStr = description
	return nil
}
func (l *DummyListing) SetAccessMode(mode ListingAccessMode) error {
	l.Access = mode
	return nil
}

func (l *DummyListing) SinglePurchasePrice() Currency      { return l.SinglePrice }
func (l *DummyListing) MonthlySubscriptionPrice() Currency { return l.MonthlySubPrice }
func (l *DummyListing) MonthlyHardwarePrice() Currency     { return l.MonthlyHWPrice }

func (l *DummyListing) SetSinglePurchasePrice(price Currency) error {
	l.SinglePrice = price
	return nil
}
func (l *DummyListing) SetMonthlySubscriptionPrice(price Currency) error {
	l.MonthlySubPrice = price
	return nil
}
func (l *DummyListing) SetMonthlyHardwarePrice(price Currency) error {
	l.MonthlyHWPrice = price
	return nil
}

func (l *DummyListing) HardwareRequirements() HardwareSpecification {
	// Dummy implementation — your real implementation goes here
	return l.HardwareSpec
}

func (l *DummyListing) ClearSinglePurchasePrice() error {
	l.SinglePrice = nil
	return nil
}

func (l *DummyListing) ClearSubscriptionPrice() error {
	l.MonthlySubPrice = nil
	return nil
}

// -----------------------
// State Specification
// -----------------------

type DummyStateSpecification struct {
	running bool
}

func (s *DummyStateSpecification) Running() bool {
	return s.running
}

func (s *DummyStateSpecification) SetRunning(state bool) error {
	s.running = state
	return nil
}

//
// =========================
// Rentals
// =========================
//

type DummyRental struct {
	ID                 int
	_SourceListingUUID string
	Title              string
	Description        string

	Hardware  *DummyHardwareSpec
	StateSpec *DummyStateSpecification
}

func (r *DummyRental) UUID() string {
	return strconv.Itoa(r.ID)
}

func (r *DummyRental) SourceListingUUID() string {
	return r._SourceListingUUID
}

func (r *DummyRental) HardwareSetup() HardwareSpecification {
	return r.Hardware
}

func (r *DummyRental) Specification() StateSpecification {
	return r.StateSpec
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

func (r *DummyRental) SetDescription(desc string) error {
	r.Description = desc
	return nil
}

//
// =========================
// Dummy Data
// =========================
//

var dummyUsers = map[string]*DummyUser{
	"alice": {IDVal: 1, UsernameStr: "alice", PasswordStr: "hash", Dev: true},
	"bob":   {IDVal: 2, UsernameStr: "bob", PasswordStr: "hash", Dev: false},
	"carol": {IDVal: 3, UsernameStr: "carol", PasswordStr: "hash", Dev: true},
}

var globalListings []Listing = []Listing{
	&DummyListing{
		IDVal:           0,
		TitleStr:        "Node.js Dev Stack",
		DescStr:         "Node.js",
		Access:          Public,
		SinglePrice:     &DummyCurrency{name: "Credits", short: "CR", value: 10},
		MonthlySubPrice: &DummyCurrency{name: "Credits", short: "CR", value: 2},
		MonthlyHWPrice:  &DummyCurrency{name: "Credits", short: "CR", value: 1},
		HardwareSpec: &DummyHardwareSpec{
			CPU:  2,
			RAM:  4 * 1024 * 1024 * 1024,  // 4GB
			Disk: 20 * 1024 * 1024 * 1024, // 20GB
		},
	},
	&DummyListing{
		IDVal:           1,
		TitleStr:        "Go Microservices",
		DescStr:         "Go + Kafka",
		Access:          Public,
		MonthlySubPrice: &DummyCurrency{name: "Credits", short: "CR", value: 3},
		MonthlyHWPrice:  &DummyCurrency{name: "Credits", short: "CR", value: 1},
		HardwareSpec: &DummyHardwareSpec{
			CPU:  2,
			RAM:  4 * 1024 * 1024 * 1024,  // 4GB
			Disk: 20 * 1024 * 1024 * 1024, // 20GB
		},
	},
}

var userListingMap = map[int][]int{
	1: {0},
	2: {},
	3: {1},
}

var dummyRentals = map[int][]Rental{
	1: {
		&DummyRental{
			ID:                 11,
			_SourceListingUUID: "0",
			Hardware: &DummyHardwareSpec{
				CPU:  2,
				RAM:  4 * 1024 * 1024 * 1024,  // 4GB
				Disk: 20 * 1024 * 1024 * 1024, // 20GB
			},
			StateSpec: &DummyStateSpecification{
				running: true,
			},
		},
		&DummyRental{
			ID:                 12,
			_SourceListingUUID: "3",
			Hardware: &DummyHardwareSpec{
				CPU:  4,
				RAM:  8 * 1024 * 1024 * 1024,  // 8GB
				Disk: 40 * 1024 * 1024 * 1024, // 40GB
			},
			StateSpec: &DummyStateSpecification{
				running: false,
			},
		},
	},

	2: {
		&DummyRental{
			ID:                 21,
			_SourceListingUUID: "2",
			Hardware: &DummyHardwareSpec{
				CPU:  8,
				RAM:  16 * 1024 * 1024 * 1024, // 16GB
				Disk: 60 * 1024 * 1024 * 1024, // 60GB
			},
			StateSpec: &DummyStateSpecification{
				running: true,
			},
		},
	},

	3: {
		&DummyRental{
			ID:                 31,
			_SourceListingUUID: "4",
			Hardware: &DummyHardwareSpec{
				CPU:  2,
				RAM:  2 * 1024 * 1024 * 1024,  // 2GB
				Disk: 10 * 1024 * 1024 * 1024, // 10GB
			},
			StateSpec: &DummyStateSpecification{
				running: false,
			},
		},
	},
}

//
// =========================
// Store (Interactor)
// =========================
//

type DummyStore struct{}

func (d *DummyStore) GetUserByID(id string) (User, error) {
	for _, u := range dummyUsers {
		if strconv.Itoa(u.IDVal) == id {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (d *DummyStore) GetUserByName(username string) (User, error) {
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}
	u, ok := dummyUsers[username]
	if !ok {
		return nil, errors.New("user not found")
	}
	return u, nil
}

func (d *DummyStore) GetPublicListing(uuid string) (Listing, error) {
	index, err := uuidToIndex(uuid)
	if err != nil {
		return nil, errors.New("invalid UUID")
	}

	if index < 0 || index >= len(globalListings) || globalListings[index] == nil {
		return nil, errors.New("listing not found")
	}

	listing := globalListings[index]

	if listing.AccessMode() != Public {
		return nil, errors.New("listing is not public")
	}

	return listing, nil
}

func (d *DummyStore) QueryPublicListings(opts *ListingQueryOptions) ([]Listing, error) {
	var results []Listing

	for _, l := range globalListings {
		if l == nil {
			continue
		}
		if l.AccessMode() != Public {
			continue
		}
		if opts == nil || opts.Text == "" {
			results = append(results, l)
			continue
		}

		q := strings.ToLower(opts.Text)
		if strings.Contains(strings.ToLower(l.Title()), q) ||
			strings.Contains(strings.ToLower(l.Description()), q) {
			results = append(results, l)
		}
	}

	return results, nil
}

//
// =========================
// Helpers
// =========================
//

func uuidToIndex(uuid string) (int, error) {
	n, err := strconv.Atoi(uuid)
	if err != nil {
		return 0, errors.New("conversion failed")
	}
	return n, nil
}
