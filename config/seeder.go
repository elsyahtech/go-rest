//nolint:dupl
package config

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

// MongoDBSeeders defines the interface that individual MongoDB seeder files must implement.
type MongoDBSeeder interface {
	Name() string
	Number() int
	Up(ctx context.Context, db *mongo.Database) error
}

type MongoDBSeederRegistry struct {
	seeders []MongoDBSeeder
}

func NewMongoDBSeederRegistry() *MongoDBSeederRegistry {
	return &MongoDBSeederRegistry{
		seeders: make([]MongoDBSeeder, 0),
	}
}

func (r *MongoDBSeederRegistry) MongoDBSeederRegister(m MongoDBSeeder) {
	r.seeders = append(r.seeders, m)
}

func (r *MongoDBSeederRegistry) GetMongoDBSeeders() []MongoDBSeeder {
	return r.seeders
}

// ======================================================================================
// MongoDBSeederRegistrar defines a function signature for registering database seeders.
// ======================================================================================.
type MongoDBSeederRegistrar func(*MongoDBSeederRegistry)
