package mongodboutbox

import (
	"context"
	"errors"
	"time"

	"common/pkg/application/outbox"
	mongodbregistry "common/pkg/infrastructure/mongodb/registry"
	"common/pkg/shared"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoOutboxRepository struct {
	collection *mongo.Collection
}

func NewMongoOutboxRepository(reg *mongodbregistry.DatabaseRegistry) *MongoOutboxRepository {
	if reg.Outbox == nil {
		panic("outbox collection must not be nil")
	}
	return &MongoOutboxRepository{
		collection: reg.Outbox,
	}
}

func (repo *MongoOutboxRepository) EnsureIndexes(ctx context.Context) error {
	indexModels := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "occurred_at", Value: 1},
				{Key: "_id", Value: 1},
			},
		},
	}
	_, err := repo.collection.Indexes().CreateMany(ctx, indexModels)
	return err
}

// -------------------- Document Mapping --------------------

type outboxEventDocument struct {
	ID         shared.OutboxEventID `bson:"_id"`
	EventType  string               `bson:"event_type"`
	OccurredAt int64                `bson:"occurred_at"`
	Payload    []byte               `bson:"payload"`
}

func mapEntityToDocument(event *outbox.OutboxEvent) (*outboxEventDocument, error) {
	if event == nil {
		return nil, errors.New("outbox event cannot be nil")
	}

	return &outboxEventDocument{
		ID:         event.ID(),
		EventType:  event.EventType(),
		OccurredAt: event.OccurredAt().UnixNano(),
		Payload:    event.Payload(),
	}, nil
}

func mapDocumentToEntity(doc *outboxEventDocument) (*outbox.OutboxEvent, error) {
	return outbox.NewOutboxEventWithID(
		doc.ID,
		doc.EventType,
		doc.Payload,
		time.Unix(0, doc.OccurredAt),
	), nil
}

// -------------------- Command Repository --------------------

func (repo *MongoOutboxRepository) Create(
	ctx context.Context,
	event *outbox.OutboxEvent,
) error {
	document, err := mapEntityToDocument(event)
	if err != nil {
		return err
	}

	_, err = repo.collection.InsertOne(ctx, document)
	if mongo.IsDuplicateKeyError(err) {
		return shared.ErrAlreadyExists
	}
	return err
}

func (repo *MongoOutboxRepository) BulkCreate(
	ctx context.Context,
	events []*outbox.OutboxEvent,
) error {
	if len(events) == 0 {
		return nil
	}

	documents := make([]any, len(events))
	for i, event := range events {
		document, err := mapEntityToDocument(event)
		if err != nil {
			return err
		}
		documents[i] = document
	}

	_, err := repo.collection.InsertMany(ctx, documents)
	return err
}

func (repo *MongoOutboxRepository) Delete(
	ctx context.Context,
	eventID uuid.UUID,
) error {
	result, err := repo.collection.DeleteOne(ctx, bson.M{"_id": eventID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return shared.ErrNotFound
	}
	return nil
}

func (repo *MongoOutboxRepository) DeleteOlderThan(
	ctx context.Context,
	cutoff time.Time,
) error {
	_, err := repo.collection.DeleteMany(
		ctx,
		bson.M{
			"occurred_at": bson.M{
				"$lt": cutoff.UnixNano(),
			},
		},
	)
	return err
}

// -------------------- Query Repository --------------------

func (repo *MongoOutboxRepository) GetByID(
	ctx context.Context,
	eventID uuid.UUID,
) (*outbox.OutboxEvent, error) {
	var doc outboxEventDocument
	err := repo.collection.FindOne(ctx, bson.M{"_id": eventID}).Decode(&doc)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, shared.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return mapDocumentToEntity(&doc)
}

func (repo *MongoOutboxRepository) Exists(
	ctx context.Context,
	eventID uuid.UUID,
) (bool, error) {
	count, err := repo.collection.CountDocuments(
		ctx,
		bson.M{"_id": eventID},
		options.Count().SetLimit(1),
	)
	return count == 1, err
}

func (repo *MongoOutboxRepository) FetchNextBatchAfter(
	ctx context.Context,
	after time.Time,
	request shared.BatchRequest[outbox.OutboxEventCursor],
) ([]*outbox.OutboxEvent, outbox.OutboxEventCursor, error) {

	filter := bson.M{
		"$or": []bson.M{
			{
				"occurred_at": bson.M{
					"$gt": after.UnixNano(),
				},
			},
			{
				"occurred_at": after.UnixNano(),
				"_id": bson.M{
					"$gt": request.Cursor.LastID,
				},
			},
		},
	}

	findOptions := options.Find().
		SetSort(bson.D{
			{Key: "occurred_at", Value: 1},
			{Key: "_id", Value: 1},
		}).
		SetLimit(int64(request.MaxBatchSize))

	cursor, err := repo.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, outbox.OutboxEventCursor{}, err
	}
	defer cursor.Close(ctx)

	var events []*outbox.OutboxEvent
	for cursor.Next(ctx) {
		var doc outboxEventDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, outbox.OutboxEventCursor{}, err
		}
		event, err := mapDocumentToEntity(&doc)
		if err != nil {
			return nil, outbox.OutboxEventCursor{}, err
		}
		events = append(events, event)
	}

	if len(events) == 0 {
		return events, outbox.OutboxEventCursor{}, nil
	}

	last := events[len(events)-1]
	return events, outbox.OutboxEventCursor{
		LastOccurredAt: last.OccurredAt(),
		LastID:         last.ID(),
	}, nil
}
