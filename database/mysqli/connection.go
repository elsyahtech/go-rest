package mysqli

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/elsyahtech/go-rest/config"

	//nolint:revive,nolintlint
	_ "github.com/go-sql-driver/mysql"
)

func MYSQLConnection(cfg *config.Config) (*sql.DB, error) {
	if cfg == nil {
		return nil, errors.New("config is nil")
	}

	// Construct the MYSQL Server connection string (DSN) using configuration parameters
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=%s",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		url.QueryEscape(cfg.App.Timezone),
	)

	const timeOutContext = 10 * time.Second

	// Open a new database handle using the mysql driver and DSN
	database, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("open mysql connection: %w", err)
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

		return nil, fmt.Errorf("ping mysql database: %w", err)
	}

	// Return the successfully established database connection pool handle
	return database, nil
}

// ==========================================================================
// ConnectDatabaseMySQL establishes a connection to the MYSQL Server
// database using global configuration settings
// ==========================================================================.
func ConnectDatabaseMySQL() (*sql.DB, error) {
	cfg := config.GlobalConfig

	// Construct the MYSQL Server connection string (DSN) using configuration parameters
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=%s",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		url.QueryEscape(cfg.App.Timezone),
	)

	const timeOutContext = 10 * time.Second

	// Open a new database handle using the mysql driver and DSN
	database, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("open mysql connection: %w", err)
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

		return nil, fmt.Errorf("ping mysql database: %w", err)
	}

	// Return the successfully established database connection pool handle
	return database, nil
}
