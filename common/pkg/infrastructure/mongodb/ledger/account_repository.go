package mongodbledger

import (
	"context"
	"errors"

	"common/pkg/domain/ledger"
	"common/pkg/shared"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoAccountRepository implements both AccountCommandRepository and AccountQueryRepository
type MongoAccountRepository struct {
	collection *mongo.Collection
}

// NewMongoAccountRepository constructs a repository backed by the given Mongo collection
func NewMongoAccountRepository(collection *mongo.Collection) *MongoAccountRepository {
	if collection == nil {
		panic("mongo collection must not be nil")
	}

	return &MongoAccountRepository{
		collection: collection,
	}
}

// EnsureAccountIndexes creates indexes for efficient queries
func (repo *MongoAccountRepository) EnsureIndexes(ctx context.Context) error {
	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "ownerType", Value: 1},
				{Key: "ownerID", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "accountType", Value: 1},
			},
		},
	}

	_, err := repo.collection.Indexes().CreateMany(ctx, indexModels)
	return err
}

// -------------------- Command Repository --------------------

// Create inserts a new account inside a transaction
func (repo *MongoAccountRepository) Create(
	ctx context.Context,
	transaction shared.Transaction,
	account *ledger.Account,
) (*ledger.Account, error) {
	if account == nil {
		return nil, errors.New("account cannot be nil")
	}

	sessionCtx := transaction.SessionContext(ctx)

	_, err := repo.collection.InsertOne(sessionCtx, account)
	if mongo.IsDuplicateKeyError(err) {
		return nil, shared.ErrAlreadyExists
	}

	return account, err
}

// Update replaces an existing account inside a transaction
func (repo *MongoAccountRepository) Update(
	ctx context.Context,
	transaction shared.Transaction,
	account *ledger.Account,
) (*ledger.Account, error) {
	if account == nil {
		return nil, errors.New("account cannot be nil")
	}

	sessionCtx := transaction.SessionContext(ctx)

	result, err := repo.collection.ReplaceOne(
		sessionCtx,
		bson.M{"_id": account.ID()},
		account,
	)
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, shared.ErrNotFound
	}

	return account, nil
}

// Delete removes an account by ID inside a transaction
func (repo *MongoAccountRepository) Delete(
	ctx context.Context,
	transaction shared.Transaction,
	accountID shared.AccountID,
) error {
	sessionCtx := transaction.SessionContext(ctx)

	result, err := repo.collection.DeleteOne(sessionCtx, bson.M{"_id": accountID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return shared.ErrNotFound
	}

	return nil
}

// -------------------- Query Repository --------------------

// GetByID retrieves an account by its ID
func (repo *MongoAccountRepository) GetByID(
	ctx context.Context,
	accountID shared.AccountID,
) (*ledger.Account, error) {
	var account ledger.Account

	err := repo.collection.FindOne(ctx, bson.M{"_id": accountID}).Decode(&account)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, shared.ErrNotFound
	}

	return &account, err
}

// GetByOwner retrieves all accounts owned by a specific owner
func (repo *MongoAccountRepository) GetByOwner(
	ctx context.Context,
	ownerType ledger.AccountOwnerType,
	ownerID string,
) ([]*ledger.Account, error) {
	filter := bson.M{
		"ownerType": ownerType,
		"ownerID":   ownerID,
	}

	cursor, err := repo.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var accounts []*ledger.Account
	for cursor.Next(ctx) {
		var account ledger.Account
		if err := cursor.Decode(&account); err != nil {
			return nil, err
		}
		accounts = append(accounts, &account)
	}

	return accounts, nil
}

// GetFirstByOwner retrieves the first account for a given owner
func (repo *MongoAccountRepository) GetFirstByOwner(
	ctx context.Context,
	ownerType ledger.AccountOwnerType,
	ownerID string,
) (*ledger.Account, error) {
	filter := bson.M{
		"ownerType": ownerType,
		"ownerID":   ownerID,
	}

	var account ledger.Account
	err := repo.collection.FindOne(ctx, filter).Decode(&account)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, shared.ErrNotFound
	}

	return &account, err
}

// GetByType retrieves all accounts of a specific account type
func (repo *MongoAccountRepository) GetByType(
	ctx context.Context,
	accountType ledger.AccountType,
) ([]*ledger.Account, error) {
	filter := bson.M{
		"accountType": accountType,
	}

	cursor, err := repo.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var accounts []*ledger.Account
	for cursor.Next(ctx) {
		var account ledger.Account
		if err := cursor.Decode(&account); err != nil {
			return nil, err
		}
		accounts = append(accounts, &account)
	}

	return accounts, nil
}
