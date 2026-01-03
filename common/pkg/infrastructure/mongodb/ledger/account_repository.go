package mongodbledger

import (
	"context"
	"errors"
	"time"

	mongodbregistry "common/pkg/infrastructure/mongodb/registry"
	"common/pkg/shared"
	"common/pkg/util"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoAccountRepository struct {
	collection *mongo.Collection
}

func NewMongoAccountRepository(reg *mongodbregistry.DatabaseRegistry) *MongoAccountRepository {
	if reg.Accounts == nil {
		panic("accounts collection must not be nil")
	}
	return &MongoAccountRepository{collection: reg.Accounts}
}

// EnsureIndexes creates MongoDB indexes for accounts.
func (repo *MongoAccountRepository) EnsureIndexes(ctx context.Context) error {
	indexModels := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "owner_type", Value: 1},
				{Key: "owner_id", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "account_type", Value: 1},
			},
		},
	}
	_, err := repo.collection.Indexes().CreateMany(ctx, indexModels)
	return err
}

// -------------------- Document Mapping --------------------

type accountDocument struct {
	ID          shared.AccountID            `bson:"_id"`
	AccountType accounting.AccountType      `bson:"account_type"`
	OwnerType   accounting.AccountOwnerType `bson:"owner_type"`
	OwnerID     string                      `bson:"owner_id"`
	CreatedAt   int64                       `bson:"created_at"` // nanoseconds
}

func mapEntityToDocument(entity *accounting.Account) *accountDocument {
	return &accountDocument{
		ID:          entity.ID(),
		AccountType: entity.Type(),
		OwnerType:   entity.OwnerType(),
		OwnerID:     entity.OwnerID(),
		CreatedAt:   entity.CreatedAt().UnixNano(),
	}
}

func mapDocumentToEntity(doc *accountDocument) (*accounting.Account, error) {
	return accounting.NewAccountWithID(
		doc.ID,
		doc.AccountType,
		doc.OwnerType,
		doc.OwnerID,
		time.Unix(0, doc.CreatedAt),
	)
}

// -------------------- Command Repository --------------------

func (repo *MongoAccountRepository) Create(
	ctx context.Context,

	accountEntity *accounting.Account,
) (*accounting.Account, error) {
	if accountEntity == nil {
		return nil, errors.New("account cannot be nil")
	}

	operationCtx := util.ResolveTransactionalContext(ctx, transaction)
	doc := mapEntityToDocument(accountEntity)

	_, err := repo.collection.InsertOne(operationCtx, doc)
	if mongo.IsDuplicateKeyError(err) {
		return nil, shared.ErrAlreadyExists
	}
	return accountEntity, err
}

func (repo *MongoAccountRepository) Update(
	ctx context.Context,

	accountEntity *accounting.Account,
) (*accounting.Account, error) {
	if accountEntity == nil {
		return nil, errors.New("account cannot be nil")
	}

	operationCtx := util.ResolveTransactionalContext(ctx, transaction)
	doc := mapEntityToDocument(accountEntity)

	result, err := repo.collection.ReplaceOne(
		operationCtx,
		bson.M{"_id": accountEntity.ID()},
		doc,
	)
	if err != nil {
		return nil, err
	}
	if result.MatchedCount == 0 {
		return nil, shared.ErrNotFound
	}
	return accountEntity, nil
}

func (repo *MongoAccountRepository) Delete(
	ctx context.Context,

	accountID shared.AccountID,
) error {
	operationCtx := util.ResolveTransactionalContext(ctx, transaction)
	result, err := repo.collection.DeleteOne(operationCtx, bson.M{"_id": accountID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return shared.ErrNotFound
	}
	return nil
}

// -------------------- Query Repository --------------------

func (repo *MongoAccountRepository) GetByID(
	ctx context.Context,
	accountID shared.AccountID,
) (*accounting.Account, error) {
	var doc accountDocument
	err := repo.collection.FindOne(ctx, bson.M{"_id": accountID}).Decode(&doc)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, shared.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return mapDocumentToEntity(&doc)
}

func (repo *MongoAccountRepository) GetByOwner(
	ctx context.Context,
	ownerType accounting.AccountOwnerType,
	ownerID string,
) ([]*accounting.Account, error) {
	filter := bson.M{"owner_type": ownerType, "owner_id": ownerID}
	cursor, err := repo.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var accounts []*accounting.Account
	for cursor.Next(ctx) {
		var doc accountDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		entity, err := mapDocumentToEntity(&doc)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, entity)
	}

	return accounts, nil
}

func (repo *MongoAccountRepository) GetFirstByOwner(
	ctx context.Context,
	ownerType accounting.AccountOwnerType,
	ownerID string,
) (*accounting.Account, error) {
	filter := bson.M{"owner_type": ownerType, "owner_id": ownerID}
	var doc accountDocument
	err := repo.collection.FindOne(ctx, filter).Decode(&doc)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, shared.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return mapDocumentToEntity(&doc)
}

func (repo *MongoAccountRepository) GetByType(
	ctx context.Context,
	accountType accounting.AccountType,
) ([]*accounting.Account, error) {
	filter := bson.M{"account_type": accountType}
	cursor, err := repo.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var accounts []*accounting.Account
	for cursor.Next(ctx) {
		var doc accountDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		entity, err := mapDocumentToEntity(&doc)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, entity)
	}

	return accounts, nil
}
