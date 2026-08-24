package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/elsyahtech/go-rest/config"

	//nolint:revive,nolintlint
	_ "github.com/lib/pq"
)

func POSTGRESQLConnection(cfg *config.Config) (*sql.DB, error) {
	if cfg == nil {
		return nil, errors.New("config is nil")
	}

	// Construct the POSTGRESQL Server connection string (DSN) using configuration parameters
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.TLS,
	)

	const timeOutContext = 10 * time.Second

	// Open a new database handle using the postgresql driver and DSN
	database, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("open postgresql connection: %w", err)
	}

	// Create a context with a timeout for the connection operation
	ctx, cancel := context.WithTimeout(context.Background(), timeOutContext)

	// Ensure the connection context cancellation is deferred
	defer cancel()

	// Ping the POSTGRESQL database to verify connection health and responsiveness
	if err := database.PingContext(ctx); err != nil {
		if errDB := database.Close(); errDB != nil {
			_ = errDB
		}

		return nil, fmt.Errorf("ping postgresql database: %w", err)
	}

	// Return the successfully established database connection pool handle
	return database, nil
}

// ==========================================================================
// ConnectDatabasePostgreSQL establishes a connection to the POSTGRESQL Server
// database using global configuration settings
// ==========================================================================.
func ConnectDatabasePostgreSQL() (*sql.DB, error) {
	cfg := config.GlobalConfig

	// Construct the POSTGRESQL Server connection string (DSN) using configuration parameters
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.TLS,
	)

	const timeOutContext = 10 * time.Second

	// Open a new database handle using the postgresql driver and DSN
	database, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("open postgresql connection: %w", err)
	}

	// Create a context with a timeout for the connection operation
	ctx, cancel := context.WithTimeout(context.Background(), timeOutContext)

	// Ensure the connection context cancellation is deferred
	defer cancel()

	// Ping the POSTGRESQL database to verify connection health and responsiveness
	if err := database.PingContext(ctx); err != nil {
		if errDB := database.Close(); errDB != nil {
			_ = errDB
		}

		return nil, fmt.Errorf("ping postgresql database: %w", err)
	}

	// Return the successfully established database connection pool handle
	return database, nil
}
