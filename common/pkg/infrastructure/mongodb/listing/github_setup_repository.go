package mongodblisting

import (
	"context"
	"errors"

	"common/pkg/domain/entities/githubsetups"
	"common/pkg/domain/repositories"
	mongodbregistry "common/pkg/infrastructure/mongodb/registry"
	"common/pkg/shared"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoGitHubSetupRepository implements both command and query repositories.
type MongoGitHubSetupRepository struct {
	collection *mongo.Collection
}

// NewMongoGitHubSetupRepository constructs a repository using the registry.
func NewMongoGitHubSetupRepository(reg *mongodbregistry.DatabaseRegistry) *MongoGitHubSetupRepository {
	if reg.GithubSetups == nil {
		panic("github setups collection must not be nil")
	}
	return &MongoGitHubSetupRepository{collection: reg.GithubSetups}
}

// EnsureGitHubSetupIndexes ensures indexes for efficient queries and pagination.
func (repo *MongoGitHubSetupRepository) EnsureIndexes(ctx context.Context) error {
	indexModels := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "listingID", Value: 1},
				{Key: "createdAt", Value: 1},
			},
		},
	}

	_, err := repo.collection.Indexes().CreateMany(ctx, indexModels)
	return err
}

// -------------------- Command Repository --------------------

func (repo *MongoGitHubSetupRepository) Create(
	ctx context.Context,

	setup *githubsetups.GitHubSetupDefinition,
) (*githubsetups.GitHubSetupDefinition, error) {
	if setup == nil {
		return nil, errors.New("setup cannot be nil")
	}

	_, err := repo.collection.InsertOne(ctx, setup)
	if mongo.IsDuplicateKeyError(err) {
		return nil, shared.ErrAlreadyExists
	}
	return setup, err
}

func (repo *MongoGitHubSetupRepository) Update(
	ctx context.Context,

	setup *githubsetups.GitHubSetupDefinition,
) (*githubsetups.GitHubSetupDefinition, error) {
	if setup == nil {
		return nil, errors.New("setup cannot be nil")
	}

	result, err := repo.collection.ReplaceOne(
		ctx,
		bson.M{"_id": setup.ID()},
		setup,
	)
	if err != nil {
		return nil, err
	}
	if result.MatchedCount == 0 {
		return nil, shared.ErrNotFound
	}
	return setup, nil
}

func (repo *MongoGitHubSetupRepository) DeleteByID(
	ctx context.Context,

	setupID shared.GithubSetupID,
) error {
	result, err := repo.collection.DeleteOne(ctx, bson.M{"_id": setupID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return shared.ErrNotFound
	}
	return nil
}

func (repo *MongoGitHubSetupRepository) DeleteByListingID(
	ctx context.Context,

	listingID shared.ListingID,
) error {
	_, err := repo.collection.DeleteMany(ctx, bson.M{"listingID": listingID})
	return err
}

// -------------------- Query Repository --------------------

func (repo *MongoGitHubSetupRepository) GetByID(
	ctx context.Context,
	setupID shared.GithubSetupID,
) (*githubsetups.GitHubSetupDefinition, error) {
	var setup githubsetups.GitHubSetupDefinition
	err := repo.collection.FindOne(ctx, bson.M{"_id": setupID}).Decode(&setup)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, shared.ErrNotFound
	}
	return &setup, err
}

func (repo *MongoGitHubSetupRepository) FetchNextBatchByListing(
	ctx context.Context,
	listingID shared.ListingID,
	request shared.BatchRequest[repositories.GitHubSetupCursor],
) ([]*githubsetups.GitHubSetupDefinition, repositories.GitHubSetupCursor, error) {
	filter := bson.M{"listingID": listingID}
	filter["createdAt"] = bson.M{"$gt": request.Cursor.LastCreatedAt}

	findOptions := options.Find().
		SetSort(bson.D{
			{Key: "createdAt", Value: 1},
		}).
		SetLimit(int64(request.MaxBatchSize))

	cursor, err := repo.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, repositories.GitHubSetupCursor{}, err
	}
	defer cursor.Close(ctx)

	var setups []*githubsetups.GitHubSetupDefinition
	for cursor.Next(ctx) {
		var s githubsetups.GitHubSetupDefinition
		if err := cursor.Decode(&s); err != nil {
			return nil, repositories.GitHubSetupCursor{}, err
		}
		setups = append(setups, &s)
	}

	if len(setups) == 0 {
		return setups, repositories.GitHubSetupCursor{}, nil
	}

	lastSetup := setups[len(setups)-1]
	return setups, repositories.GitHubSetupCursor{
		LastCreatedAt: lastSetup.CreatedAt(),
	}, nil
}
