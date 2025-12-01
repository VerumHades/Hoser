package database

import (
	"errors"
	"strconv"
	"strings"
)

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
func (c *DummyCurrency) Value() float32    { return c.value }
func (c *DummyCurrency) AsNumber() float32 { return c.value }

func (c *DummyCurrency) SetName(name string)    { c.name = name }
func (c *DummyCurrency) SetShort(short string)  { c.short = short }
func (c *DummyCurrency) SetValue(value float32) { c.value = value }

func (i *DummyStore) NewCurrency(name string, short string, value float32) Currency {
	return &DummyCurrency{
		name:  name,
		short: short,
		value: value,
	}
}

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

// =========================
// Pricing
// =========================

type DummyPricing struct {
	IDVal     string
	TypeVal   PricingType
	AmountVal Currency
}

func (p *DummyPricing) UUID() string                { return p.IDVal }
func (p *DummyPricing) Type() PricingType           { return p.TypeVal }
func (p *DummyPricing) Amount() Currency            { return p.AmountVal }
func (p *DummyPricing) SetType(v PricingType) error { p.TypeVal = v; return nil }
func (p *DummyPricing) SetAmount(a Currency) error  { p.AmountVal = a; return nil }

// =========================
// PricingList
// =========================

type DummyPricingList struct {
	pricings map[string]Pricing
	counter  int // pricing IDs local to listing
}

func NewDummyPricingList() *DummyPricingList {
	return &DummyPricingList{
		pricings: map[string]Pricing{},
		counter:  0,
	}
}

func (pl *DummyPricingList) nextID() string {
	id := strconv.Itoa(pl.counter)
	pl.counter++
	return id
}

func (pl *DummyPricingList) AddPricing(pt PricingType, amount Currency) (Pricing, error) {
	id := pl.nextID()
	pl.pricings[id] = &DummyPricing{
		IDVal:     id,
		TypeVal:   pt,
		AmountVal: amount,
	}
	return pl.pricings[id], nil
}

func (pl *DummyPricingList) GetPricing(id string) Pricing {
	return pl.pricings[id]
}

func (pl *DummyPricingList) RemovePricing(id string) error {
	if _, ok := pl.pricings[id]; !ok {
		return errors.New("pricing not found")
	}
	delete(pl.pricings, id)
	return nil
}
func (pl *DummyPricingList) All() []Pricing {
	list := make([]Pricing, 0, len(pl.pricings))
	for _, p := range pl.pricings {
		list = append(list, p)
	}
	return list
}

// =========================
// Listing
// =========================
type DummyListing struct {
	IDVal    int
	TitleStr string
	DescStr  string
	Access   ListingAccessMode

	HardwareSpec *DummyHardwareSpec

	PricingList *DummyPricingList
}

func (l *DummyListing) UUID() string                  { return strconv.Itoa(l.IDVal) }
func (l *DummyListing) Title() string                 { return l.TitleStr }
func (l *DummyListing) Description() string           { return l.DescStr }
func (l *DummyListing) AccessMode() ListingAccessMode { return l.Access }

func (l *DummyListing) SetTitle(t string) error       { l.TitleStr = t; return nil }
func (l *DummyListing) SetDescription(d string) error { l.DescStr = d; return nil }
func (l *DummyListing) SetAccessMode(m ListingAccessMode) error {
	l.Access = m
	return nil
}

func (l *DummyListing) Pricing() PricingList {
	return l.PricingList
}

func (l *DummyListing) HardwareRequirements() HardwareSpecification {
	return l.HardwareSpec
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
	"alice": {IDVal: 1, UsernameStr: "alice", PasswordStr: "$2a$12$lMgdt2I6sY7wrd/cJtUN1ecjdtSiqIAuiMOuYKwMPiyZck068uGU6", Dev: true},
	"bob":   {IDVal: 2, UsernameStr: "bob", PasswordStr: "$2a$12$lMgdt2I6sY7wrd/cJtUN1ecjdtSiqIAuiMOuYKwMPiyZck068uGU6", Dev: false},
	"carol": {IDVal: 3, UsernameStr: "carol", PasswordStr: "$2a$12$lMgdt2I6sY7wrd/cJtUN1ecjdtSiqIAuiMOuYKwMPiyZck068uGU6", Dev: true},
}

