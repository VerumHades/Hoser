package billing

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BillingServiceInMongo struct {
	interserverSecret string
	jwtSecret         string
	client            *mongo.Client
	database          string
	accountsColl      *mongo.Collection
	paymentsColl      *mongo.Collection
}

func NewBillingServiceMongo(mongoURI, database, interserverSecret, jwtSecret string) (*BillingServiceInMongo, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	service := &BillingServiceInMongo{
		interserverSecret: interserverSecret,
		jwtSecret:         jwtSecret,
		client:            client,
		database:          database,
		accountsColl:      client.Database(database).Collection("accounts"),
		paymentsColl:      client.Database(database).Collection("payments"),
	}

	return service, nil
}

// CreateAccount creates a billing account for a given userID
func (s *BillingServiceInMongo) CreateAccount(userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := s.accountsColl.CountDocuments(ctx, bson.M{"userID": userID})
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("account already exists")
	}

	_, err = s.accountsColl.InsertOne(ctx, bson.M{
		"userID":    userID,
		"billingID": fmt.Sprintf("billing-%s", userID),
	})
	return err
}

// Bill adds a payment record for a given listingID and billingID
func (s *BillingServiceInMongo) Bill(listingID string, billingID string) error {
	if billingID == "" || listingID == "" {
		return errors.New("listingID and billingID required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.paymentsColl.UpdateOne(
		ctx,
		bson.M{"billingID": billingID},
		bson.M{"$push": bson.M{"listingIDs": listingID}},
		options.Update().SetUpsert(true),
	)
	return err
}

// ListPayments lists all listingIDs for a given userID
func (s *BillingServiceInMongo) ListPayments(userID string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var account struct {
		BillingID string `bson:"billingID"`
	}
	err := s.accountsColl.FindOne(ctx, bson.M{"userID": userID}).Decode(&account)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("account not found")
		}
		return nil, err
	}

	var payment struct {
		ListingIDs []string `bson:"listingIDs"`
	}
	err = s.paymentsColl.FindOne(ctx, bson.M{"billingID": account.BillingID}).Decode(&payment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return []string{}, nil
		}
		return nil, err
	}

	return payment.ListingIDs, nil
}
