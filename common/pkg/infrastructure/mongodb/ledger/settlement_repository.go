package mongodbledger

import (
	"context"
	"errors"
	"time"

	"common/pkg/domain/entities/accounting"
	mongodbregistry "common/pkg/infrastructure/mongodb/registry"
	"common/pkg/shared"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoSettlementRepository struct {
	collection *mongo.Collection
}

func NewMongoSettlementRepository(reg *mongodbregistry.DatabaseRegistry) *MongoSettlementRepository {
	if reg.Settlements == nil {
		panic("settlements collection must not be nil")
	}
	return &MongoSettlementRepository{collection: reg.Settlements}
}

// EnsureIndexes creates MongoDB indexes for settlements
func (repo *MongoSettlementRepository) EnsureIndexes(ctx context.Context) error {
	indexModels := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "ledger_transaction_id", Value: 1}},
		},
		{
			Keys: bson.D{
				{Key: "ledger_transaction_id", Value: 1},
				{Key: "created_at", Value: -1},
			},
		},
	}
	_, err := repo.collection.Indexes().CreateMany(ctx, indexModels)
	return err
}

// -------------------- Document Mapping --------------------

type settlementDocument struct {
	ID                  shared.SettlementID         `bson:"_id"`
	LedgerTransactionID shared.LedgerTransactionID  `bson:"ledger_transaction_id"`
	AccountID           shared.AccountID            `bson:"account_id"`
	AmountInMinorUnits  int64                       `bson:"amount_in_minor_units"`
	Status              accounting.SettlementStatus `bson:"status"`
	ReferenceID         string                      `bson:"reference_id"`
	CreatedAt           int64                       `bson:"created_at"` // nanoseconds
	UpdatedAt           int64                       `bson:"updated_at"` // nanoseconds
}

func settlementMapEntityToDocument(entity *accounting.Settlement) *settlementDocument {
	return &settlementDocument{
		ID:                  entity.ID(),
		LedgerTransactionID: entity.LedgerTransactionID(),
		AccountID:           entity.AccountID(),
		AmountInMinorUnits:  entity.Amount(),
		Status:              entity.Status(),
		ReferenceID:         entity.ReferenceID(),
		CreatedAt:           entity.CreatedAt().UnixNano(),
		UpdatedAt:           entity.UpdatedAt().UnixNano(),
	}
}

func settlementMapDocumentToEntity(doc *settlementDocument) *accounting.Settlement {
	return accounting.NewSettlementWithID(
		doc.ID,
		doc.LedgerTransactionID,
		doc.AccountID,
		doc.AmountInMinorUnits,
		doc.Status,
		doc.ReferenceID,
		time.Unix(0, doc.CreatedAt),
		time.Unix(0, doc.UpdatedAt),
	)
}

// -------------------- Command Repository --------------------

func (repo *MongoSettlementRepository) Create(
	ctx context.Context,

	settlement *accounting.Settlement,
) (*accounting.Settlement, error) {
	if settlement == nil {
		return nil, errors.New("settlement cannot be nil")
	}

	doc := settlementMapEntityToDocument(settlement)

	_, err := repo.collection.InsertOne(ctx, doc)
	if mongo.IsDuplicateKeyError(err) {
		return nil, shared.ErrAlreadyExists
	}

	return settlement, err
}

func (repo *MongoSettlementRepository) Update(
	ctx context.Context,

	settlement *accounting.Settlement,
) (*accounting.Settlement, error) {
	if settlement == nil {
		return nil, errors.New("settlement cannot be nil")
	}

	doc := settlementMapEntityToDocument(settlement)

	result, err := repo.collection.ReplaceOne(
		ctx,
		bson.M{"_id": settlement.ID()},
		doc,
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

	settlementID shared.SettlementID,
) error {

	result, err := repo.collection.DeleteOne(ctx, bson.M{"_id": settlementID})
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
) (*accounting.Settlement, error) {
	var doc settlementDocument
	err := repo.collection.FindOne(ctx, bson.M{"_id": settlementID}).Decode(&doc)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, shared.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return settlementMapDocumentToEntity(&doc), nil
}

func (repo *MongoSettlementRepository) GetByLedgerTransactionID(
	ctx context.Context,
	ledgerTransactionID shared.LedgerTransactionID,
) ([]*accounting.Settlement, error) {
	filter := bson.M{"ledger_transaction_id": ledgerTransactionID}

	cursor, err := repo.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var settlements []*accounting.Settlement
	for cursor.Next(ctx) {
		var doc settlementDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		settlements = append(settlements, settlementMapDocumentToEntity(&doc))
	}

	return settlements, nil
}

func (repo *MongoSettlementRepository) GetLastByLedgerTransactionID(
	ctx context.Context,
	ledgerTransactionID shared.LedgerTransactionID,
) (*accounting.Settlement, error) {
	filter := bson.M{"ledger_transaction_id": ledgerTransactionID}
	opts := options.FindOne().SetSort(bson.D{{Key: "created_at", Value: -1}})

	var doc settlementDocument
	err := repo.collection.FindOne(ctx, filter, opts).Decode(&doc)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, shared.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return settlementMapDocumentToEntity(&doc), nil
}
