package mongodbuser

import (
	"context"
	"errors"
	"time"

	"common/pkg/domain/user"
	mongodbregistry "common/pkg/infrastructure/mongodb/registry"
	"common/pkg/shared"
	"common/pkg/util"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoSavedListingRepository struct {
	collection *mongo.Collection
}

func NewMongoSavedListingRepository(databaseRegistry *mongodbregistry.DatabaseRegistry) *MongoSavedListingRepository {
	if databaseRegistry.Libraries == nil {
		panic("libraries collection must not be nil")
	}
	return &MongoSavedListingRepository{collection: databaseRegistry.Libraries}
}

func (repo *MongoSavedListingRepository) EnsureIndexes(ctx context.Context) error {
	indexModels := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "listing_id", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "created_at", Value: 1},
			},
		},
	}

	_, err := repo.collection.Indexes().CreateMany(ctx, indexModels)
	return err
}

type savedListingDocument struct {
	ID        shared.SavedListingID `bson:"_id"`
	UserID    shared.UserID         `bson:"user_id"`
	ListingID shared.ListingID      `bson:"listing_id"`
	CreatedAt int64                 `bson:"created_at"`
}

func mapSavedListingEntityToDocument(entity *user.SavedListing) *savedListingDocument {
	return &savedListingDocument{
		ID:        entity.ID(),
		UserID:    entity.UserID(),
		ListingID: entity.ListingID(),
		CreatedAt: entity.CreatedAt().UnixNano(),
	}
}

func mapSavedListingDocumentToEntity(document *savedListingDocument) (*user.SavedListing, error) {
	return user.NewSavedListingWithID(
		document.ID,
		document.UserID,
		document.ListingID,
		time.Unix(0, document.CreatedAt),
	)
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

	operationContext := util.ResolveTransactionalContext(ctx, transaction)
	document := mapSavedListingEntityToDocument(item)

	_, err := repo.collection.InsertOne(operationContext, document)
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
	operationContext := util.ResolveTransactionalContext(ctx, transaction)

	result, err := repo.collection.DeleteOne(
		operationContext,
		bson.M{"_id": itemID},
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

func (repo *MongoSavedListingRepository) GetByID(
	ctx context.Context,
	itemID shared.SavedListingID,
) (*user.SavedListing, error) {
	var document savedListingDocument

	err := repo.collection.FindOne(
		ctx,
		bson.M{"_id": itemID},
	).Decode(&document)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, shared.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return mapSavedListingDocumentToEntity(&document)
}

func (repo *MongoSavedListingRepository) GetByUserAndListing(
	ctx context.Context,
	userID shared.UserID,
	listingID shared.ListingID,
) (*user.SavedListing, error) {
	var document savedListingDocument

	err := repo.collection.FindOne(
		ctx,
		bson.M{"listing_id": listingID, "user_id": userID},
	).Decode(&document)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, shared.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return mapSavedListingDocumentToEntity(&document)
}

func (repo *MongoSavedListingRepository) ExistsByUserAndListing(
	ctx context.Context,
	userID shared.UserID,
	listingID shared.ListingID,
) (bool, error) {
	filter := bson.M{
		"user_id":    userID,
		"listing_id": listingID,
	}

	count, err := repo.collection.CountDocuments(
		ctx,
		filter,
		options.Count().SetLimit(1),
	)
	return count == 1, err
}

func (repo *MongoSavedListingRepository) FetchNextBatchByUser(
	ctx context.Context,
	userID shared.UserID,
	request shared.BatchRequest[user.SavedListingCursor],
) ([]*user.SavedListing, user.SavedListingCursor, error) {
	filter := bson.M{
		"user_id": userID,
	}

	filter["created_at"] = bson.M{
		"$gt": request.Cursor.LastCreatedAt.UnixNano(),
	}

	findOptions := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: 1}}).
		SetLimit(int64(request.MaxBatchSize))

	cursor, err := repo.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, user.SavedListingCursor{}, err
	}
	defer cursor.Close(ctx)

	var items []*user.SavedListing
	for cursor.Next(ctx) {
		var document savedListingDocument
		if err := cursor.Decode(&document); err != nil {
			return nil, user.SavedListingCursor{}, err
		}

		entity, err := mapSavedListingDocumentToEntity(&document)
		if err != nil {
			return nil, user.SavedListingCursor{}, err
		}

		items = append(items, entity)
	}

	if len(items) == 0 {
		return items, user.SavedListingCursor{}, nil
	}

	last := items[len(items)-1]

	return items, user.SavedListingCursor{
		LastCreatedAt: last.CreatedAt(),
	}, nil
}
