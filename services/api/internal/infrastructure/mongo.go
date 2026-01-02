package infrastructure

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// connectMongo handles the initialization and health check of the MongoDB driver.
// It returns the client and a cleanup function (cancel) for the context.
func ConnectMongo(cfg Configuration) (*mongo.Client, context.Context, context.CancelFunc) {
	// 1. Setup Context with timeout for the connection phase
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// 2. Build Connection String
	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s",
		cfg.MongoDatabaseUserName,
		cfg.MongoDatabasePassword,
		cfg.MongoDatabaseHost,
		cfg.MongoDatabasePort,
		cfg.MongoDatabaseName,
	)

	// 3. Initialize Client
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		cancel()
		log.Fatalf("Failed to create MongoDB client: %v", err)
	}

	// 4. Verify Connection (Ping)
	if err := client.Ping(ctx, nil); err != nil {
		cancel()
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	fmt.Println("Successfully connected to MongoDB")
	return client, ctx, cancel
}
