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

// MongoLedgerTransactionRepository implements both command and query repositories
type MongoLedgerTransactionRepository struct {
	collection *mongo.Collection
}

// NewMongoLedgerTransactionRepository constructs a repository backed by the given Mongo collection
func NewMongoLedgerTransactionRepository(collection *mongo.Collection) *MongoLedgerTransactionRepository {
	if collection == nil {
		panic("mongo collection must not be nil")
	}
	return &MongoLedgerTransactionRepository{collection: collection}
}

// EnsureLedgerTransactionIndexes creates indexes for efficient queries
func (repo *MongoLedgerTransactionRepository) EnsureIndexes(ctx context.Context) error {
	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "referenceType", Value: 1},
				{Key: "referenceID", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "entries.accountID", Value: 1},
				{Key: "referenceID", Value: 1},
			},
		},
	}

	_, err := repo.collection.Indexes().CreateMany(ctx, indexModels)
	return err
}

// -------------------- Command Repository --------------------

func (repo *MongoLedgerTransactionRepository) Create(
	ctx context.Context,
	transaction shared.Transaction,
	ledgerTransaction *ledger.LedgerTransaction,
) error {
	if ledgerTransaction == nil {
		return errors.New("ledger transaction cannot be nil")
	}

	sessionCtx := transaction.SessionContext(ctx)

	_, err := repo.collection.InsertOne(sessionCtx, ledgerTransaction)
	if mongo.IsDuplicateKeyError(err) {
		return shared.ErrAlreadyExists
	}

	return err
}

func (repo *MongoLedgerTransactionRepository) Update(
	ctx context.Context,
	transaction shared.Transaction,
	ledgerTransaction *ledger.LedgerTransaction,
) error {
	if ledgerTransaction == nil {
		return errors.New("ledger transaction cannot be nil")
	}

	sessionCtx := transaction.SessionContext(ctx)

	result, err := repo.collection.ReplaceOne(
		sessionCtx,
		bson.M{"_id": ledgerTransaction.ID()},
		ledgerTransaction,
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return shared.ErrNotFound
	}

	return nil
}

func (repo *MongoLedgerTransactionRepository) Delete(
	ctx context.Context,
	transaction shared.Transaction,
	ledgerTransactionID shared.LedgerTransactionID,
) error {
	sessionCtx := transaction.SessionContext(ctx)

	result, err := repo.collection.DeleteOne(sessionCtx, bson.M{"_id": ledgerTransactionID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return shared.ErrNotFound
	}
	return nil
}

// -------------------- Query Repository --------------------

func (repo *MongoLedgerTransactionRepository) GetByID(
	ctx context.Context,
	ledgerTransactionID shared.LedgerTransactionID,
) (*ledger.LedgerTransaction, error) {
	var tx ledger.LedgerTransaction
	err := repo.collection.FindOne(ctx, bson.M{"_id": ledgerTransactionID}).Decode(&tx)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, shared.ErrNotFound
	}
	return &tx, err
}

func (repo *MongoLedgerTransactionRepository) GetByReference(
	ctx context.Context,
	referenceType ledger.ReferenceType,
	referenceID string,
) ([]*ledger.LedgerTransaction, error) {
	filter := bson.M{
		"referenceType": referenceType,
		"referenceID":   referenceID,
	}

	cursor, err := repo.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var transactions []*ledger.LedgerTransaction
	for cursor.Next(ctx) {
		var tx ledger.LedgerTransaction
		if err := cursor.Decode(&tx); err != nil {
			return nil, err
		}
		transactions = append(transactions, &tx)
	}

	return transactions, nil
}

func (repo *MongoLedgerTransactionRepository) GetLatestByReferenceAndAccount(
	ctx context.Context,
	accountID shared.AccountID,
	referenceID string,
) (*ledger.LedgerTransaction, error) {
	filter := bson.M{
		"referenceID":       referenceID,
		"entries.accountID": accountID,
	}

	opts := options.FindOne().SetSort(bson.D{{Key: "createdAt", Value: -1}})

	var tx ledger.LedgerTransaction
	err := repo.collection.FindOne(ctx, filter, opts).Decode(&tx)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, shared.ErrNotFound
	}

	return &tx, err
}
