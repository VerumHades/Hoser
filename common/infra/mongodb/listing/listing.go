package listingmongodb

import (
	"context"
	"errors"

	"common/pkg/hardware"
	"common/pkg/listing"
	"common/pkg/money"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// listingDocument represents the MongoDB persistence model for a Listing.
type listingDocument struct {
	ID                    string                          `bson:"_id"`
	AuthorID              string                          `bson:"author_id"`
	Title                 string                          `bson:"title"`
	Description           string                          `bson:"description"`
	AccessMode            listing.ListingAccessMode       `bson:"access_mode"`
	HardwareSpecification *hardware.HardwareSpecification `bson:"hardware_spec,omitempty"`
	PriceValue            float64                         `bson:"price_value"`
	PriceCurrency         string                          `bson:"price_currency"`
}

// ListingRepository implements listing.ListingRepository using MongoDB.
type ListingRepository struct {
	collection *mongo.Collection
}

// NewListingRepository creates a new MongoDB-backed listing repository.
func NewListingRepository(collection *mongo.Collection) *ListingRepository {
	return &ListingRepository{
		collection: collection,
	}
}

// Save inserts or updates a listing.
func (r *ListingRepository) Save(listingEntity *listing.Listing) error {
	doc := listingDocument{
		ID:                    listingEntity.ID,
		AuthorID:              listingEntity.AuthorID,
		Title:                 listingEntity.Title,
		Description:           listingEntity.Description,
		AccessMode:            listingEntity.AccessMode,
		HardwareSpecification: listingEntity.HardwareSpecification,
		PriceValue:            listingEntity.Price.Amount,
		PriceCurrency:         listingEntity.Price.CurrencyCode,
	}

	_, err := r.collection.UpdateByID(
		context.Background(),
		doc.ID,
		bson.M{"$set": doc},
		options.Update().SetUpsert(true),
	)
	return err
}

// ListAll returns all listings in the collection.
func (r *ListingRepository) ListAll() ([]*listing.Listing, error) {
	cursor, err := r.collection.Find(context.Background(), bson.M{}) // empty filter matches all documents
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var results []*listing.Listing
	for cursor.Next(context.Background()) {
		var doc listingDocument
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

// GetByID retrieves a listing by ID.
func (r *ListingRepository) GetByID(listingID string) (*listing.Listing, error) {
	var doc listingDocument
	err := r.collection.FindOne(context.Background(), bson.M{"_id": listingID}).Decode(&doc)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return mapDocumentToDomain(doc), nil
}

// Delete removes a listing by ID.
func (r *ListingRepository) Delete(listingID string) error {
	_, err := r.collection.DeleteOne(context.Background(), bson.M{"_id": listingID})
	return err
}

// ListByAuthor returns all listings authored by the given author ID.
func (r *ListingRepository) ListByAuthor(authorID string) ([]*listing.Listing, error) {
	cursor, err := r.collection.Find(context.Background(), bson.M{"author_id": authorID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var results []*listing.Listing
	for cursor.Next(context.Background()) {
		var doc listingDocument
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

// mapDocumentToDomain converts a MongoDB document to the domain Listing.
func mapDocumentToDomain(doc listingDocument) *listing.Listing {
	return &listing.Listing{
		ID:                    doc.ID,
		AuthorID:              doc.AuthorID,
		Title:                 doc.Title,
		Description:           doc.Description,
		AccessMode:            doc.AccessMode,
		HardwareSpecification: doc.HardwareSpecification,
		Price:                 money.Money{Amount: doc.PriceValue, CurrencyCode: doc.PriceCurrency},
	}
}
