package mongodbinstance

import (
	"context"
	"errors"
	"time"

	"common/pkg/domain/instance"
	mongodbregistry "common/pkg/infrastructure/mongodb/registry"
	"common/pkg/shared"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoInstanceRentalContractRepository implements both command and query repositories.
type MongoInstanceRentalContractRepository struct {
	collection *mongo.Collection
}

// NewMongoInstanceRentalContractRepository constructs the repository.
func NewMongoInstanceRentalContractRepository(databaseRegistry *mongodbregistry.DatabaseRegistry) *MongoInstanceRentalContractRepository {
	if databaseRegistry.InstanceContracts == nil {
		panic("instance contracts collection must not be nil")
	}
	return &MongoInstanceRentalContractRepository{collection: databaseRegistry.InstanceContracts}
}

// EnsureIndexes creates the required MongoDB indexes.
func (repo *MongoInstanceRentalContractRepository) EnsureIndexes(ctx context.Context) error {
	indexModels := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "listingID", Value: 1},
				{Key: "createdAt", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "ownerID", Value: 1},
				{Key: "createdAt", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "periodStart", Value: 1},
				{Key: "periodEnd", Value: 1},
				{Key: "createdAt", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "renewalDueAt", Value: 1},
				{Key: "createdAt", Value: 1},
			},
		},
	}

	_, err := repo.collection.Indexes().CreateMany(ctx, indexModels)
	return err
}

// -------------------- Command Repository --------------------

func (repo *MongoInstanceRentalContractRepository) Create(
	ctx context.Context,
	contract *instance.InstanceRentalContract,
) error {
	if contract == nil {
		return errors.New("contract cannot be nil")
	}

	sessionCtx := transaction.SessionContext(ctx)
	_, err := repo.collection.InsertOne(sessionCtx, contract)
	if mongo.IsDuplicateKeyError(err) {
		return shared.ErrAlreadyExists
	}
	return err
}

func (repo *MongoInstanceRentalContractRepository) Update(
	ctx context.Context,
	contract *instance.InstanceRentalContract,
) error {
	if contract == nil {
		return errors.New("contract cannot be nil")
	}

	sessionCtx := transaction.SessionContext(ctx)
	result, err := repo.collection.ReplaceOne(sessionCtx, bson.M{"_id": contract.ID()}, contract)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return shared.ErrNotFound
	}
	return nil
}

func (repo *MongoInstanceRentalContractRepository) Delete(
	ctx context.Context,

	contractID shared.InstanceRentalContractID,
) error {
	sessionCtx := transaction.SessionContext(ctx)
	result, err := repo.collection.DeleteOne(sessionCtx, bson.M{"_id": contractID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return shared.ErrNotFound
	}
	return nil
}

// -------------------- Query Repository --------------------

func (repo *MongoInstanceRentalContractRepository) GetByID(
	ctx context.Context,
	contractID shared.InstanceRentalContractID,
) (*instance.InstanceRentalContract, error) {
	var contract instance.InstanceRentalContract
	err := repo.collection.FindOne(ctx, bson.M{"_id": contractID}).Decode(&contract)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, shared.ErrNotFound
	}
	return &contract, err
}

func (repo *MongoInstanceRentalContractRepository) Exists(
	ctx context.Context,
	contractID shared.InstanceRentalContractID,
) (bool, error) {
	count, err := repo.collection.CountDocuments(ctx, bson.M{"_id": contractID}, options.Count().SetLimit(1))
	return count == 1, err
}

func (repo *MongoInstanceRentalContractRepository) FetchNextBatchByListing(
	ctx context.Context,
	listingID shared.ListingID,
	request shared.BatchRequest[instance.InstanceRentalContractCursor],
) ([]*instance.InstanceRentalContract, instance.InstanceRentalContractCursor, error) {
	filter := bson.M{"listingID": listingID}
	return repo.fetchBatch(ctx, filter, request)
}

func (repo *MongoInstanceRentalContractRepository) FetchNextBatchByOwner(
	ctx context.Context,
	ownerID shared.UserID,
	request shared.BatchRequest[instance.InstanceRentalContractCursor],
) ([]*instance.InstanceRentalContract, instance.InstanceRentalContractCursor, error) {
	filter := bson.M{"ownerID": ownerID}
	return repo.fetchBatch(ctx, filter, request)
}

func (repo *MongoInstanceRentalContractRepository) FetchNextBatchActiveAt(
	ctx context.Context,
	at time.Time,
	request shared.BatchRequest[instance.InstanceRentalContractCursor],
) ([]*instance.InstanceRentalContract, instance.InstanceRentalContractCursor, error) {
	filter := bson.M{
		"periodStart": bson.M{"$lte": at},
		"periodEnd":   bson.M{"$gt": at},
	}
	return repo.fetchBatch(ctx, filter, request)
}

func (repo *MongoInstanceRentalContractRepository) FetchNextBatchExpiredBefore(
	ctx context.Context,
	cutoffTime time.Time,
	request shared.BatchRequest[instance.InstanceRentalContractCursor],
) ([]*instance.InstanceRentalContract, instance.InstanceRentalContractCursor, error) {
	filter := bson.M{"periodEnd": bson.M{"$lt": cutoffTime}}
	return repo.fetchBatch(ctx, filter, request)
}

func (repo *MongoInstanceRentalContractRepository) FetchNextBatchPendingRenewal(
	ctx context.Context,
	cutoffTime time.Time,
	request shared.BatchRequest[instance.InstanceRentalContractCursor],
) ([]*instance.InstanceRentalContract, instance.InstanceRentalContractCursor, error) {
	filter := bson.M{"renewalDueAt": bson.M{"$lte": cutoffTime}}
	return repo.fetchBatch(ctx, filter, request)
}

// -------------------- Internal --------------------

func (repo *MongoInstanceRentalContractRepository) fetchBatch(
	ctx context.Context,
	filter bson.M,
	request shared.BatchRequest[instance.InstanceRentalContractCursor],
) ([]*instance.InstanceRentalContract, instance.InstanceRentalContractCursor, error) {
	filter["createdAt"] = bson.M{"$gt": request.Cursor.LastCreatedAt}

	findOptions := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: 1}}).
		SetLimit(int64(request.MaxBatchSize))

	cursor, err := repo.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, instance.InstanceRentalContractCursor{}, err
	}
	defer cursor.Close(ctx)

	var contracts []*instance.InstanceRentalContract
	for cursor.Next(ctx) {
		var c instance.InstanceRentalContract
		if err := cursor.Decode(&c); err != nil {
			return nil, instance.InstanceRentalContractCursor{}, err
		}
		contracts = append(contracts, &c)
	}

	if len(contracts) == 0 {
		return contracts, instance.InstanceRentalContractCursor{}, nil
	}

	last := contracts[len(contracts)-1]
	return contracts, instance.InstanceRentalContractCursor{
		LastCreatedAt: last.CreatedAt(),
	}, nil
}
