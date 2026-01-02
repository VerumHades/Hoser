package mongodbregistry

import "go.mongodb.org/mongo-driver/mongo"

type DatabaseRegistry struct {
	Users              *mongo.Collection
	Listings           *mongo.Collection
	Libraries          *mongo.Collection
	Accounts           *mongo.Collection
	LedgerTransactions *mongo.Collection
	Settlements        *mongo.Collection
	InstanceContracts  *mongo.Collection
	HardwareRates      *mongo.Collection
	GithubSetups       *mongo.Collection
}

// NewDatabaseRegistry initializes all collections in one place
func NewDatabaseRegistry(db *mongo.Database) *DatabaseRegistry {
	return &DatabaseRegistry{
		Users:              db.Collection("users"),
		Listings:           db.Collection("listings"),
		Libraries:          db.Collection("libraries"),
		Accounts:           db.Collection("billing_accounts"),
		LedgerTransactions: db.Collection("ledger_transaction"),
		Settlements:        db.Collection("settlement"),
		InstanceContracts:  db.Collection("instances"),
		HardwareRates:      db.Collection("hardware_costs"),
		GithubSetups:       db.Collection("github_setups"),
	}
}
