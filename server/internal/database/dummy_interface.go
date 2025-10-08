package database

import "errors"

// DummyInteractor implements Interactor
type DummyInteractor struct{}

// Returns a dummy user if username is not empty
func (d *DummyInteractor) GetUserByName(username string) (*User, error) {
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}
	// passwords: password1234
	return &User{Username: username, PasswordHash: "$2y$10$LHg6JpCH8z5YHKsdYBAtiukBbWLLuR.Gyowp7KzX8R3rOm/hjHPJa"}, nil
}

// Returns some dummy listings
func (d *DummyInteractor) GetUserListings(user *User) ([]*Listing, error) {
	if user == nil {
		return nil, errors.New("user cannot be nil")
	}
	return []*Listing{
		{ID: 1, Title: "Listing 1"},
		{ID: 2, Title: "Listing 2"},
	}, nil
}

// Returns some dummy rentals
func (d *DummyInteractor) GetUserRentals(user *User) ([]*Rental, error) {
	if user == nil {
		return nil, errors.New("user cannot be nil")
	}
	return []*Rental{
		{ID: 1},
		{ID: 2},
	}, nil
}

// Returns some dummy public rentals
func (d *DummyInteractor) QueryPublicRentals() ([]*Rental, error) {
	return []*Rental{
		{ID: 101},
		{ID: 102},
	}, nil
}
