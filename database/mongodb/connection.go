package mongodb

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/elsyahtech/go-rest/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// ===============================================================================================================
// ConnectDatabaseMongoDB establishes a connection to the MongoDB database using global configuration settings
// ===============================================================================================================.
func MONGODBConnection(cfg *config.Config) (*mongo.Client, error) {
	if cfg == nil {
		return nil, errors.New("config is nil")
	}

	dsn := getDsn(cfg)

	const timeOut = 10 * time.Second

	// Create a context with a timeout for the connection operation
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)

	// Ensure the connection context cancellation is deferred
	defer cancel()

	// Configure MongoDB client options with URI and server selection timeout
	clientOpts := options.Client().ApplyURI(dsn).SetServerSelectionTimeout(timeOut)

	// Connect to the MongoDB server instance
	client, err := mongo.Connect(clientOpts)
	if err != nil {
		return nil, fmt.Errorf("open mongodb connection: %w", err)
	}

	// Ping the MongoDB server to verify connection health and responsiveness
	if err := client.Ping(ctx, nil); err != nil {
		if errDB := client.Disconnect(ctx); errDB != nil {
			_ = errDB
		}

		return nil, fmt.Errorf("ping mongodb database: %w", err)
	}

	// Return a initialized MongoDBConnection wrapper instance
	return client, nil
}

// ====================================================================================================
// getDSNMongoDB constructs the MongoDB connection string (DSN) based on configuration parameters.
// ====================================================================================================.
func getDsn(cfg *config.Config) string {
	if cfg.Database.Username != "" && cfg.Database.Password != "" {
		// Escape the database username to handle special characters safely
		username := url.QueryEscape(cfg.Database.Username)

		// Escape the database password to handle special characters safely
		password := url.QueryEscape(cfg.Database.Password)

		// Construct and return authenticated MongoDB DSN with admin authSource
		return fmt.Sprintf(
			"mongodb://%s:%s@%s:%s/%s?authSource=admin",
			username,
			password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.Name,
		)
	}

	// Construct and return unauthenticated MongoDB DSN
	return fmt.Sprintf(
		"mongodb://%s:%s/%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)
}
