package mongodbledger

import (
	"context"
	"errors"
	"time"

	"common/pkg/domain/ledger"
	mongodbregistry "common/pkg/infrastructure/mongodb/registry"
	"common/pkg/shared"
	"common/pkg/util"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoLedgerTransactionRepository struct {
	collection *mongo.Collection
}

func NewMongoLedgerTransactionRepository(reg *mongodbregistry.DatabaseRegistry) *MongoLedgerTransactionRepository {
	if reg.LedgerTransactions == nil {
		panic("ledger transactions collection must not be nil")
	}
	return &MongoLedgerTransactionRepository{collection: reg.LedgerTransactions}
}

// EnsureIndexes creates MongoDB indexes for ledger transactions
func (repo *MongoLedgerTransactionRepository) EnsureIndexes(ctx context.Context) error {
	indexModels := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "reference_type", Value: 1},
				{Key: "reference_id", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "entries.account_id", Value: 1},
				{Key: "reference_id", Value: 1},
			},
		},
	}

	_, err := repo.collection.Indexes().CreateMany(ctx, indexModels)
	return err
}

// -------------------- Document Mapping --------------------

type ledgerEntryDocument struct {
	AccountID          shared.AccountID `bson:"account_id"`
	AmountInMinorUnits int64            `bson:"amount_in_minor_units"`
	CreatedAt          int64            `bson:"created_at"` // nanoseconds
}

type ledgerTransactionDocument struct {
	ID            shared.LedgerTransactionID `bson:"_id"`
	ReferenceType ledger.ReferenceType       `bson:"reference_type"`
	ReferenceID   string                     `bson:"reference_id"`
	CreatedAt     int64                      `bson:"created_at"` // nanoseconds
	Entries       []ledgerEntryDocument      `bson:"entries"`
}

func ledgerMapEntityToDocument(entity *ledger.LedgerTransaction) *ledgerTransactionDocument {
	entries := make([]ledgerEntryDocument, len(entity.Entries()))
	for i, e := range entity.Entries() {
		entries[i] = ledgerEntryDocument{
			AccountID:          e.AccountID(),
			AmountInMinorUnits: e.AmountInMinorUnits(),
			CreatedAt:          e.CreatedAt().UnixNano(),
		}
	}

	return &ledgerTransactionDocument{
		ID:            entity.ID(),
		ReferenceType: entity.ReferenceType(),
		ReferenceID:   entity.ReferenceID(),
		CreatedAt:     entity.CreatedAt().UnixNano(),
		Entries:       entries,
	}
}

func ledgerMapDocumentToEntity(doc *ledgerTransactionDocument) (*ledger.LedgerTransaction, error) {
	entries := make([]*ledger.LedgerEntry, len(doc.Entries))
	for i, e := range doc.Entries {
		entry, err := ledger.NewLedgerEntryWithTime(
			e.AccountID,
			e.AmountInMinorUnits,
			time.Unix(0, e.CreatedAt),
		)
		if err != nil {
			return nil, err
		}
		entries[i] = entry
	}

	return ledger.NewLedgerTransactionWithID(
		doc.ID,
		doc.ReferenceType,
		doc.ReferenceID,
		entries,
		time.Unix(0, doc.CreatedAt),
	)
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

	operationCtx := util.ResolveTransactionalContext(ctx, transaction)
	doc := ledgerMapEntityToDocument(ledgerTransaction)

	_, err := repo.collection.InsertOne(operationCtx, doc)
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

	operationCtx := util.ResolveTransactionalContext(ctx, transaction)
	doc := ledgerMapEntityToDocument(ledgerTransaction)

	result, err := repo.collection.ReplaceOne(
		operationCtx,
		bson.M{"_id": ledgerTransaction.ID()},
		doc,
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
	operationCtx := util.ResolveTransactionalContext(ctx, transaction)
	result, err := repo.collection.DeleteOne(operationCtx, bson.M{"_id": ledgerTransactionID})
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
	var doc ledgerTransactionDocument
	err := repo.collection.FindOne(ctx, bson.M{"_id": ledgerTransactionID}).Decode(&doc)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, shared.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return ledgerMapDocumentToEntity(&doc)
}

func (repo *MongoLedgerTransactionRepository) GetByReference(
	ctx context.Context,
	referenceType ledger.ReferenceType,
	referenceID string,
) ([]*ledger.LedgerTransaction, error) {
	filter := bson.M{
		"reference_type": referenceType,
		"reference_id":   referenceID,
	}

	cursor, err := repo.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var transactions []*ledger.LedgerTransaction
	for cursor.Next(ctx) {
		var doc ledgerTransactionDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		entity, err := ledgerMapDocumentToEntity(&doc)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, entity)
	}

	return transactions, nil
}

func (repo *MongoLedgerTransactionRepository) GetLatestByReferenceAndAccount(
	ctx context.Context,
	accountID shared.AccountID,
	referenceID string,
) (*ledger.LedgerTransaction, error) {
	filter := bson.M{
		"reference_id":       referenceID,
		"entries.account_id": accountID,
	}

	opts := options.FindOne().SetSort(bson.D{{Key: "created_at", Value: -1}})
	var doc ledgerTransactionDocument
	err := repo.collection.FindOne(ctx, filter, opts).Decode(&doc)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, shared.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return ledgerMapDocumentToEntity(&doc)
}
