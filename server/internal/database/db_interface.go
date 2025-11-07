package database

type User interface {
	GetUsername() string
	GetPasswordHash() string

	GetListings() ([]Listing, error)
	GetListing(uuid string) (Listing, error)
	DeleteListing(uuid string) (Listing, error)
	CreateListing(title string, description string) (Listing, error)

	GetRentals() ([]Rental, error)
	IsDeveloper() bool
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

	SetTitle(title string) error
	SetDescription(description string) error
	SetAccessMode(mode ListingAccessMode) error
}

type Rental interface {
	GetTitle() string
	GetDescription() string
}

type RentalQueryOptions struct {
	Text string
}

type Interactor interface {
	GetUserByName(username string) (User, error)
	QueryPublicRentals(options *RentalQueryOptions) ([]Rental, error)
}
