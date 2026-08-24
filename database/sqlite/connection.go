package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/elsyahtech/go-rest/config"

	//nolint:revive,nolintlint
	_ "modernc.org/sqlite"
)

// ==========================================================================
// ConnectDatabaseSQLite establishes a connection to the SQLite
// database using global configuration settings
// ==========================================================================.
func SQLITEConnection() (*sql.DB, error) {
	cfg := config.GlobalConfig
	dsn := cfg.Database.Name

	const timeOutContext = 10 * time.Second

	// Open a new database handle using the sqlite driver and DSN
	database, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite connection: %w", err)
	}

	// Create a context with a timeout for the connection operation
	ctx, cancel := context.WithTimeout(context.Background(), timeOutContext)

	// Ensure the connection context cancellation is deferred
	defer cancel()

	// Ping the SQLITE database to verify connection health and responsiveness
	if err := database.PingContext(ctx); err != nil {
		if errDB := database.Close(); errDB != nil {
			_ = errDB
		}

		return nil, fmt.Errorf("ping sqlite database: %w", err)
	}

	// Return the successfully established database connection pool handle
	return database, nil
}

// ==========================================================================
// ConnectDatabaseSQLite establishes a connection to the SQLite
// database using global configuration settings
// ==========================================================================.
func ConnectDatabaseSQLite(cfg *config.Config) (*sql.DB, error) {
	if cfg == nil {
		return nil, errors.New("config is nil")
	}

	dsn := cfg.Database.Name

	const timeOutContext = 10 * time.Second

	// Open a new database handle using the sqlite driver and DSN
	database, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite connection: %w", err)
	}

	// Create a context with a timeout for the connection operation
	ctx, cancel := context.WithTimeout(context.Background(), timeOutContext)

	// Ensure the connection context cancellation is deferred
	defer cancel()

	// Ping the SQLITE database to verify connection health and responsiveness
	if err := database.PingContext(ctx); err != nil {
		if errDB := database.Close(); errDB != nil {
			_ = errDB
		}

		return nil, fmt.Errorf("ping sqlite database: %w", err)
	}

	// Return the successfully established database connection pool handle
	return database, nil
}
