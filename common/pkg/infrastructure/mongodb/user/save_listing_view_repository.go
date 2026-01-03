package mongodbuser

import (
	"context"
	"time"

	"common/pkg/domain/entities/user"
	mongodbregistry "common/pkg/infrastructure/mongodb/registry"
	"common/pkg/shared"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoUserSavedListingViewRepository implements UserSavedListingViewRepository
// using MongoDB aggregation with a $lookup join.
type MongoUserSavedListingViewRepository struct {
	savedListingCollection *mongo.Collection
	listingCollection      *mongo.Collection
}

// NewMongoUserSavedListingViewRepository constructs the repository using the registry.
// This resolves the "multiple parameters of same type" error in Wire.
func NewMongoUserSavedListingViewRepository(reg *mongodbregistry.DatabaseRegistry) *MongoUserSavedListingViewRepository {
	if reg.Libraries == nil || reg.Listings == nil {
		panic("libraries or listings collection must not be nil")
	}
	return &MongoUserSavedListingViewRepository{
		savedListingCollection: reg.Libraries,
		listingCollection:      reg.Listings,
	}
}

// FetchNextBatchOfUserSavedListings fetches a batch of saved listings joined with public listing info.
func (repo *MongoUserSavedListingViewRepository) FetchNextBatchOfUserSavedListings(
	ctx context.Context,
	userID shared.UserID,
	request shared.BatchRequest[user.UserSavedListingViewCursor],
) ([]*user.UserSavedListingView, user.UserSavedListingViewCursor, error) {

	matchFilter := bson.M{
		"user_id": userID,
	}

	if !request.Cursor.LastCreatedAt.IsZero() {
		matchFilter["createdAt"] = bson.M{
			"$gt": request.Cursor.LastCreatedAt,
		}
	}

	pipeline := mongo.Pipeline{
		// Match saved listings for the user after the cursor
		{{Key: "$match", Value: matchFilter}},
		// Join with listings collection
		{{Key: "$lookup", Value: bson.M{
			"from":         repo.listingCollection.Name(),
			"localField":   "listing_id",
			"foreignField": "_id",
			"as":           "listing",
		}}},
		// Unwind the joined listing array
		{{Key: "$unwind", Value: "$listing"}},
		// Sort by createdAt ascending
		{{Key: "$sort", Value: bson.D{{Key: "createdAt", Value: 1}}}},
		// Limit to batch size
		{{Key: "$limit", Value: int64(request.MaxBatchSize)}},
		// Project into the DTO shape
		{{Key: "$project", Value: bson.M{
			"SavedListingID": "$_id",
			"ListingID":      "$listing._id",
			"Title":          "$listing.title",
			"Description":    "$listing.description",
			"CreatedAt":      "$createdAt",
		}}},
	}

	cursor, err := repo.savedListingCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, user.UserSavedListingViewCursor{}, err
	}
	defer cursor.Close(ctx)

	var items []*user.UserSavedListingView
	var lastCreatedAt time.Time

	for cursor.Next(ctx) {
		var view user.UserSavedListingView
		if err := cursor.Decode(&view); err != nil {
			return nil, user.UserSavedListingViewCursor{}, err
		}
		items = append(items, &view)
		if view.CreatedAt.After(lastCreatedAt) {
			lastCreatedAt = view.CreatedAt
		}
	}

	nextCursor := user.UserSavedListingViewCursor{LastCreatedAt: lastCreatedAt}
	return items, nextCursor, nil
}
