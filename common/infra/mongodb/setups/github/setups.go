package githubsetups

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"common/pkg/util"
)

// =================== DOMAIN ===================

// GitHubSetupDefinition represents a setup linked to a listing.
type GitHubSetupDefinition struct {
	ID          string    `bson:"_id"`          // unique ID
	ListingID   string    `bson:"listing_id"`   // associated listing
	RepoURL     string    `bson:"repo_url"`     // GitHub repository URL
	AccessToken string    `bson:"access_token"` // access token for private repos
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
}

// =================== REPOSITORY ===================

type GitHubSetupRepository interface {
	Save(definition *GitHubSetupDefinition) error
	GetByListingID(listingID string) (*GitHubSetupDefinition, error)
	Delete(listingID string) error
	ListAll() ([]*GitHubSetupDefinition, error)
}

// =================== SERVICE ===================

type GitHubSetupService struct {
	repo GitHubSetupRepository
}

// NewService creates a new GitHubSetupService instance.
func NewService(repo GitHubSetupRepository) *GitHubSetupService {
	return &GitHubSetupService{repo: repo}
}

// GuardedSave validates and saves a setup definition.
func (s *GitHubSetupService) GuardedSave(def *GitHubSetupDefinition) error {
	if def.ListingID == "" {
		return errors.New("listing ID cannot be empty")
	}
	if def.RepoURL == "" {
		return errors.New("repository URL cannot be empty")
	}
	now := time.Now()
	if def.ID == "" {
		def.ID = util.GenerateUUID()
		def.CreatedAt = now
	}
	def.UpdatedAt = now

	return s.repo.Save(def)
}

// GetByListing retrieves a setup definition for a listing.
func (s *GitHubSetupService) GetByListing(listingID string) (*GitHubSetupDefinition, error) {
	if listingID == "" {
		return nil, errors.New("listing ID cannot be empty")
	}
	return s.repo.GetByListingID(listingID)
}

// RemoveByListing deletes a setup definition for a listing.
func (s *GitHubSetupService) RemoveByListing(listingID string) error {
	if listingID == "" {
		return errors.New("listing ID cannot be empty")
	}
	return s.repo.Delete(listingID)
}

// =================== MONGO IMPLEMENTATION ===================

type MongoRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(collection *mongo.Collection) *MongoRepository {
	return &MongoRepository{collection: collection}
}

func (r *MongoRepository) Save(def *GitHubSetupDefinition) error {
	_, err := r.collection.UpdateByID(
		context.Background(),
		def.ID,
		bson.M{"$set": def},
		options.Update().SetUpsert(true),
	)
	return err
}

func (r *MongoRepository) GetByListingID(listingID string) (*GitHubSetupDefinition, error) {
	var doc GitHubSetupDefinition
	err := r.collection.FindOne(context.Background(), bson.M{"listing_id": listingID}).Decode(&doc)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

func (r *MongoRepository) Delete(listingID string) error {
	_, err := r.collection.DeleteOne(context.Background(), bson.M{"listing_id": listingID})
	return err
}

func (r *MongoRepository) ListAll() ([]*GitHubSetupDefinition, error) {
	cursor, err := r.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var results []*GitHubSetupDefinition
	for cursor.Next(context.Background()) {
		var doc GitHubSetupDefinition
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
