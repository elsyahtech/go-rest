package mssql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/elsyahtech/go-rest/config"

	//nolint:revive,nolintlint
	_ "github.com/microsoft/go-mssqldb"
)

// ==========================================================================
// MSSQLConnection establishes a connection to the Microsoft SQL Server
// database using global configuration settings
// ==========================================================================.
func MSSQLConnection(cfg *config.Config) (*sql.DB, error) {
	if cfg == nil {
		return nil, errors.New("config is nil")
	}

	// Construct the SQL Server connection string (DSN) using configuration parameters
	dsn := fmt.Sprintf(
		"server=%s;port=%s;database=%s;user id=%s;password=%s;encrypt=%s;trustservercertificate=%t",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.TLS,
		cfg.Database.TrustServerCertificate,
	)

	const timeOutContext = 10 * time.Second

	// Open a new database handle using the sqlserver driver and DSN
	database, err := sql.Open("sqlserver", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlserver connection: %w", err)
	}

	// Create a context with a timeout for the connection operation
	ctx, cancel := context.WithTimeout(context.Background(), timeOutContext)

	// Ensure the connection context cancellation is deferred
	defer cancel()

	// Ping the SQL Server database to verify connection health and responsiveness
	if err := database.PingContext(ctx); err != nil {
		if errDB := database.Close(); errDB != nil {
			_ = errDB
		}

		return nil, fmt.Errorf("ping sqlserver database: %w", err)
	}

	// Return the successfully established database connection pool handle
	return database, nil
}

// ==========================================================================
// ConnectDatabaseMsSQL establishes a connection to the Microsoft SQL Server
// database using global configuration settings
// ==========================================================================.
func ConnectDatabaseMsSQL() (*sql.DB, error) {
	cfg := config.GlobalConfig

	// Construct the SQL Server connection string (DSN) using configuration parameters
	dsn := fmt.Sprintf(
		"server=%s;port=%s;database=%s;user id=%s;password=%s;encrypt=%s;trustservercertificate=%t",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.TLS,
		cfg.Database.TrustServerCertificate,
	)

	const timeOutContext = 10 * time.Second

	// Open a new database handle using the sqlserver driver and DSN
	database, err := sql.Open("sqlserver", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlserver connection: %w", err)
	}

	// Create a context with a timeout for the connection operation
	ctx, cancel := context.WithTimeout(context.Background(), timeOutContext)

	// Ensure the connection context cancellation is deferred
	defer cancel()

	// Ping the SQL Server database to verify connection health and responsiveness
	if err := database.PingContext(ctx); err != nil {
		if errDB := database.Close(); errDB != nil {
			_ = errDB
		}

		return nil, fmt.Errorf("ping sqlserver database: %w", err)
	}

	// Return the successfully established database connection pool handle
	return database, nil
}
