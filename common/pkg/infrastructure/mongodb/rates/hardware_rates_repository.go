package mongodbrates

import (
	"context"
	"errors"
	"time"

	"common/pkg/domain/entities/rates"
	"common/pkg/domain/repositories"
	mongodbregistry "common/pkg/infrastructure/mongodb/registry"
	"common/pkg/shared"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoHardwareCostRateRepository struct {
	collection *mongo.Collection
}

func NewMongoHardwareCostRateRepository(databaseRegistry *mongodbregistry.DatabaseRegistry) *MongoHardwareCostRateRepository {
	if databaseRegistry.HardwareRates == nil {
		panic("hardware rates collection must not be nil")
	}
	return &MongoHardwareCostRateRepository{collection: databaseRegistry.HardwareRates}
}

func (repo *MongoHardwareCostRateRepository) EnsureIndexes(ctx context.Context) error {
	indexModels := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "resource", Value: 1},
				{Key: "valid_from", Value: 1},
			},
		},
	}

	_, err := repo.collection.Indexes().CreateMany(ctx, indexModels)
	return err
}

type hardwareCostRateDocument struct {
	ID          shared.HardwareCostRateID  `bson:"_id"`
	Resource    rates.HardwareResourceType `bson:"resource"`
	ValidFrom   int64                      `bson:"valid_from"`
	CostInCents int64                      `bson:"cost_in_cents"`
	CreatedAt   int64                      `bson:"created_at"`
}

func mapEntityToDocument(entity *rates.HardwareCostRate) *hardwareCostRateDocument {
	return &hardwareCostRateDocument{
		ID:          entity.ID(),
		Resource:    entity.Resource(),
		ValidFrom:   entity.ValidFrom().UnixNano(),
		CostInCents: entity.CostInCents(),
		CreatedAt:   entity.CreatedAt().UnixNano(),
	}
}

func mapDocumentToEntity(doc *hardwareCostRateDocument) (*rates.HardwareCostRate, error) {
	return rates.NewHardwareCostRateWithID(
		doc.ID,
		doc.Resource,
		doc.CostInCents,
		time.Unix(0, doc.ValidFrom),
		time.Unix(0, doc.CreatedAt),
	)
}

// -------------------- Command Repository --------------------

func (repo *MongoHardwareCostRateRepository) Create(
	ctx context.Context,

	rate *rates.HardwareCostRate,
) (*rates.HardwareCostRate, error) {
	if rate == nil {
		return nil, errors.New("hardware cost rate cannot be nil")
	}

	document := mapEntityToDocument(rate)

	_, err := repo.collection.InsertOne(ctx, document)
	if mongo.IsDuplicateKeyError(err) {
		return nil, shared.ErrAlreadyExists
	}

	return rate, err
}

func (repo *MongoHardwareCostRateRepository) Update(
	ctx context.Context,

	rate *rates.HardwareCostRate,
) (*rates.HardwareCostRate, error) {
	if rate == nil {
		return nil, errors.New("hardware cost rate cannot be nil")
	}

	document := mapEntityToDocument(rate)

	result, err := repo.collection.ReplaceOne(
		ctx,
		bson.M{"_id": rate.ID()},
		document,
	)
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
		"resource":   resourceType,
		"valid_from": bson.M{"$lte": at.UnixNano()},
	}
	findOptions := options.FindOne().SetSort(bson.D{{Key: "valid_from", Value: -1}})

	var doc hardwareCostRateDocument
	err := repo.collection.FindOne(ctx, filter, findOptions).Decode(&doc)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, shared.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return mapDocumentToEntity(&doc)
}

func (repo *MongoHardwareCostRateRepository) FetchNextBatchOrderedByEffectiveDate(
	ctx context.Context,
	request shared.BatchRequest[repositories.HardwareCostRateCursor],
) ([]*rates.HardwareCostRate, repositories.HardwareCostRateCursor, error) {
	filter := bson.M{}

	filter["valid_from"] = bson.M{
		"$gt": request.Cursor.LastEffectiveDate.UnixNano(),
	}

	findOptions := options.Find().
		SetSort(bson.D{{Key: "valid_from", Value: 1}}).
		SetLimit(int64(request.MaxBatchSize))

	cursor, err := repo.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, repositories.HardwareCostRateCursor{}, err
	}
	defer cursor.Close(ctx)

	var items []*rates.HardwareCostRate
	for cursor.Next(ctx) {
		var doc hardwareCostRateDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, repositories.HardwareCostRateCursor{}, err
		}

		entity, err := mapDocumentToEntity(&doc)
		if err != nil {
			return nil, repositories.HardwareCostRateCursor{}, err
		}

		items = append(items, entity)
	}

	if len(items) == 0 {
		return items, repositories.HardwareCostRateCursor{}, nil
	}

	last := items[len(items)-1]

	return items, repositories.HardwareCostRateCursor{
		LastEffectiveDate: last.ValidFrom(),
	}, nil
}
