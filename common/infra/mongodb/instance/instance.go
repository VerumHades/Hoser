package mongodb

import (
	"common/pkg/instance"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoInstanceRepository implements instance.InstanceRepository using MongoDB.
type MongoInstanceRepository struct {
	collection *mongo.Collection
	timeout    time.Duration
}

// NewMongoInstanceRepository creates a new MongoInstanceRepository.
func NewMongoInstanceRepository(collection *mongo.Collection, timeout time.Duration) *MongoInstanceRepository {
	return &MongoInstanceRepository{
		collection: collection,
		timeout:    timeout,
	}
}

// Save upserts an instance into MongoDB.
func (r *MongoInstanceRepository) Save(inst *instance.Instance) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	filter := bson.M{"id": inst.ID}
	update := bson.M{"$set": inst}
	opts := options.Update().SetUpsert(true)

	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

// GetByID fetches an instance by its ID.
func (r *MongoInstanceRepository) GetByID(id string) (*instance.Instance, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	filter := bson.M{"id": id}
	var inst instance.Instance
	err := r.collection.FindOne(ctx, filter).Decode(&inst)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &inst, err
}

// Delete removes an instance by its ID.
func (r *MongoInstanceRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	_, err := r.collection.DeleteOne(ctx, bson.M{"id": id})
	return err
}

// ListByListing returns all instances with the given listing ID.
func (r *MongoInstanceRepository) ListByListing(listingID string) ([]*instance.Instance, error) {
	return r.listByFilter(bson.M{"listingid": listingID})
}

// ListByBilling returns all instances with the given billing ID.
func (r *MongoInstanceRepository) ListByBilling(billingID string) ([]*instance.Instance, error) {
	return r.listByFilter(bson.M{"billingid": billingID})
}

// ListExpired returns instances whose expiry is before the cutoff time.
func (r *MongoInstanceRepository) ListExpired(cutoff time.Time) ([]*instance.Instance, error) {
	return r.listByFilter(bson.M{"expiry": bson.M{"$lt": cutoff}})
}

// ListByState returns instances matching the given state and contract state.
func (r *MongoInstanceRepository) ListByState(state instance.InstanceState, contractState instance.ContractState) ([]*instance.Instance, error) {
	filter := bson.M{
		"state":         state,
		"contractstate": contractState,
	}
	return r.listByFilter(filter)
}

// ListByPendingHardwareUpdate returns instances with pending hardware updates.
func (r *MongoInstanceRepository) ListByPendingHardwareUpdate() ([]*instance.Instance, error) {
	return r.listByFilter(bson.M{"updatehardware": true})
}

// listByFilter is a helper for queries returning multiple instances.
func (r *MongoInstanceRepository) listByFilter(filter bson.M) ([]*instance.Instance, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var instances []*instance.Instance
	for cursor.Next(ctx) {
		var inst instance.Instance
		if err := cursor.Decode(&inst); err != nil {
			return nil, err
		}
		instances = append(instances, &inst)
	}

	return instances, nil
}
