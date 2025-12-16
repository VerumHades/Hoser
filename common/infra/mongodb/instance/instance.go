package instancemongodb

import (
	"context"
	"errors"

	"common/pkg/hardware"
	"common/pkg/instance"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// instanceDocument represents the MongoDB persistence model for Instance.
type instanceDocument struct {
	ID                    string                          `bson:"_id"`
	ListingID             string                          `bson:"listing_id"`
	BillingID             string                          `bson:"billing_id"`
	State                 instance.InstanceState          `bson:"state"`
	HardwareSpecification *hardware.HardwareSpecification `bson:"hardware_spec,omitempty"`
}

// InstanceRepository implements instance.InstanceRepository using MongoDB.
type InstanceRepository struct {
	collection *mongo.Collection
}

// NewInstanceRepository creates a new MongoDB-backed instance repository.
func NewInstanceRepository(collection *mongo.Collection) *InstanceRepository {
	return &InstanceRepository{
		collection: collection,
	}
}

// Save inserts or updates an instance.
func (r *InstanceRepository) Save(instanceEntity *instance.Instance) error {
	doc := instanceDocument{
		ID:                    instanceEntity.ID,
		ListingID:             instanceEntity.ListingID,
		BillingID:             instanceEntity.BillingID,
		State:                 instanceEntity.State,
		HardwareSpecification: instanceEntity.HardwareSpecification,
	}

	_, err := r.collection.UpdateByID(
		context.Background(),
		doc.ID,
		bson.M{"$set": doc},
		options.Update().SetUpsert(true),
	)
	return err
}

// GetByID retrieves an instance by ID.
func (r *InstanceRepository) GetByID(id string) (*instance.Instance, error) {
	var doc instanceDocument
	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&doc)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return mapDocumentToDomain(doc), nil
}

// Delete removes an instance by ID.
func (r *InstanceRepository) Delete(id string) error {
	_, err := r.collection.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}

// ListByListing returns all instances for a given listing.
func (r *InstanceRepository) ListByListing(listingID string) ([]*instance.Instance, error) {
	cursor, err := r.collection.Find(context.Background(), bson.M{"listing_id": listingID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var results []*instance.Instance
	for cursor.Next(context.Background()) {
		var doc instanceDocument
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

// ListByBilling returns all instances for a given billing ID.
func (r *InstanceRepository) ListByBilling(billingID string) ([]*instance.Instance, error) {
	cursor, err := r.collection.Find(context.Background(), bson.M{"billing_id": billingID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var results []*instance.Instance
	for cursor.Next(context.Background()) {
		var doc instanceDocument
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

// mapDocumentToDomain converts a MongoDB document into a domain Instance.
func mapDocumentToDomain(doc instanceDocument) *instance.Instance {
	return &instance.Instance{
		ID:                    doc.ID,
		ListingID:             doc.ListingID,
		BillingID:             doc.BillingID,
		State:                 doc.State,
		HardwareSpecification: doc.HardwareSpecification,
	}
}
