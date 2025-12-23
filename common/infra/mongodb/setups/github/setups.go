package listing

import (
	githubsetups "common/pkg/setups/github"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoGitHubSetupRepository implements GitHubSetupRepository using MongoDB.
type MongoGitHubSetupRepository struct {
	collection *mongo.Collection
}

// NewMongoGitHubSetupRepository creates a new repository instance.
func NewMongoGitHubSetupRepository(collection *mongo.Collection) *MongoGitHubSetupRepository {
	return &MongoGitHubSetupRepository{collection: collection}
}

// Save inserts or updates a githubsetups.GitHubSetupDefinition.
func (r *MongoGitHubSetupRepository) Save(setup *githubsetups.GitHubSetupDefinition) error {
	if setup == nil {
		return errors.New("setup cannot be nil")
	}
	if setup.ID == "" {
		return errors.New("setup ID cannot be empty")
	}
	now := time.Now()
	if setup.CreatedAt.IsZero() {
		setup.CreatedAt = now
	}
	setup.UpdatedAt = now

	_, err := r.collection.UpdateByID(
		context.Background(),
		setup.ID,
		bson.M{"$set": setup},
		options.Update().SetUpsert(true),
	)
	return err
}

// GetByID retrieves a setup by its ID.
func (r *MongoGitHubSetupRepository) GetByID(id string) (*githubsetups.GitHubSetupDefinition, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	var doc githubsetups.GitHubSetupDefinition
	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

// GetByListingID retrieves all setups for a listing.
func (r *MongoGitHubSetupRepository) GetByListingID(listingID string) ([]*githubsetups.GitHubSetupDefinition, error) {
	if listingID == "" {
		return nil, errors.New("listingID cannot be empty")
	}
	cursor, err := r.collection.Find(context.Background(), bson.M{"listingid": listingID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var results []*githubsetups.GitHubSetupDefinition
	for cursor.Next(context.Background()) {
		var doc githubsetups.GitHubSetupDefinition
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		results = append(results, &doc)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// DeleteByID deletes a setup by its ID.
func (r *MongoGitHubSetupRepository) DeleteByID(id string) error {
	if id == "" {
		return errors.New("id cannot be empty")
	}
	_, err := r.collection.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}

// DeleteByListingID deletes all setups for a listing.
func (r *MongoGitHubSetupRepository) DeleteByListingID(listingID string) error {
	if listingID == "" {
		return errors.New("listingID cannot be empty")
	}
	_, err := r.collection.DeleteMany(context.Background(), bson.M{"listingid": listingID})
	return err
}
