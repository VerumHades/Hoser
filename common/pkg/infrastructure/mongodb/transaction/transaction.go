package mongotransaction

import (
	"context"
	"errors"

	"common/pkg/shared"

	"go.mongodb.org/mongo-driver/mongo"
)

type MongoTransaction struct {
	session mongo.Session
	active  bool
}

type MongoTransactionProvider struct {
	mongoClient *mongo.Client
}

func NewMongoTransactionProvider(
	mongoClient *mongo.Client,
) *MongoTransactionProvider {
	if mongoClient == nil {
		panic("mongo client must not be nil")
	}

	return &MongoTransactionProvider{
		mongoClient: mongoClient,
	}
}

func (provider *MongoTransactionProvider) BeginTransaction(
	ctx context.Context,
) (shared.Transaction, error) {
	session, err := provider.mongoClient.StartSession()
	if err != nil {
		return nil, err
	}

	if err := session.StartTransaction(); err != nil {
		session.EndSession(ctx)
		return nil, err
	}

	return &MongoTransaction{
		session: session,
		active:  true,
	}, nil
}

func (transaction *MongoTransaction) SessionContext(
	ctx context.Context,
) context.Context {
	if !transaction.active {
		panic("attempted to use inactive mongo transaction")
	}

	return mongo.NewSessionContext(ctx, transaction.session)
}

func (transaction *MongoTransaction) Commit() error {
	if !transaction.active {
		return errors.New("transaction already completed")
	}

	err := transaction.session.CommitTransaction(context.Background())
	transaction.end()

	return err
}

func (transaction *MongoTransaction) Rollback() error {
	if !transaction.active {
		return errors.New("transaction already completed")
	}

	err := transaction.session.AbortTransaction(context.Background())
	transaction.end()

	return err
}

func (transaction *MongoTransaction) end() {
	transaction.active = false
	transaction.session.EndSession(context.Background())
}
