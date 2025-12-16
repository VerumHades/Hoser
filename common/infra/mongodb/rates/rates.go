package ratesmongodb

import (
	"context"
	"errors"
	"time"

	"common/pkg/money"
	"common/pkg/rates"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// hardwareCostRateDocument represents the MongoDB persistence model.
type hardwareCostRateDocument struct {
	ID        string      `bson:"_id"`
	CPUCost   money.Money `bson:"cpu_cost"`
	RAMCost   money.Money `bson:"ram_cost"`
	DiskCost  money.Money `bson:"disk_cost"`
	ValidFrom time.Time   `bson:"valid_from"`
	ValidTo   *time.Time  `bson:"valid_to,omitempty"`
}

// HardwareCostRepository implements rates.HardwareCostRepository using MongoDB.
type HardwareCostRepository struct {
	collection *mongo.Collection
}

// NewHardwareCostRepository creates a new MongoDB-backed hardware cost repository.
func NewHardwareCostRepository(collection *mongo.Collection) *HardwareCostRepository {
	return &HardwareCostRepository{collection: collection}
}

// Save inserts or updates a hardware cost rate.
func (r *HardwareCostRepository) Save(rate *rates.HardwareCostRate) error {
	doc := hardwareCostRateDocument{
		ID:        rate.ID,
		CPUCost:   rate.CPUCost,
		RAMCost:   rate.RAMCost,
		DiskCost:  rate.DiskCost,
		ValidFrom: rate.ValidFrom,
		ValidTo:   rate.ValidTo,
	}

	_, err := r.collection.UpdateByID(
		context.Background(),
		doc.ID,
		bson.M{"$set": doc},
		options.Update().SetUpsert(true),
	)
	return err
}

// GetActiveRate returns the rate active at the specified time.
func (r *HardwareCostRepository) GetActiveRate(at time.Time) (*rates.HardwareCostRate, error) {
	filter := bson.M{
		"valid_from": bson.M{"$lte": at},
		"$or": []bson.M{
			{"valid_to": bson.M{"$gt": at}},
			{"valid_to": nil},
		},
	}

	var doc hardwareCostRateDocument
	err := r.collection.FindOne(context.Background(), filter).Decode(&doc)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return mapDocumentToDomain(doc), nil
}

// ListAll returns all hardware cost rates.
func (r *HardwareCostRepository) ListAll() ([]*rates.HardwareCostRate, error) {
	cursor, err := r.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var results []*rates.HardwareCostRate
	for cursor.Next(context.Background()) {
		var doc hardwareCostRateDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		results = append(results, mapDocumentToDomain(doc))
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// mapDocumentToDomain converts a MongoDB document into a domain HardwareCostRate.
func mapDocumentToDomain(doc hardwareCostRateDocument) *rates.HardwareCostRate {
	return &rates.HardwareCostRate{
		ID:        doc.ID,
		CPUCost:   doc.CPUCost,
		RAMCost:   doc.RAMCost,
		DiskCost:  doc.DiskCost,
		ValidFrom: doc.ValidFrom,
		ValidTo:   doc.ValidTo,
	}
}
