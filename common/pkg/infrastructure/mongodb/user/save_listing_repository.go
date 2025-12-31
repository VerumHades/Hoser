package mongodbuser

import (
	"context"
	"errors"

	"common/pkg/domain/user"
	"common/pkg/shared"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoSavedListingRepository implements command and query repositories for user saved listings.
type MongoSavedListingRepository struct {
	collection *mongo.Collection
}

// NewMongoSavedListingRepository constructs a repository backed by a Mongo collection.
func NewMongoSavedListingRepository(collection *mongo.Collection) *MongoSavedListingRepository {
	if collection == nil {
		panic("mongo collection must not be nil")
	}
	return &MongoSavedListingRepository{collection: collection}
}

// EnsureIndexes creates required MongoDB indexes for saved listings.
func (repo *MongoSavedListingRepository) EnsureIndexes(ctx context.Context) error {
	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "userID", Value: 1}, {Key: "listingID", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "userID", Value: 1},
				{Key: "createdAt", Value: 1},
			},
		},
	}

	_, err := repo.collection.Indexes().CreateMany(ctx, indexModels)
	return err
}

// -------------------- Command Repository --------------------

func (repo *MongoSavedListingRepository) Create(
	ctx context.Context,
	transaction shared.Transaction,
	item *user.SavedListing,
) (*user.SavedListing, error) {
	if item == nil {
		return nil, errors.New("saved listing cannot be nil")
	}

	sessionCtx := transaction.SessionContext(ctx)
	_, err := repo.collection.InsertOne(sessionCtx, item)
	if mongo.IsDuplicateKeyError(err) {
		return nil, shared.ErrAlreadyExists
	}
	return item, err
}

func (repo *MongoSavedListingRepository) Delete(
	ctx context.Context,
	transaction shared.Transaction,
	itemID shared.SavedListingID,
) error {
	sessionCtx := transaction.SessionContext(ctx)
	result, err := repo.collection.DeleteOne(sessionCtx, bson.M{"_id": itemID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return shared.ErrNotFound
	}
	return nil
}

// -------------------- Query Repository --------------------

func (repo *MongoSavedListingRepository) GetByID(
	ctx context.Context,
	itemID shared.SavedListingID,
) (*user.SavedListing, error) {
	var item user.SavedListing
	err := repo.collection.FindOne(ctx, bson.M{"_id": itemID}).Decode(&item)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, shared.ErrNotFound
	}
	return &item, err
}

func (repo *MongoSavedListingRepository) ExistsByUserAndListing(
	ctx context.Context,
	userID shared.UserID,
	listingID shared.ListingID,
) (bool, error) {
	filter := bson.M{"userID": userID, "listingID": listingID}
	count, err := repo.collection.CountDocuments(ctx, filter, options.Count().SetLimit(1))
	return count == 1, err
}

func (repo *MongoSavedListingRepository) FetchNextBatchByUser(
	ctx context.Context,
	userID shared.UserID,
	request shared.BatchRequest[user.SavedListingCursor],
) ([]*user.SavedListing, user.SavedListingCursor, error) {
	filter := bson.M{"userID": userID}
	filter["createdAt"] = bson.M{"$gt": request.Cursor.LastCreatedAt}

	findOptions := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: 1}}).
		SetLimit(int64(request.MaxBatchSize))

	cursor, err := repo.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, user.SavedListingCursor{}, err
	}
	defer cursor.Close(ctx)

	var items []*user.SavedListing
	for cursor.Next(ctx) {
		var s user.SavedListing
		if err := cursor.Decode(&s); err != nil {
			return nil, user.SavedListingCursor{}, err
		}
		items = append(items, &s)
	}

	if len(items) == 0 {
		return items, user.SavedListingCursor{}, nil
	}

	last := items[len(items)-1]
	return items, user.SavedListingCursor{LastCreatedAt: last.CreatedAt()}, nil
}
