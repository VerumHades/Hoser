package database

type Currency interface {
	GetAsNumber() float32
	GetName() string
	GetShort() string
}

type User interface {
	GetUsername() string
	GetPasswordHash() string

	GetListings() ([]Listing, error)
	GetListing(uuid string) (Listing, error)
	DeleteListing(uuid string) (Listing, error)
	CreateListing(title string, description string) (Listing, error)

	GetRentals() ([]Rental, error)
	IsDeveloper() bool

	Rent(listingUUID string) error
}

type ListingAccessMode int

const (
	Private ListingAccessMode = iota
	Public
)

type Listing interface {
	GetUUID() string

	GetTitle() string
	GetDescription() string
	GetAccessMode() ListingAccessMode

	GetSinglePurchasePrice() Currency
	GetMonthlySubscriptionPrice() Currency
	GetMonthlyHardwarePrice() Currency

	SetTitle(title string) error
	SetDescription(description string) error
	SetAccessMode(mode ListingAccessMode) error

	SetSinglePurchasePrice(price Currency) error
	SetMonthlySubscriptionPrice(price Currency) error
	SetMonthlyHardwarePrice(price Currency) error
}

type Rental interface {
	GetUUID() string
	GetSourceListingUUID() string

	GetTitle() string
	GetDescription() string

	SetTitle(title string) error
	SetDescription(description string) error
}

type RentalQueryOptions struct {
	Text string
}

type Interactor interface {
	GetUserByName(username string) (User, error)
	GetPublicListing(uuid string) (Listing, error)
	QueryPublicListings(options *RentalQueryOptions) ([]Listing, error)
}
