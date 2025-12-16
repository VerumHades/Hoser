package accountmongodb

import (
	"context"
	"errors"
	"time"

	"common/pkg/billing/account"
	"common/pkg/billing/payments"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// billingAccountDocument represents the MongoDB persistence model.
// It is intentionally separate from the domain model.
type billingAccountDocument struct {
	ID                string                       `bson:"_id"`
	OwnerID           string                       `bson:"owner_id"`
	Status            account.BillingAccountStatus `bson:"status"`
	PaymentProvider   payments.PaymentProvider     `bson:"payment_provider"`
	ProviderAccountID string                       `bson:"provider_account_id"`
	CreatedAt         time.Time                    `bson:"created_at"`
}

// BillingAccountRepository is a MongoDB-backed implementation of account.BillingAccountRepository.
type BillingAccountRepository struct {
	collection *mongo.Collection
}

// NewBillingAccountRepository creates a new MongoDB billing account repository.
func NewBillingAccountRepository(collection *mongo.Collection) *BillingAccountRepository {
	return &BillingAccountRepository{
		collection: collection,
	}
}

// Save persists a billing account using upsert semantics.
func (r *BillingAccountRepository) Save(accountEntity *account.BillingAccount) error {
	document := billingAccountDocument{
		ID:                accountEntity.ID,
		OwnerID:           accountEntity.OwnerID,
		Status:            accountEntity.Status,
		PaymentProvider:   accountEntity.PaymentProvider,
		ProviderAccountID: accountEntity.ProviderAccountID,
		CreatedAt:         accountEntity.CreatedAt,
	}

	_, err := r.collection.UpdateByID(
		context.Background(),
		document.ID,
		bson.M{"$set": document},
		options.Update().SetUpsert(true),
	)

	return err
}

// GetByID retrieves a billing account by its ID.
func (r *BillingAccountRepository) GetByID(accountID string) (*account.BillingAccount, error) {
	var document billingAccountDocument

	err := r.collection.FindOne(
		context.Background(),
		bson.M{"_id": accountID},
	).Decode(&document)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return mapDocumentToDomain(document), nil
}

// ListByOwner returns all billing accounts owned by the given owner ID.
func (r *BillingAccountRepository) ListByOwner(ownerID string) ([]*account.BillingAccount, error) {
	cursor, err := r.collection.Find(
		context.Background(),
		bson.M{"owner_id": ownerID},
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var results []*account.BillingAccount

	for cursor.Next(context.Background()) {
		var document billingAccountDocument
		if err := cursor.Decode(&document); err != nil {
			return nil, err
		}

		results = append(results, mapDocumentToDomain(document))
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func mapDocumentToDomain(document billingAccountDocument) *account.BillingAccount {
	return &account.BillingAccount{
		ID:                document.ID,
		OwnerID:           document.OwnerID,
		Status:            document.Status,
		PaymentProvider:   document.PaymentProvider,
		ProviderAccountID: document.ProviderAccountID,
		CreatedAt:         document.CreatedAt,
	}
}
