package database

import (
	"context"
	"database/sql"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

// ===============================================================================================================
// GlobalDB holds the global database connection instance for the application.
// ===============================================================================================================.
var GlobalDB *Database

type Database struct {
	SQLDB   *sql.DB
	MongoDB *mongo.Client
	DBName  string
	Timeout time.Duration
}

type DatabaseImpl interface {
	Exec(ctx context.Context, query string, args ...any) error
	Query(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	Collection(collectionName string) *mongo.Collection
	Close() error
}
