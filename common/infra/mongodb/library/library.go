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

// NewSavedListingRepository creates a new MongoDB-backed saved listing repository
// and ensures required indexes exist.
func NewSavedListingRepository(collection *mongo.Collection) (*SavedListingRepository, error) {
	repository := &SavedListingRepository{
		collection: collection,
	}

	if err := repository.ensureIndexes(); err != nil {
		return nil, err
	}

	return repository, nil
}

// Save inserts a saved listing.
// Uniqueness is enforced by a compound database index.
func (r *SavedListingRepository) Save(item *library.SavedListing) error {
	doc := savedListingDocument{
		ID:        item.ID,
		UserID:    item.UserID,
		ListingID: item.ListingID,
	}

	_, err := r.collection.InsertOne(context.Background(), doc)
	return err
}

// ExistsByUserAndListing checks whether a user already saved a specific listing.
func (r *SavedListingRepository) ExistsByUserAndListing(userID string, listingID string) (bool, error) {
	filter := bson.M{
		"user_id":    userID,
		"listing_id": listingID,
	}

	err := r.collection.
		FindOne(
			context.Background(),
			filter,
			options.FindOne().SetProjection(bson.M{"_id": 1}),
		).
		Err()

	if errors.Is(err, mongo.ErrNoDocuments) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
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

	err := r.collection.
		FindOne(context.Background(), bson.M{"_id": itemID}).
		Decode(&doc)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return mapDocumentToDomain(doc), nil
}

// ensureIndexes creates required MongoDB indexes for correctness and performance.
func (r *SavedListingRepository) ensureIndexes() error {
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "user_id", Value: 1},
			{Key: "listing_id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}

	_, err := r.collection.Indexes().CreateOne(context.Background(), indexModel)
	return err
}

// mapDocumentToDomain converts a MongoDB document into a domain SavedListing.
func mapDocumentToDomain(doc savedListingDocument) *library.SavedListing {
	return &library.SavedListing{
		ID:        doc.ID,
		UserID:    doc.UserID,
		ListingID: doc.ListingID,
	}
}
