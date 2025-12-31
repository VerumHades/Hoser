package mongodbrates

import (
	"context"
	"errors"
	"time"

	"common/pkg/domain/rates"
	"common/pkg/shared"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoHardwareCostRateRepository implements command and query repositories for hardware cost rates.
type MongoHardwareCostRateRepository struct {
	collection *mongo.Collection
}

// NewMongoHardwareCostRateRepository constructs a repository backed by a Mongo collection.
func NewMongoHardwareCostRateRepository(collection *mongo.Collection) *MongoHardwareCostRateRepository {
	if collection == nil {
		panic("mongo collection must not be nil")
	}
	return &MongoHardwareCostRateRepository{collection: collection}
}

// EnsureIndexes creates required indexes for hardware cost rates.
func (repo *MongoHardwareCostRateRepository) EnsureIndexes(ctx context.Context) error {
	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "resource", Value: 1},
				{Key: "validFrom", Value: 1},
			},
		},
	}
	_, err := repo.collection.Indexes().CreateMany(ctx, indexModels)
	return err
}

// -------------------- Command Repository --------------------

func (repo *MongoHardwareCostRateRepository) Create(
	ctx context.Context,
	transaction shared.Transaction,
	rate *rates.HardwareCostRate,
) (*rates.HardwareCostRate, error) {
	if rate == nil {
		return nil, errors.New("hardware cost rate cannot be nil")
	}
	sessionCtx := transaction.SessionContext(ctx)
	_, err := repo.collection.InsertOne(sessionCtx, rate)
	if mongo.IsDuplicateKeyError(err) {
		return nil, shared.ErrAlreadyExists
	}
	return rate, err
}

func (repo *MongoHardwareCostRateRepository) Update(
	ctx context.Context,
	transaction shared.Transaction,
	rate *rates.HardwareCostRate,
) (*rates.HardwareCostRate, error) {
	if rate == nil {
		return nil, errors.New("hardware cost rate cannot be nil")
	}
	sessionCtx := transaction.SessionContext(ctx)
	result, err := repo.collection.ReplaceOne(sessionCtx, bson.M{"_id": rate.ID()}, rate)
	if err != nil {
		return nil, err
	}
	if result.MatchedCount == 0 {
		return nil, shared.ErrNotFound
	}
	return rate, nil
}

// -------------------- Query Repository --------------------

func (repo *MongoHardwareCostRateRepository) GetActiveRate(
	ctx context.Context,
	resourceType rates.HardwareResourceType,
	at time.Time,
) (*rates.HardwareCostRate, error) {
	filter := bson.M{
		"resource":  resourceType,
		"validFrom": bson.M{"$lte": at},
	}
	opts := options.FindOne().SetSort(bson.D{{Key: "validFrom", Value: -1}})
	var rate rates.HardwareCostRate
	err := repo.collection.FindOne(ctx, filter, opts).Decode(&rate)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, shared.ErrNotFound
	}
	return &rate, err
}

func (repo *MongoHardwareCostRateRepository) FetchNextBatchOrderedByEffectiveDate(
	ctx context.Context,
	request shared.BatchRequest[rates.HardwareCostRateCursor],
) ([]*rates.HardwareCostRate, rates.HardwareCostRateCursor, error) {
	filter := bson.M{}
	filter["validFrom"] = bson.M{"$gt": request.Cursor.LastEffectiveDate}

	findOptions := options.Find().
		SetSort(bson.D{{Key: "validFrom", Value: 1}}).
		SetLimit(int64(request.MaxBatchSize))

	cursor, err := repo.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, rates.HardwareCostRateCursor{}, err
	}
	defer cursor.Close(ctx)

	var ratesList []*rates.HardwareCostRate
	for cursor.Next(ctx) {
		var r rates.HardwareCostRate
		if err := cursor.Decode(&r); err != nil {
			return nil, rates.HardwareCostRateCursor{}, err
		}
		ratesList = append(ratesList, &r)
	}

	if len(ratesList) == 0 {
		return ratesList, rates.HardwareCostRateCursor{}, nil
	}

	last := ratesList[len(ratesList)-1]
	return ratesList, rates.HardwareCostRateCursor{
		LastEffectiveDate: last.ValidFrom(),
	}, nil
}
