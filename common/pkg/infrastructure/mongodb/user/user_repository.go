package mongodbuser

import (
	"context"
	"errors"

	"common/pkg/domain/user"
	mongodbregistry "common/pkg/infrastructure/mongodb/registry"
	"common/pkg/shared"
	"common/pkg/util"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoUserRepository struct {
	collection *mongo.Collection
}

func NewMongoUserRepository(databaseRegistry *mongodbregistry.DatabaseRegistry) *MongoUserRepository {
	if databaseRegistry.Users == nil {
		panic("users collection must not be nil")
	}
	return &MongoUserRepository{collection: databaseRegistry.Users}
}

func (repo *MongoUserRepository) EnsureIndexes(ctx context.Context) error {
	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}

	_, err := repo.collection.Indexes().CreateMany(ctx, indexModels)
	return err
}

type userDocument struct {
	ID           shared.UserID `bson:"_id"`
	Username     string        `bson:"username"`
	PasswordHash string        `bson:"password_hash"`
	Developer    bool          `bson:"developer"`
}

func mapEntityToDocument(userEntity *user.User) *userDocument {
	return &userDocument{
		ID:           userEntity.ID(),
		Username:     userEntity.Username(),
		PasswordHash: userEntity.PasswordHash(),
		Developer:    userEntity.IsDeveloper(),
	}
}

func mapDocumentToEntity(document *userDocument) (*user.User, error) {
	return user.NewUserWithID(
		document.ID,
		document.Username,
		document.PasswordHash,
		document.Developer,
	)
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

	operationContext := util.ResolveTransactionalContext(ctx, transaction)
	document := mapEntityToDocument(userEntity)

	_, err := repo.collection.InsertOne(operationContext, document)
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

	operationContext := util.ResolveTransactionalContext(ctx, transaction)
	document := mapEntityToDocument(userEntity)

	result, err := repo.collection.ReplaceOne(
		operationContext,
		bson.M{"_id": userEntity.ID()},
		document,
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
	operationContext := util.ResolveTransactionalContext(ctx, transaction)

	result, err := repo.collection.DeleteOne(
		operationContext,
		bson.M{"_id": userID},
	)
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
	var document userDocument

	err := repo.collection.FindOne(
		ctx,
		bson.M{"_id": userID},
	).Decode(&document)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, shared.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	userEntity, err := mapDocumentToEntity(&document)
	if err != nil {
		return nil, err
	}

	return userEntity, nil
}

func (repo *MongoUserRepository) Exists(
	ctx context.Context,
	userID shared.UserID,
) (bool, error) {
	count, err := repo.collection.CountDocuments(
		ctx,
		bson.M{"_id": userID},
		options.Count().SetLimit(1),
	)
	return count == 1, err
}

func (repo *MongoUserRepository) GetByUsername(
	ctx context.Context,
	username string,
) (*user.User, error) {
	var document userDocument

	err := repo.collection.FindOne(
		ctx,
		bson.M{"username": username},
	).Decode(&document)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, shared.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	userEntity, err := mapDocumentToEntity(&document)
	if err != nil {
		return nil, err
	}

	return userEntity, nil
}
