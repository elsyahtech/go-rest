package config

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

// MongoDBMigration defines the interface that individual MongoDB migration files must implement.
type MongoDBMigration interface {
	Name() string
	Number() int
	Up(ctx context.Context, db *mongo.Database) error
}

// MigrationRegistry acts as the collector/container for migrations (similar to *fiber.App).
type MongoDBMigrationRegistry struct {
	migrations []MongoDBMigration
}

func NewMigrationRegistry() *MongoDBMigrationRegistry {
	return &MongoDBMigrationRegistry{
		migrations: make([]MongoDBMigration, 0),
	}
}

func (r *MongoDBMigrationRegistry) Register(m MongoDBMigration) {
	r.migrations = append(r.migrations, m)
}

func (r *MongoDBMigrationRegistry) GetMigrations() []MongoDBMigration {
	return r.migrations
}

// ======================================================================================
// MigrationRegistrar defines a function signature for registering database migrations.
// ======================================================================================.
type MongoDBMigrationRegistrar func(*MongoDBMigrationRegistry)
