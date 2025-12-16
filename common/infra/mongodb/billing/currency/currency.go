package currencymongodb

import (
	"context"
	"errors"

	"common/pkg/billing/currency"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// currencyDocument is the MongoDB persistence model.
type currencyDocument struct {
	Code   string `bson:"_id"` // use Code as the primary key
	Name   string `bson:"name"`
	Symbol string `bson:"symbol,omitempty"`
}

// CurrencyRepository is a MongoDB-backed implementation of currency.CurrencyRepository.
type CurrencyRepository struct {
	collection *mongo.Collection
}

// NewCurrencyRepository creates a new MongoDB-backed currency repository.
func NewCurrencyRepository(collection *mongo.Collection) *CurrencyRepository {
	return &CurrencyRepository{
		collection: collection,
	}
}

// Save inserts or updates a currency in the database.
func (r *CurrencyRepository) Save(currencyEntity *currency.Currency) error {
	doc := currencyDocument{
		Code:   currencyEntity.Code,
		Name:   currencyEntity.Name,
		Symbol: currencyEntity.Symbol,
	}

	_, err := r.collection.UpdateByID(
		context.Background(),
		doc.Code,
		bson.M{"$set": doc},
		options.Update().SetUpsert(true),
	)

	return err
}

// GetByCode retrieves a currency by its code.
func (r *CurrencyRepository) GetByCode(code string) (*currency.Currency, error) {
	var doc currencyDocument

	err := r.collection.FindOne(
		context.Background(),
		bson.M{"_id": code},
	).Decode(&doc)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return mapDocumentToDomain(doc), nil
}

// ListAll returns all currencies in the database.
func (r *CurrencyRepository) ListAll() ([]*currency.Currency, error) {
	cursor, err := r.collection.Find(
		context.Background(),
		bson.M{},
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var results []*currency.Currency
	for cursor.Next(context.Background()) {
		var doc currencyDocument
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

// mapDocumentToDomain converts the MongoDB document to the domain model.
func mapDocumentToDomain(doc currencyDocument) *currency.Currency {
	return &currency.Currency{
		Code:   doc.Code,
		Name:   doc.Name,
		Symbol: doc.Symbol,
	}
}
