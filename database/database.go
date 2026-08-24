package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

func NewDatabase(sqlDB *sql.DB, mongoDB *mongo.Client, dbName string, timeout time.Duration) *Database {
	return &Database{
		SQLDB:   sqlDB,
		MongoDB: mongoDB,
		DBName:  dbName,
		Timeout: timeout,
	}
}

func (db *Database) Exec(ctx context.Context, query string, args ...any) (string, error) {
	troubleShootMsg := ""

	if db.SQLDB == nil {
		troubleShootMsg = "ensure the database connection is properly initialized before calling Exec. Check it in ./system/bootstrap/database"

		return troubleShootMsg, errors.New("database connection is nil")
	}

	result, err := db.SQLDB.ExecContext(ctx, query, args...)
	if err != nil {
		troubleShootMsg = "verify the SQL syntax, table/column names, and parameter types in your query."

		return troubleShootMsg, fmt.Errorf("failed to execute SQL: %w", err)
	}

	if result == nil {
		troubleShootMsg = "the SQL driver returned a nil result. Check database driver behavior in ./app/config/database or connection stability."

		return troubleShootMsg, errors.New("no result returned from SQL execution")
	}

	return troubleShootMsg, nil
}

func (db *Database) Query(ctx context.Context, query string, args ...any) (*sql.Rows, string, error) {
	troubleShootMsg := ""

	if db.SQLDB == nil {
		troubleShootMsg = "ensure the database connection is properly initialized before calling Exec. Check it in ./system/bootstrap/database"

		return nil, troubleShootMsg, errors.New("sql database connection is nil")
	}

	result, err := db.SQLDB.QueryContext(ctx, query, args...)
	if err != nil {
		troubleShootMsg = "verify your SQL query syntax, parameters, and database schema. Ensure the target table exists and is accessible."

		return nil, troubleShootMsg, fmt.Errorf("failed to execute SQL: %w", err)
	}

	if result == nil {
		troubleShootMsg = "the SQL driver returned a nil result. Check database driver behavior in ./app/config/database or connection stability."

		return nil, troubleShootMsg, errors.New("no result returned from SQL execution")
	}

	return result, troubleShootMsg, nil
}

func (db *Database) Collection(collectionName string) (*mongo.Collection, string, error) {
	troubleShootMsg := ""

	if db == nil || db.MongoDB == nil {
		troubleShootMsg = "ensure the MongoDB connection is properly initialized before accessing collections. Check it in ./system/bootstrap/database"

		return nil, troubleShootMsg, errors.New("mongodb connection is nil")
	}

	result := db.MongoDB.Database(db.DBName).Collection(collectionName)
	if result == nil {
		troubleShootMsg = "failed to retrieve the MongoDB collection. Verify the database name and collection name."

		return nil, troubleShootMsg, errors.New("mongo collection is nil")
	}

	return result, troubleShootMsg, nil
}

func (db *Database) Close() (string, error) {
	troubleShootMsg := ""

	if db == nil {
		troubleShootMsg = "the database instance is nil. Ensure the instance is allocated before attempting to close it. " +
			"Check it in ./system/bootstrap/database"

		return troubleShootMsg, errors.New("database instance is nil")
	}

	var errs []error

	if db.SQLDB != nil {
		if err := db.SQLDB.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close SQL connection: %w", err))
		}
	}

	if db.MongoDB != nil {
		ctx, cancel := context.WithTimeout(context.Background(), db.Timeout)
		defer cancel()

		if err := db.MongoDB.Disconnect(ctx); err != nil {
			errs = append(errs, fmt.Errorf("failed to disconnect MongoDB: %w", err))
		}
	}

	if len(errs) > 0 {
		troubleShootMsg = "Check network connectivity or timeout configurations during the database disconnection phase. " +
			"Check timeout configuration ini ./app/config/database."

		return troubleShootMsg, errs[0]
	}

	return troubleShootMsg, nil
}