var globalListings []Listing = []Listing{
	&DummyListing{
		IDVal:    0,
		TitleStr: "Node.js Full-Stack Environment",
		DescStr:  "Production-ready Node.js 20 environment with Express, Redis, and PM2 preconfigured.",
		Access:   Public,
		HardwareSpec: &DummyHardwareSpec{
			CPU:  2,
			RAM:  8 * 1024 * 1024 * 1024,
			Disk: 40 * 1024 * 1024 * 1024,
		},
		PricingList: func() *DummyPricingList {
			pl := NewDummyPricingList()
			pl.AddPricing(OneTime, &DummyCurrency{name: "USD", short: "USD", value: 12})
			pl.AddPricing(Monthly, &DummyCurrency{name: "USD", short: "USD", value: 3})
			pl.AddPricing(Yearly, &DummyCurrency{name: "USD", short: "USD", value: 30})
			return pl
		}(),
	},

	&DummyListing{
		IDVal:    1,
		TitleStr: "Go Microservices Platform",
		DescStr:  "Optimized Go 1.22 microservice playground with built-in metrics, tracing, and Kafka connectors.",
		Access:   Public,
		HardwareSpec: &DummyHardwareSpec{
			CPU:  2,
			RAM:  4 * 1024 * 1024 * 1024,
			Disk: 20 * 1024 * 1024 * 1024,
		},
		PricingList: func() *DummyPricingList {
			pl := NewDummyPricingList()
			pl.AddPricing(Monthly, &DummyCurrency{name: "EUR", short: "€", value: 4})
			return pl
		}(),
	},

	&DummyListing{
		IDVal:    2,
		TitleStr: "Rust High-Performance Backend",
		DescStr:  "A blazing-fast Rust environment with Tokio, Axum, SQLx and optimized debug tooling.",
		Access:   Public,
		HardwareSpec: &DummyHardwareSpec{
			CPU:  4,
			RAM:  8 * 1024 * 1024 * 1024,
			Disk: 30 * 1024 * 1024 * 1024,
		},
		PricingList: func() *DummyPricingList {
			pl := NewDummyPricingList()
			pl.AddPricing(Monthly, &DummyCurrency{name: "USD", short: "USD", value: 6})
			pl.AddPricing(Yearly, &DummyCurrency{name: "USD", short: "USD", value: 60})
			return pl
		}(),
	},

	&DummyListing{
		IDVal:    3,
		TitleStr: "Python Data Science Lab",
		DescStr:  "JupyterLab with Python 3.12, Pandas, NumPy, scikit-learn, and GPU-ready environment.",
		Access:   Public,
		HardwareSpec: &DummyHardwareSpec{
			CPU:  4,
			RAM:  16 * 1024 * 1024 * 1024,
			Disk: 60 * 1024 * 1024 * 1024,
		},
		PricingList: func() *DummyPricingList {
			pl := NewDummyPricingList()
			pl.AddPricing(OneTime, &DummyCurrency{name: "EUR", short: "€", value: 15})
			pl.AddPricing(Monthly, &DummyCurrency{name: "EUR", short: "€", value: 5})
			return pl
		}(),
	},

	&DummyListing{
		IDVal:    4,
		TitleStr: "DevOps Sandbox",
		DescStr:  "Kubernetes-in-a-box with Docker, Helm, Grafana, and logging stack prewired.",
		Access:   Public,
		HardwareSpec: &DummyHardwareSpec{
			CPU:  4,
			RAM:  8 * 1024 * 1024 * 1024,
			Disk: 50 * 1024 * 1024 * 1024,
		},
		PricingList: func() *DummyPricingList {
			pl := NewDummyPricingList()
			pl.AddPricing(Monthly, &DummyCurrency{name: "Credits", short: "CR", value: 5})
			return pl
		}(),
	},

	&DummyListing{
		IDVal:    5,
		TitleStr: "AI / LLM Experimentation Rig",
		DescStr:  "Environment tailored for ML tests, including PyTorch, Transformers, and CUDA support.",
		Access:   Public,
		HardwareSpec: &DummyHardwareSpec{
			CPU:  8,
			RAM:  32 * 1024 * 1024 * 1024,
			Disk: 120 * 1024 * 1024 * 1024,
		},
		PricingList: func() *DummyPricingList {
			pl := NewDummyPricingList()
			pl.AddPricing(OneTime, &DummyCurrency{name: "USD", short: "USD", value: 20})
			pl.AddPricing(Monthly, &DummyCurrency{name: "USD", short: "USD", value: 10})
			return pl
		}(),
	},

	&DummyListing{
		IDVal:    6,
		TitleStr: "PHP + MySQL Legacy Stack",
		DescStr:  "Stable LAMP environment ideal for legacy apps or migrations.",
		Access:   Public,
		HardwareSpec: &DummyHardwareSpec{
			CPU:  2,
			RAM:  2 * 1024 * 1024 * 1024,
			Disk: 15 * 1024 * 1024 * 1024,
		},
		PricingList: func() *DummyPricingList {
			pl := NewDummyPricingList()
			pl.AddPricing(Monthly, &DummyCurrency{name: "EUR", short: "€", value: 2})
			return pl
		}(),
	},

	&DummyListing{
		IDVal:    7,
		TitleStr: "C# .NET Cloud API Environment",
		DescStr:  "ASP.NET Core 8 project runner with EF Core, Redis cache, and OpenAPI tools.",
		Access:   Public,
		HardwareSpec: &DummyHardwareSpec{
			CPU:  4,
			RAM:  8 * 1024 * 1024 * 1024,
			Disk: 40 * 1024 * 1024 * 1024,
		},
		PricingList: func() *DummyPricingList {
			pl := NewDummyPricingList()
			pl.AddPricing(Yearly, &DummyCurrency{name: "USD", short: "USD", value: 45})
			return pl
		}(),
	},

	&DummyListing{
		IDVal:    8,
		TitleStr: "Elixir Phoenix Realtime Platform",
		DescStr:  "Optimized BEAM VM system for concurrent workloads, perfect for chat apps or live updates.",
		Access:   Public,
		HardwareSpec: &DummyHardwareSpec{
			CPU:  2,
			RAM:  4 * 1024 * 1024 * 1024,
			Disk: 25 * 1024 * 1024 * 1024,
		},
		PricingList: func() *DummyPricingList {
			pl := NewDummyPricingList()
			pl.AddPricing(Monthly, &DummyCurrency{name: "Credits", short: "CR", value: 4})
			return pl
		}(),
	},

	&DummyListing{
		IDVal:    9,
		TitleStr: "Blockchain Solidity Playground",
		DescStr:  "Smart contract development VM with Foundry, Hardhat, and test networks included.",
		Access:   Public,
		HardwareSpec: &DummyHardwareSpec{
			CPU:  2,
			RAM:  4 * 1024 * 1024 * 1024,
			Disk: 30 * 1024 * 1024 * 1024,
		},
		PricingList: func() *DummyPricingList {
			pl := NewDummyPricingList()
			pl.AddPricing(OneTime, &DummyCurrency{name: "USD", short: "USD", value: 14})
			pl.AddPricing(Monthly, &DummyCurrency{name: "USD", short: "USD", value: 4})
			return pl
		}(),
	},
}

var userListingMap = map[int][]int{
	1: {0, 1, 2, 3, 4, 5},
	2: {},
	3: {},
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
