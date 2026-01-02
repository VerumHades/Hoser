package mongodblisting

import (
	"context"
	"errors"
	"time"

	"common/pkg/domain/listing"
	mongodbregistry "common/pkg/infrastructure/mongodb/registry"
	"common/pkg/shared"
	"common/pkg/util"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoListingRepository struct {
	collection *mongo.Collection
}

func NewMongoListingRepository(reg *mongodbregistry.DatabaseRegistry) *MongoListingRepository {
	if reg.Listings == nil {
		panic("listings collection must not be nil")
	}
	return &MongoListingRepository{collection: reg.Listings}
}

// EnsureIndexes creates MongoDB indexes for listings.
func (repo *MongoListingRepository) EnsureIndexes(ctx context.Context) error {
	indexModels := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "author_id", Value: 1},
				{Key: "created_at", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "created_at", Value: 1},
			},
		},
	}
	_, err := repo.collection.Indexes().CreateMany(ctx, indexModels)
	return err
}

// -------------------- Document Mapping --------------------

type listingDocument struct {
	ID                    shared.ListingID              `bson:"_id"`
	AuthorID              shared.UserID                 `bson:"author_id"`
	Title                 string                        `bson:"title"`
	Description           string                        `bson:"description"`
	AccessMode            listing.ListingAccessMode     `bson:"access_mode"`
	HardwareSpecification *shared.HardwareSpecification `bson:"hardware_spec"`
	PriceInMinorUnits     int64                         `bson:"price_in_minor_units"`
	CreatedAt             int64                         `bson:"created_at"` // nanoseconds
}

func mapEntityToDocument(entity *listing.Listing) *listingDocument {
	return &listingDocument{
		ID:                    entity.ID(),
		AuthorID:              entity.AuthorID(),
		Title:                 entity.Title(),
		Description:           entity.Description(),
		AccessMode:            entity.AccessMode(),
		HardwareSpecification: entity.HardwareSpecification(),
		PriceInMinorUnits:     entity.PriceInMinorUnits(),
		CreatedAt:             entity.CreatedAt().UnixNano(),
	}
}

func mapDocumentToEntity(doc *listingDocument) (*listing.Listing, error) {
	return listing.NewListingWithID(
		doc.ID,
		doc.AuthorID,
		doc.Title,
		doc.Description,
		doc.AccessMode,
		doc.HardwareSpecification,
		doc.PriceInMinorUnits,
		time.Unix(0, doc.CreatedAt),
	)
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

	operationCtx := util.ResolveTransactionalContext(ctx, transaction)
	document := mapEntityToDocument(listingEntity)

	_, err := repo.collection.InsertOne(operationCtx, document)
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

	operationCtx := util.ResolveTransactionalContext(ctx, transaction)
	document := mapEntityToDocument(listingEntity)

	result, err := repo.collection.ReplaceOne(
		operationCtx,
		bson.M{"_id": listingEntity.ID()},
		document,
	)
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
	operationCtx := util.ResolveTransactionalContext(ctx, transaction)
	result, err := repo.collection.DeleteOne(operationCtx, bson.M{"_id": listingID})
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
	var doc listingDocument
	err := repo.collection.FindOne(ctx, bson.M{"_id": listingID}).Decode(&doc)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, shared.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return mapDocumentToEntity(&doc)
}

func (repo *MongoListingRepository) GetByIDAndAuthor(
	ctx context.Context,
	listingID shared.ListingID,
	userID shared.UserID,
) (*listing.Listing, error) {
	var doc listingDocument
	filter := bson.M{"_id": listingID, "author_id": userID}
	err := repo.collection.FindOne(ctx, filter).Decode(&doc)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, shared.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return mapDocumentToEntity(&doc)
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
	filter := bson.M{"author_id": authorID}
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
	filter["created_at"] = bson.M{"$gt": request.Cursor.LastCreatedAt.UnixNano()}

	findOptions := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: 1}}).
		SetLimit(int64(request.MaxBatchSize))

	cursor, err := repo.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, listing.ListingCursor{}, err
	}
	defer cursor.Close(ctx)

	var listings []*listing.Listing
	for cursor.Next(ctx) {
		var doc listingDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, listing.ListingCursor{}, err
		}
		entity, err := mapDocumentToEntity(&doc)
		if err != nil {
			return nil, listing.ListingCursor{}, err
		}
		listings = append(listings, entity)
	}

	if len(listings) == 0 {
		return listings, listing.ListingCursor{}, nil
	}

	last := listings[len(listings)-1]
	return listings, listing.ListingCursor{
		LastCreatedAt: last.CreatedAt(),
	}, nil
}
