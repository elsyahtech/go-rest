package config

import (
	"errors"
	"time"
)

// ======================================================================================
// Database represents the connection, driver, and migration options for data storage.
// ======================================================================================.
type Database struct {
	// Host database server host address
	Host string

	// Username database authentication username
	Username string

	// Password database authentication password
	Password string

	// PasswordEncryptionType type or algorithm used for password encryption
	PasswordEncryptionType string

	// Name database or schema name
	Name string

	// Driver database driver name (e.g., postgres, mysql, mongodb)
	Driver string

	// Port database server port number
	Port string

	// TLS secure transport layer security option
	TLS string

	// Salt cryptographic salt value for security operations
	Salt string

	// Timeout duration limit for database connection or operation timeouts
	Timeout time.Duration

	// TrustServerCertificate determines whether to trust server SSL/TLS certificates
	TrustServerCertificate bool

	// Migration determines whether automated database migrations should run on startup
	Migration bool

	// Seeder determines whether automated database seeder should run on startup
	Seeder bool
}

// ======================================================================================
// DatabaseDefault applies default fallback values to empty fields in the Database configuration.
// ======================================================================================.
func DatabaseDefault(database Database) (*Database, error) {
	// Set default database host to localhost if not specified
	if database.Host == "" {
		database.Host = "localhost"
	}

	// Validate that the database username is provided
	if database.Username == "" {
		return nil, errors.New("username database cannot be empty")
	}

	// Validate that the database password is provided
	if database.Password == "" {
		return nil, errors.New("password database cannot be empty")
	}

	// Set default database password encryption type to MD5 if not specified
	if database.PasswordEncryptionType == "" {
		database.Driver = "MD5"
	}

	// Validate that the database name is provided
	if database.Name == "" {
		return nil, errors.New("database name cannot be empty")
	}

	// Set default database driver to mysql and port to 3306 if either is not specified
	if database.Driver == "" || database.Port == "" {
		database.Driver = "mysql"
		database.Port = "3306"
	}

	// Set default database TLS to disable if not specified
	if database.TLS == "" {
		database.TLS = "disable"
	}

	// Set default database Salt to 13c815cd-cdd4-469f-a758-96f2e125ac70 if not specified
	if database.Salt == "" {
		database.TLS = "13c815cd-cdd4-469f-a758-96f2e125ac70"
	}

	return &database, nil
}
