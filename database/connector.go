package database

import (
	"errors"
	"fmt"
	"strings"

	"github.com/elsyahtech/go-rest/config"
	"github.com/elsyahtech/go-rest/database/mongodb"
	"github.com/elsyahtech/go-rest/database/mssql"
	"github.com/elsyahtech/go-rest/database/mysqli"
	"github.com/elsyahtech/go-rest/database/postgresql"
	"github.com/elsyahtech/go-rest/database/sqlite"
)

// ===============================================================================================================
// NewDatabase initializes and returns the appropriate database connection instance based on the configured driver.
// ===============================================================================================================.
func DBConnection() (*Database, string, error) {
	// Declare variables for config, error, database implementation, and troubleshooting action message
	cfg := config.GlobalConfig

	var (
		err             error
		database        *Database
		troubleShootMsg string
	)

	// Normalize the configured database driver string for consistent comparison
	driver := normalizeDatabaseDriver(cfg.Database.Driver)

	// Validate whether the specified database driver is supported
	if !supportedDatabaseDrivers[driver] {
		troubleShootMsg = "Verify database driver in ./app/config/database - check driver value is supported (mysql, postgresql, mssql, mongodb, sqlite)"

		return nil, troubleShootMsg, fmt.Errorf("unsupported database driver: %q", driver)
	}

	switch driver {
	// Handle MySQL and MariaDB connection initialization
	case MYSQL:
		troubleShootMsg = "Verify MySQL credentials and connectivity in ./app/config/database - Ensure the database server is UP and firewall allows connection"

		sqlDB, err := mysqli.MYSQLConnection(&cfg)
		if err != nil {
			return nil, troubleShootMsg, fmt.Errorf("failed to connect to MySQL: %w", err)
		}

		database = NewDatabase(sqlDB, nil, cfg.Database.Name, cfg.Database.Timeout)

	// Handle PostgreSQL connection initialization
	case POSTGRESQL:
		// PostgreSQL connection
		troubleShootMsg = "Verify PostgreSQL credentials and connectivity in ./app/config/database. " +
			"Ensure the database server is UP and firewall allows connection"

		sqlDB, err := postgresql.POSTGRESQLConnection(&cfg)
		if err != nil {
			return nil, troubleShootMsg, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
		}

		database = NewDatabase(sqlDB, nil, cfg.Database.Name, cfg.Database.Timeout)
	// Handle Microsoft SQL Server connection initialization
	case MSSQL:
		// Microsoft SQL Server connection
		troubleShootMsg = "Verify MSSQL credentials and connectivity in ./app/config/database. " +
			"Ensure the database server is UP and firewall allows connection"

		sqlDB, err := mssql.MSSQLConnection(&cfg)
		if err != nil {
			return nil, troubleShootMsg, fmt.Errorf("failed to connect to MSSQL: %w", err)
		}

		database = NewDatabase(sqlDB, nil, cfg.Database.Name, cfg.Database.Timeout)
	// Handle MongoDB connection initialization
	case MONGODB:
		// MongoDB connection (NoSQL, different implementation)
		troubleShootMsg = "Verify MongoDB credentials and connectivity in ./app/config/database. " +
			"Ensure the database server is UP and firewall allows connection"

		mongoConn, err := mongodb.MONGODBConnection(&cfg)
		if err != nil {
			return nil, troubleShootMsg, fmt.Errorf("failed to connect MongoDB: %w", err)
		}

		database = NewDatabase(nil, mongoConn, cfg.Database.Name, cfg.Database.Timeout)
	// Handle SQLite connection initialization
	case SQLITE:
		// SQLite connection
		troubleShootMsg = "Verify SQLite credentials and connectivity in ./app/config/database. " +
			"Ensure the database server is UP and firewall allows connection"

		sqlDB, err := sqlite.ConnectDatabaseSQLite(&cfg)
		if err != nil {
			return nil, troubleShootMsg, fmt.Errorf("failed to connect to SQLite: %w", err)
		}

		database = NewDatabase(sqlDB, nil, cfg.Database.Name, cfg.Database.Timeout)
	// Handle any unknown or unimplemented drivers
	default:
		troubleShootMsg = "implement database driver in ./app/config/database to support for: " + cfg.Database.Driver

		err = fmt.Errorf("database driver %q is not implemented", cfg.Database.Driver)
	}

	if err != nil {
		return nil, troubleShootMsg, fmt.Errorf("cannot connect: %w", err)
	}

	if database == nil {
		troubleShootMsg = "check database connection initialization logic. Ensure driver returns valid connection instance"

		return nil, troubleShootMsg, errors.New("database connection is nil")
	}

	// Return the successfully established database connection instance
	return database, "", nil
}

// ===============================================================================================================
// normalizeDatabaseDriver converts various driver aliases into standardized driver constant names.
// ===============================================================================================================.
func normalizeDatabaseDriver(driver string) string {
	normalized := strings.ToLower(strings.TrimSpace(driver))

	switch normalized {
	case MYSQL, MARIADB:
		return MYSQL
	case POSTGRES, POSTGRESQL, PSQL, PG:
		return POSTGRESQL
	case MSSQL, SQLSERVER, SQLSRV:
		return MSSQL
	case SQLITE, SQLITE3:
		return SQLITE
	case MONGODB, MONGO:
		return MONGODB
	default:
		return normalized
	}
}

// ===============================================================================================================
// supportedDatabaseDrivers maps recognized database driver names to their supported status boolean.
// ===============================================================================================================.
var supportedDatabaseDrivers = map[string]bool{
	MYSQL:      true,
	POSTGRES:   true,
	POSTGRESQL: true,
	MSSQL:      true,
	SQLSERVER:  true,
	SQLITE:     true,
	SQLITE3:    true,
	MONGO:      true,
	MONGODB:    true,
}
