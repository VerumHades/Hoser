package librarymongodb

import (
	"context"
	"errors"

	"common/pkg/library"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// savedListingDocument represents the MongoDB persistence model for a saved listing.
type savedListingDocument struct {
	ID        string `bson:"_id"`
	UserID    string `bson:"user_id"`
	ListingID string `bson:"listing_id"`
}

// SavedListingRepository implements library.SavedListingRepository using MongoDB.
type SavedListingRepository struct {
	collection *mongo.Collection
}

// NewSavedListingRepository creates a new MongoDB-backed saved listing repository.
func NewSavedListingRepository(collection *mongo.Collection) *SavedListingRepository {
	return &SavedListingRepository{
		collection: collection,
	}
}

// Save inserts or updates a saved listing.
func (r *SavedListingRepository) Save(item *library.SavedListing) error {
	doc := savedListingDocument{
		ID:        item.ID,
		UserID:    item.UserID,
		ListingID: item.ListingID,
	}

	_, err := r.collection.UpdateByID(
		context.Background(),
		doc.ID,
		bson.M{"$set": doc},
		options.Update().SetUpsert(true),
	)
	return err
}

// GetByUserID retrieves all saved listings for a user.
func (r *SavedListingRepository) GetByUserID(userID string) ([]*library.SavedListing, error) {
	cursor, err := r.collection.Find(context.Background(), bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var results []*library.SavedListing
	for cursor.Next(context.Background()) {
		var doc savedListingDocument
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

// Delete removes a saved listing by its ID.
func (r *SavedListingRepository) Delete(itemID string) error {
	_, err := r.collection.DeleteOne(context.Background(), bson.M{"_id": itemID})
	return err
}

// GetByID retrieves a saved listing by its ID.
func (r *SavedListingRepository) GetByID(itemID string) (*library.SavedListing, error) {
	var doc savedListingDocument
	err := r.collection.FindOne(context.Background(), bson.M{"_id": itemID}).Decode(&doc)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return mapDocumentToDomain(doc), nil
}

// mapDocumentToDomain converts a MongoDB document into a domain SavedListing.
func mapDocumentToDomain(doc savedListingDocument) *library.SavedListing {
	return &library.SavedListing{
		ID:        doc.ID,
		UserID:    doc.UserID,
		ListingID: doc.ListingID,
	}
}
