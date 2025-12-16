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

// NewUserRepository creates a new MongoDB-backed user repository.
func NewUserRepository(collection *mongo.Collection) *UserRepository {
	return &UserRepository{collection: collection}
}

// Save inserts or updates a user.
func (r *UserRepository) Save(userEntity *user.User) error {
	doc := userDocument{
		ID:           userEntity.ToPrivateView().ID,
		Username:     userEntity.ToPrivateView().Username,
		PasswordHash: userEntity.ToPrivateView().PasswordHash,
		Developer:    userEntity.ToPrivateView().IsDev,
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
	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&doc)
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
	err := r.collection.FindOne(context.Background(), bson.M{"username": username}).Decode(&doc)
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

// mapDocumentToDomain converts a MongoDB document to the domain User.
func mapDocumentToDomain(doc userDocument) *user.User {
	return user.NewUser(doc.ID, doc.Username, doc.PasswordHash, doc.Developer)
}
