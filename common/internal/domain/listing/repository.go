package listing

import "common/internal/shared"

type ListingRepository interface {
	Save(listing *Listing) error

	GetByID(listingID shared.ListingID) (*Listing, error)

	ListByAuthor(authorID shared.UserID) ([]*Listing, error)
	ListAll() ([]*Listing, error) // new method to list all listings

	Delete(listingID shared.ListingID) error
}

// GitHubSetupRepository defines persistence operations for GitHub setups.
type GitHubSetupRepository interface {
	Save(setup *GitHubSetupDefinition) error

	GetByID(id shared.SetupID) (*GitHubSetupDefinition, error)
	GetByListingID(listingID shared.ListingID) ([]*GitHubSetupDefinition, error)

	DeleteByID(id shared.SetupID) error
	DeleteByListingID(listingID shared.ListingID) error
}
