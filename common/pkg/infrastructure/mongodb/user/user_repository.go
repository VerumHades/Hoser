package mongodbuser

import (
	"common/pkg/domain/user"
	"context"
	"errors"

	"common/pkg/shared"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoUserRepository implements UserCommandRepository and UserQueryRepository.
type MongoUserRepository struct {
	collection *mongo.Collection
}

// NewMongoUserRepository constructs a repository backed by a Mongo collection.
func NewMongoUserRepository(collection *mongo.Collection) *MongoUserRepository {
	if collection == nil {
		panic("mongo collection must not be nil")
	}
	return &MongoUserRepository{collection: collection}
}

// EnsureIndexes sets up required indexes for the users collection.
func (repo *MongoUserRepository) EnsureIndexes(ctx context.Context) error {
	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}

	_, err := repo.collection.Indexes().CreateMany(ctx, indexModels)
	return err
}

// -------------------- Command Repository --------------------

func (repo *MongoUserRepository) Create(
	ctx context.Context,
	transaction shared.Transaction,
	userEntity *user.User,
) (*user.User, error) {
	if userEntity == nil {
		return nil, errors.New("user cannot be nil")
	}

	sessionCtx := transaction.SessionContext(ctx)
	_, err := repo.collection.InsertOne(sessionCtx, userEntity)
	if mongo.IsDuplicateKeyError(err) {
		return nil, shared.ErrAlreadyExists
	}
	return userEntity, err
}

func (repo *MongoUserRepository) Update(
	ctx context.Context,
	transaction shared.Transaction,
	userEntity *user.User,
) (*user.User, error) {
	if userEntity == nil {
		return nil, errors.New("user cannot be nil")
	}

	sessionCtx := transaction.SessionContext(ctx)
	result, err := repo.collection.ReplaceOne(
		sessionCtx,
		bson.M{"_id": userEntity.ID()},
		userEntity,
	)
	if err != nil {
		return nil, err
	}
	if result.MatchedCount == 0 {
		return nil, shared.ErrNotFound
	}
	return userEntity, nil
}

func (repo *MongoUserRepository) Delete(
	ctx context.Context,
	transaction shared.Transaction,
	userID shared.UserID,
) error {
	sessionCtx := transaction.SessionContext(ctx)
	result, err := repo.collection.DeleteOne(sessionCtx, bson.M{"_id": userID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return shared.ErrNotFound
	}
	return nil
}

// -------------------- Query Repository --------------------

func (repo *MongoUserRepository) GetByID(
	ctx context.Context,
	userID shared.UserID,
) (*user.User, error) {
	var u user.User
	err := repo.collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&u)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, shared.ErrNotFound
	}
	return &u, err
}

func (repo *MongoUserRepository) Exists(
	ctx context.Context,
	userID shared.UserID,
) (bool, error) {
	count, err := repo.collection.CountDocuments(ctx, bson.M{"_id": userID}, options.Count().SetLimit(1))
	return count == 1, err
}

func (repo *MongoUserRepository) GetByUsername(
	ctx context.Context,
	username string,
) (*user.User, error) {
	var u user.User
	err := repo.collection.FindOne(ctx, bson.M{"username": username}).Decode(&u)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, shared.ErrNotFound
	}
	return &u, err
}
