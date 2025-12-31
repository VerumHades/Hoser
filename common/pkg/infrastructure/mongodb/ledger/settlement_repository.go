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

// MongoSettlementRepository implements SettlementCommandRepository and SettlementQueryRepository
type MongoSettlementRepository struct {
	collection *mongo.Collection
}

// NewMongoSettlementRepository constructs a repository backed by a Mongo collection
func NewMongoSettlementRepository(collection *mongo.Collection) *MongoSettlementRepository {
	if collection == nil {
		panic("mongo collection must not be nil")
	}
	return &MongoSettlementRepository{collection: collection}
}

// EnsureSettlementIndexes creates indexes for efficient queries
func (repo *MongoSettlementRepository) EnsureIndexes(ctx context.Context) error {
	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "ledgerTransactionID", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "ledgerTransactionID", Value: 1},
				{Key: "createdAt", Value: -1},
			},
		},
	}

	_, err := repo.collection.Indexes().CreateMany(ctx, indexModels)
	return err
}

// -------------------- Command Repository --------------------

func (repo *MongoSettlementRepository) Create(
	ctx context.Context,
	transaction shared.Transaction,
	settlement *ledger.Settlement,
) (*ledger.Settlement, error) {
	if settlement == nil {
		return nil, errors.New("settlement cannot be nil")
	}

	sessionCtx := transaction.SessionContext(ctx)

	_, err := repo.collection.InsertOne(sessionCtx, settlement)
	if mongo.IsDuplicateKeyError(err) {
		return nil, shared.ErrAlreadyExists
	}

	return settlement, err
}

func (repo *MongoSettlementRepository) Update(
	ctx context.Context,
	transaction shared.Transaction,
	settlement *ledger.Settlement,
) (*ledger.Settlement, error) {
	if settlement == nil {
		return nil, errors.New("settlement cannot be nil")
	}

	sessionCtx := transaction.SessionContext(ctx)

	result, err := repo.collection.ReplaceOne(
		sessionCtx,
		bson.M{"_id": settlement.ID()},
		settlement,
	)
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, shared.ErrNotFound
	}

	return settlement, nil
}

func (repo *MongoSettlementRepository) Delete(
	ctx context.Context,
	transaction shared.Transaction,
	settlementID shared.SettlementID,
) error {
	sessionCtx := transaction.SessionContext(ctx)

	result, err := repo.collection.DeleteOne(sessionCtx, bson.M{"_id": settlementID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return shared.ErrNotFound
	}

	return nil
}

// -------------------- Query Repository --------------------

func (repo *MongoSettlementRepository) GetByID(
	ctx context.Context,
	settlementID shared.SettlementID,
) (*ledger.Settlement, error) {
	var settlement ledger.Settlement
	err := repo.collection.FindOne(ctx, bson.M{"_id": settlementID}).Decode(&settlement)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, shared.ErrNotFound
	}

	return &settlement, err
}

func (repo *MongoSettlementRepository) GetByLedgerTransactionID(
	ctx context.Context,
	ledgerTransactionID shared.LedgerTransactionID,
) ([]*ledger.Settlement, error) {
	filter := bson.M{"ledgerTransactionID": ledgerTransactionID}

	cursor, err := repo.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var settlements []*ledger.Settlement
	for cursor.Next(ctx) {
		var s ledger.Settlement
		if err := cursor.Decode(&s); err != nil {
			return nil, err
		}
		settlements = append(settlements, &s)
	}

	return settlements, nil
}

func (repo *MongoSettlementRepository) GetLastByLedgerTransactionID(
	ctx context.Context,
	ledgerTransactionID shared.LedgerTransactionID,
) (*ledger.Settlement, error) {
	filter := bson.M{"ledgerTransactionID": ledgerTransactionID}

	opts := options.FindOne().SetSort(bson.D{{Key: "createdAt", Value: -1}})

	var settlement ledger.Settlement
	err := repo.collection.FindOne(ctx, filter, opts).Decode(&settlement)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, shared.ErrNotFound
	}

	return &settlement, err
}
