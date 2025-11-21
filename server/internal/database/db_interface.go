package database

// Currency represents any currency value your system supports.
type Currency interface {
	// Getters
	Name() string      // e.g., "US Dollar"
	Short() string     // e.g., "USD"
	AsNumber() float32 // numeric value, no formatting

	// Setters
	SetName(name string)
	SetShort(short string)
	SetValue(value float32)
}

// HardwareSpecification represents hardware configuration requirements.
type HardwareSpecification interface {
	CPUCount() int
	RAMBytes() int64
	DiskBytes() int64

	SetCPUCount(count int) error
	SetRAMBytes(bytes int64) error
	SetDiskBytes(bytes int64) error
}

type ListingAccessMode int

const (
	Private ListingAccessMode = iota
	Public
)

// Listing represents a hardware offering posted by a user.
type Listing interface {
	UUID() string

	Title() string
	Description() string
	AccessMode() ListingAccessMode

	SetTitle(title string) error
	SetDescription(desc string) error
	SetAccessMode(mode ListingAccessMode) error

	SinglePurchasePrice() Currency
	MonthlySubscriptionPrice() Currency
	MonthlyHardwarePrice() Currency

	HardwareRequirements() HardwareSpecification
}

type StateSpecification interface {
	Running() bool
	SetRunning(state bool) error
}

// Rental is an instance of a user renting a listing.
type Rental interface {
	UUID() string
	SourceListingUUID() string
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
	GetUserByID(id string) (User, error)
	GetUserByName(username string) (User, error)

	GetPublicListing(uuid string) (Listing, error)
	QueryPublicListings(opts *ListingQueryOptions) ([]Listing, error)
}
