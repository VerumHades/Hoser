package database

type User struct {
	ID           int
	Username     string
	PasswordHash string
	IsDeveloper  bool
}

type Listing struct {
	ID          int
	Title       string
	Description string
}

type Rental struct {
	ID          int
	Title       string
	Description string
}

type RentalQueryOptions struct {
	Text string
}

type Interactor interface {
	// Returns a user by name
	GetUserByName(username string) (*User, error)
	// Returns all listings the user has access to
	GetUserListings(user *User) ([]*Listing, error)
	// Returns all rentals the user has
	GetUserRentals(user *User) ([]*Rental, error)
	// Returns public rentals
	QueryPublicRentals(options *RentalQueryOptions) ([]*Rental, error)
}
