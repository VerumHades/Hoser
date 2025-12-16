package paymentmongodb

import (
	"context"
	"errors"
	"time"

	"common/pkg/billing/payments/payment"
	"common/pkg/money"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// paymentDocument represents the MongoDB persistence model for Payment.
type paymentDocument struct {
	ID               string                `bson:"_id"`
	BillingAccountID string                `bson:"billing_account_id"`
	AmountValue      float64               `bson:"amount_value"` // store smallest unit
	AmountCurrency   string                `bson:"amount_currency"`
	Status           payment.PaymentStatus `bson:"status"`
	Kind             payment.PaymentKind   `bson:"kind"`
	CreatedAt        time.Time             `bson:"created_at"`
	PaidAt           *time.Time            `bson:"paid_at,omitempty"`
}

// oneTimeMetadataDocument stores one-time payment metadata.
type oneTimeMetadataDocument struct {
	PaymentID string `bson:"_id"`
	ListingID string `bson:"listing_id"`
	UserID    string `bson:"user_id"`
}

// subscriptionMetadataDocument stores subscription payment metadata.
type subscriptionMetadataDocument struct {
	PaymentID        string                          `bson:"_id"`
	InstanceID       string                          `bson:"instance_id"`
	CurrentPeriodEnd time.Time                       `bson:"current_period_end"`
	Type             payment.SubscriptionPaymentType `bson:"type"`
}

// PaymentRepository implements payment.PaymentRepository with MongoDB.
type PaymentRepository struct {
	collection                     *mongo.Collection
	oneTimeMetadataCollection      *mongo.Collection
	subscriptionMetadataCollection *mongo.Collection
}

// NewPaymentRepository creates a new MongoDB-backed payment repository.
func NewPaymentRepository(
	collection *mongo.Collection,
	oneTimeMetadataCollection *mongo.Collection,
	subscriptionMetadataCollection *mongo.Collection,
) *PaymentRepository {
	return &PaymentRepository{
		collection:                     collection,
		oneTimeMetadataCollection:      oneTimeMetadataCollection,
		subscriptionMetadataCollection: subscriptionMetadataCollection,
	}
}

// Save persists a payment.
func (r *PaymentRepository) Save(paymentEntity *payment.Payment) error {
	doc := paymentDocument{
		ID:               paymentEntity.ID,
		BillingAccountID: paymentEntity.BillingAccountID,
		AmountValue:      paymentEntity.Amount.Amount,
		AmountCurrency:   paymentEntity.Amount.CurrencyCode,
		Status:           paymentEntity.Status,
		Kind:             paymentEntity.Kind,
		CreatedAt:        paymentEntity.CreatedAt,
		PaidAt:           paymentEntity.PaidAt,
	}

	_, err := r.collection.UpdateByID(
		context.Background(),
		doc.ID,
		bson.M{"$set": doc},
		options.Update().SetUpsert(true),
	)
	return err
}

// GetByID retrieves a payment by ID.
func (r *PaymentRepository) GetByID(paymentID string) (*payment.Payment, error) {
	var doc paymentDocument
	err := r.collection.FindOne(context.Background(), bson.M{"_id": paymentID}).Decode(&doc)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return mapDocumentToDomain(doc), nil
}

// ListByBillingAccount returns all payments for a billing account.
func (r *PaymentRepository) ListByBillingAccount(accountID string) ([]*payment.Payment, error) {
	cursor, err := r.collection.Find(context.Background(), bson.M{"billing_account_id": accountID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var results []*payment.Payment
	for cursor.Next(context.Background()) {
		var doc paymentDocument
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

// AttachOneTimeMetadata stores one-time payment metadata.
func (r *PaymentRepository) AttachOneTimeMetadata(paymentID string, metadata payment.OneTimePaymentMetadata) error {
	doc := oneTimeMetadataDocument{
		PaymentID: paymentID,
		ListingID: metadata.ListingID,
		UserID:    metadata.UserID,
	}

	_, err := r.oneTimeMetadataCollection.UpdateByID(
		context.Background(),
		paymentID,
		bson.M{"$set": doc},
		options.Update().SetUpsert(true),
	)
	return err
}

// AttachSubscriptionMetadata stores subscription payment metadata.
func (r *PaymentRepository) AttachSubscriptionMetadata(paymentID string, metadata payment.SubscriptionPaymentMetadata) error {
	doc := subscriptionMetadataDocument{
		PaymentID:        paymentID,
		InstanceID:       metadata.InstanceID,
		CurrentPeriodEnd: metadata.CurrentPeriodEnd,
		Type:             metadata.Type,
	}

	_, err := r.subscriptionMetadataCollection.UpdateByID(
		context.Background(),
		paymentID,
		bson.M{"$set": doc},
		options.Update().SetUpsert(true),
	)
	return err
}

// GetOneTimeMetadata retrieves one-time payment metadata.
func (r *PaymentRepository) GetOneTimeMetadata(paymentID string) (payment.OneTimePaymentMetadata, error) {
	var doc oneTimeMetadataDocument
	err := r.oneTimeMetadataCollection.FindOne(context.Background(), bson.M{"_id": paymentID}).Decode(&doc)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return payment.OneTimePaymentMetadata{}, nil
	}
	if err != nil {
		return payment.OneTimePaymentMetadata{}, err
	}

	return payment.OneTimePaymentMetadata{
		ListingID: doc.ListingID,
		UserID:    doc.UserID,
	}, nil
}

// GetSubscriptionMetadata retrieves subscription payment metadata.
func (r *PaymentRepository) GetSubscriptionMetadata(paymentID string) (payment.SubscriptionPaymentMetadata, error) {
	var doc subscriptionMetadataDocument
	err := r.subscriptionMetadataCollection.FindOne(context.Background(), bson.M{"_id": paymentID}).Decode(&doc)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return payment.SubscriptionPaymentMetadata{}, nil
	}
	if err != nil {
		return payment.SubscriptionPaymentMetadata{}, err
	}

	return payment.SubscriptionPaymentMetadata{
		InstanceID:       doc.InstanceID,
		CurrentPeriodEnd: doc.CurrentPeriodEnd,
		Type:             doc.Type,
	}, nil
}

// mapDocumentToDomain converts a MongoDB document to the domain Payment.
func mapDocumentToDomain(doc paymentDocument) *payment.Payment {
	return &payment.Payment{
		ID:               doc.ID,
		BillingAccountID: doc.BillingAccountID,
		Amount:           money.Money{Amount: doc.AmountValue, CurrencyCode: doc.AmountCurrency},
		Status:           doc.Status,
		Kind:             doc.Kind,
		CreatedAt:        doc.CreatedAt,
		PaidAt:           doc.PaidAt,
	}
}
