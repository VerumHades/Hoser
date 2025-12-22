package usermongodb

import (
	"context"
	"errors"

	"common/pkg/user"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// userDocument is the MongoDB persistence representation for a User.
type userDocument struct {
	ID           string `bson:"_id"`
	Username     string `bson:"username"`
	PasswordHash string `bson:"password_hash"`
	Developer    bool   `bson:"developer"`
}

// UserRepository implements user.UserRepository using MongoDB.
type UserRepository struct {
	collection *mongo.Collection
}

// NewUserRepository creates a new MongoDB-backed user repository
// and ensures required indexes exist.
func NewUserRepository(collection *mongo.Collection) (*UserRepository, error) {
	repository := &UserRepository{
		collection: collection,
	}

	if err := repository.ensureIndexes(); err != nil {
		return nil, err
	}

	return repository, nil
}

// Save inserts or updates a user.
// Username uniqueness is enforced by a database index.
func (r *UserRepository) Save(userEntity *user.User) error {
	privateView := userEntity.ToPrivateView()

	doc := userDocument{
		ID:           privateView.ID,
		Username:     privateView.Username,
		PasswordHash: privateView.PasswordHash,
		Developer:    privateView.IsDev,
	}

	_, err := r.collection.UpdateByID(
		context.Background(),
		doc.ID,
		bson.M{"$set": doc},
		options.Update().SetUpsert(true),
	)

	return err
}

// GetByID retrieves a user by ID.
func (r *UserRepository) GetByID(id string) (*user.User, error) {
	var doc userDocument

	err := r.collection.
		FindOne(context.Background(), bson.M{"_id": id}).
		Decode(&doc)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return mapDocumentToDomain(doc), nil
}

// GetByUsername retrieves a user by username.
func (r *UserRepository) GetByUsername(username string) (*user.User, error) {
	var doc userDocument

	err := r.collection.
		FindOne(context.Background(), bson.M{"username": username}).
		Decode(&doc)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return mapDocumentToDomain(doc), nil
}

// Delete removes a user by ID.
func (r *UserRepository) Delete(id string) error {
	_, err := r.collection.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}

// ensureIndexes creates required MongoDB indexes for correctness and performance.
func (r *UserRepository) ensureIndexes() error {
	usernameIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "username", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}

	_, err := r.collection.Indexes().CreateOne(context.Background(), usernameIndex)
	return err
}

// mapDocumentToDomain converts a MongoDB document to the domain User.
func mapDocumentToDomain(doc userDocument) *user.User {
	return user.NewUser(doc.ID, doc.Username, doc.PasswordHash, doc.Developer)
}
