package database

type Currency interface {
	Name() string  // e.g., "US Dollar"
	Short() string // e.g., "USD"

	Value() float32    // numeric value, no formatting
	AsNumber() float32 // identical to Value

	SetName(name string)
	SetShort(short string)
	SetValue(value float32)
}

type PricingType int

const (
	OneTime PricingType = iota
	Monthly
	Yearly
)

type Pricing interface {
	UUID() string
	Type() PricingType
	Amount() Currency

	SetType(value PricingType) error
	SetAmount(amount Currency) error
}

type PricingList interface {
	AddPricing(pricing_type PricingType, amount Currency) (Pricing, error)
	GetPricing(id string) Pricing
	RemovePricing(id string) error
	All() []Pricing
}

// HardwareSpecification represents hardware configuration requirements.
type HardwareSpecification interface {
	CPUCount() int
	RAMBytes() int64
	DiskBytes() int64

	SetCPUCount(count int) error
	SetRAMBytes(bytes int64) error
	SetDiskBytes(bytes int64) error

	MonthlyPrice() Currency
}

type ListingAccessMode int

const (
	Private ListingAccessMode = iota
	Public
)

// Listing represents a hardware offering posted by a user.
type Listing interface {
	UUID() string

	Author() User

	Title() string
	Description() string
	AccessMode() ListingAccessMode

	SetTitle(title string) error
	SetDescription(desc string) error
	SetAccessMode(mode ListingAccessMode) error

	Pricing() PricingList
	HardwareRequirements() HardwareSpecification
}

type StateSpecification interface {
	Running() bool
	SetRunning(state bool) error
}

// Rental is an instance of a user renting a listing.
type Rental interface {
	UUID() string
	SourceListing() Listing
	HardwareSetup() HardwareSpecification
	Specification() StateSpecification
}

// User represents a system account with access to listings and rentals.
type User interface {
	ID() string
	Username() string
	PasswordHash() string

	Listings() ([]Listing, error)
	Listing(uuid string) (Listing, error)
	CreateListing(title, description string) (Listing, error)
	DeleteListing(uuid string) error

	Rentals() ([]Rental, error)
	IsDeveloper() bool

	Rent(listingUUID string) error
}

type ListingQueryOptions struct {
	Text string // simple full-text search
}

// Store defines all operations for retrieving users and public listings.
type Store interface {
	NewCurrency(name string, short string, value float32) Currency

	GetUserByID(id string) (User, error)
	GetUserByName(username string) (User, error)

	GetPublicListing(uuid string) (Listing, error)
	QueryPublicListings(opts *ListingQueryOptions) ([]Listing, error)
}
