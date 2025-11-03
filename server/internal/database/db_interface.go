package database

type User interface {
	GetUsername() string
	GetPasswordHash() string
	GetListings() ([]Listing, error)
	CreateListing(title string, description string) (Listing, error)

	GetRentals() ([]Rental, error)
	IsDeveloper() bool
}

type Listing interface {
	GetTitle() string
	GetDescription() string
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
