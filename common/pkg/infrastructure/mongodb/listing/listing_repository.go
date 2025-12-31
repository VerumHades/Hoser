package mongodblisting

import (
	"common/pkg/domain/listing"
	"common/pkg/shared"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoListingRepository implements command and query repositories for listings.
type MongoListingRepository struct {
	collection *mongo.Collection
}

// NewMongoListingRepository constructs a repository backed by a Mongo collection.
func NewMongoListingRepository(collection *mongo.Collection) *MongoListingRepository {
	if collection == nil {
		panic("mongo collection must not be nil")
	}
	return &MongoListingRepository{collection: collection}
}

// EnsureIndexes creates required MongoDB indexes for listings.
func (repo *MongoListingRepository) EnsureIndexes(ctx context.Context) error {
	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "authorID", Value: 1},
				{Key: "createdAt", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "createdAt", Value: 1},
			},
		},
	}

	_, err := repo.collection.Indexes().CreateMany(ctx, indexModels)
	return err
}

// -------------------- Command Repository --------------------

func (repo *MongoListingRepository) Create(
	ctx context.Context,
	transaction shared.Transaction,
	listingEntity *listing.Listing,
) error {
	if listingEntity == nil {
		return errors.New("listing cannot be nil")
	}

	sessionCtx := transaction.SessionContext(ctx)
	_, err := repo.collection.InsertOne(sessionCtx, listingEntity)
	if mongo.IsDuplicateKeyError(err) {
		return shared.ErrAlreadyExists
	}
	return err
}

func (repo *MongoListingRepository) Update(
	ctx context.Context,
	transaction shared.Transaction,
	listingEntity *listing.Listing,
) error {
	if listingEntity == nil {
		return errors.New("listing cannot be nil")
	}

	sessionCtx := transaction.SessionContext(ctx)
	result, err := repo.collection.ReplaceOne(sessionCtx, bson.M{"_id": listingEntity.ID()}, listingEntity)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return shared.ErrNotFound
	}
	return nil
}

func (repo *MongoListingRepository) Delete(
	ctx context.Context,
	transaction shared.Transaction,
	listingID shared.ListingID,
) error {
	sessionCtx := transaction.SessionContext(ctx)
	result, err := repo.collection.DeleteOne(sessionCtx, bson.M{"_id": listingID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return shared.ErrNotFound
	}
	return nil
}

// -------------------- Query Repository --------------------

func (repo *MongoListingRepository) GetByID(
	ctx context.Context,
	listingID shared.ListingID,
) (*listing.Listing, error) {
	var l listing.Listing
	err := repo.collection.FindOne(ctx, bson.M{"_id": listingID}).Decode(&l)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, shared.ErrNotFound
	}
	return &l, err
}

func (repo *MongoListingRepository) GetByIDAndAuthor(
	ctx context.Context,
	listingID shared.ListingID,
	userID shared.UserID,
) (*listing.Listing, error) {
	var l listing.Listing
	filter := bson.M{"_id": listingID, "authorID": userID}
	err := repo.collection.FindOne(ctx, filter).Decode(&l)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, shared.ErrNotFound
	}
	return &l, err
}

func (repo *MongoListingRepository) Exists(
	ctx context.Context,
	listingID shared.ListingID,
) (bool, error) {
	count, err := repo.collection.CountDocuments(ctx, bson.M{"_id": listingID}, options.Count().SetLimit(1))
	return count == 1, err
}

func (repo *MongoListingRepository) FetchNextBatchByAuthor(
	ctx context.Context,
	authorID shared.UserID,
	request shared.BatchRequest[listing.ListingCursor],
) ([]*listing.Listing, listing.ListingCursor, error) {
	filter := bson.M{"authorID": authorID}
	return repo.fetchBatch(ctx, filter, request)
}

func (repo *MongoListingRepository) FetchNextBatchAll(
	ctx context.Context,
	request shared.BatchRequest[listing.ListingCursor],
) ([]*listing.Listing, listing.ListingCursor, error) {
	filter := bson.M{}
	return repo.fetchBatch(ctx, filter, request)
}

// -------------------- Internal --------------------

func (repo *MongoListingRepository) fetchBatch(
	ctx context.Context,
	filter bson.M,
	request shared.BatchRequest[listing.ListingCursor],
) ([]*listing.Listing, listing.ListingCursor, error) {
	filter["createdAt"] = bson.M{"$gt": request.Cursor.LastCreatedAt}

	findOptions := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: 1}}).
		SetLimit(int64(request.MaxBatchSize))

	cursor, err := repo.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, listing.ListingCursor{}, err
	}
	defer cursor.Close(ctx)

	var listings []*listing.Listing
	for cursor.Next(ctx) {
		var l listing.Listing
		if err := cursor.Decode(&l); err != nil {
			return nil, listing.ListingCursor{}, err
		}
		listings = append(listings, &l)
	}

	if len(listings) == 0 {
		return listings, listing.ListingCursor{}, nil
	}

	last := listings[len(listings)-1]
	return listings, listing.ListingCursor{
		LastCreatedAt: last.CreatedAt(),
	}, nil
}
