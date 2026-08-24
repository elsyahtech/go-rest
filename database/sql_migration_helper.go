package database

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/elsyahtech/go-rest/config"
)

// readMigrationFiles reads all SQL files from the migration directory
// Only reads .sql files, ignores directories and other file types
// Returns migrations sorted by number (001, 002, 003, etc)
// If directory doesn't exist, returns empty list (not an error).
func readMigrationFiles(migrationDir string) ([]MigrationFile, string, error) {
	troubleShootMsg := ""

	// Try to read directory
	entries, err := os.ReadDir(migrationDir)
	if err != nil {
		troubleShootMsg = fmt.Sprintf("ensure the directory %s exists, has proper read permissions, "+
			"and the path is correctly configured.", migrationDir)

		// Directory doesn't exist is OK (no migrations yet)
		if os.IsNotExist(err) {
			return make([]MigrationFile, 0), troubleShootMsg, nil
		}

		return nil, troubleShootMsg, fmt.Errorf("read directory %s failed: %w", migrationDir, err)
	}

	var migrationFiles []MigrationFile

	// Process each file in directory
	for _, entry := range entries {
		// Skip directories
		if entry.IsDir() {
			continue
		}

		// Only process .sql files
		if !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		// Read file content
		filePath := filepath.Join(migrationDir, entry.Name())

		content, err := os.ReadFile(filePath)
		if err != nil {
			troubleShootMsg = fmt.Sprintf("check file %s read permissions or if the file is currently locked/corrupted.", filePath)

			return nil, troubleShootMsg, fmt.Errorf("read file %s failed: %w", filePath, err)
		}

		// Create migration file struct
		migration := MigrationFile{
			Name:    entry.Name(),
			Path:    filePath,
			Content: string(content),
			Number:  extractMigrationNumber(entry.Name()),
		}

		migrationFiles = append(migrationFiles, migration)
	}

	// Sort by number (001 < 002 < 003, etc)
	// Ensures migrations execute in correct order
	slices.SortFunc(migrationFiles, func(a, b MigrationFile) int {
		if a.Number < b.Number {
			return -1
		}

		if a.Number > b.Number {
			return 1
		}

		return 0
	})

	return migrationFiles, troubleShootMsg, nil
}

// extractMigrationNumber extracts the number prefix from migration filename
// Example: "001_create_users_table.sql" → 1
// Example: "042_add_column_to_products.sql" → 42
// Used for sorting migrations in execution order.
func extractMigrationNumber(filename string) int {
	// Remove .sql extension
	name := strings.TrimSuffix(filename, ".sql")

	// Split by underscore, first part is the number
	parts := strings.Split(name, "_")
	if len(parts) == 0 {
		return 0
	}

	// Parse first part as integer
	var number int
	fmt.Sscanf(parts[0], "%d", &number) //nolint:errcheck

	return number
}

// createMigrationHistoryTable creates migration_history table if it doesn't exist
// This table tracks which migrations have been executed
// Syntax differs by database driver (MySQL, PostgreSQL, MSSQL)
// MongoDB doesn't use SQL history (skipped).
func createMigrationHistoryTable(ctx context.Context) (string, error) {
	troubleShootMsg := ""
	driver := strings.ToLower(config.GlobalConfig.Database.Driver)

	var createTableSQL string

	// Build CREATE TABLE SQL based on driver
	//nolint:revive
	switch driver {
	case MYSQL, MARIADB:
		createTableSQL = `
		CREATE TABLE IF NOT EXISTS migration_history (
			id INT AUTO_INCREMENT PRIMARY KEY,
			migration_name VARCHAR(255) NOT NULL UNIQUE,
			executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
		`
	case POSTGRESQL, POSTGRES, PSQL, PG:
		createTableSQL = `
		CREATE TABLE IF NOT EXISTS migration_history (
			id SERIAL PRIMARY KEY,
			migration_name VARCHAR(255) NOT NULL UNIQUE,
			executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
		`
	case MSSQL, SQLSERVER, SQLSRV:
		createTableSQL = `
		IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'migration_history')
		CREATE TABLE migration_history (
			id INT PRIMARY KEY IDENTITY(1,1),
			migration_name VARCHAR(255) NOT NULL UNIQUE,
			executed_at DATETIME DEFAULT GETDATE()
		)
		`
	case SQLITE, SQLITE3:
		createTableSQL = `
		CREATE TABLE IF NOT EXISTS migration_history (
			id INT AUTO_INCREMENT PRIMARY KEY,
			migration_name VARCHAR(255) NOT NULL UNIQUE,
			executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
		`
	case "mongodb", "mongo":
		// MongoDB does not use an SQL history table (NoSQL).
		// MongoDB migrations run separately on mongodb_migration.
		troubleShootMsg = ""

		return troubleShootMsg, nil
	default:
		troubleShootMsg = fmt.Sprintf("the database driver '%s' is not supported for SQL-based migrations. "+
			"Check your configuration in ./app/config/database and use a supported driver (mysql, postgres, sqlite, mssql).", driver)

		return troubleShootMsg, fmt.Errorf("unsupported database driver: %s", driver)
	}

	db := GlobalDB

	// Execute CREATE TABLE statement
	if troubleShootMsg, err := db.Exec(ctx, createTableSQL); err != nil {
		msg := fmt.Sprintf("failed to execute SQL statement for creating the 'migration_history' table using driver '%s'. "+
			"Verify database user privileges, connection state, syntax compatibility or %s", driver, troubleShootMsg)

		return msg, fmt.Errorf("%w", err)
	}

	return troubleShootMsg, nil
}

// getExecutedMigrations retrieves list of already-executed migrations from database
// Queries migration_history table and returns migration filenames
// Returns empty list if table doesn't exist (first run).
func getExecutedMigrations(ctx context.Context) ([]string, string, error) {
	troubleShootMsg := ""
	db := GlobalDB

	query := "SELECT migration_name FROM migration_history ORDER BY executed_at ASC"

	// Execute query
	rows, troubleShootMsg, err := db.Query(ctx, query)
	if err != nil {
		msg := fmt.Sprintf("query the 'migration_history' table failed. Ensure that the table exists, "+
			"the database connection is healthy, and the SQL query syntax is correct or %s", troubleShootMsg)

		return nil, msg, fmt.Errorf("%w", err)
	}

	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "close rows failed: %v\n", err)
		}
	}()

	var executedMigrations []string

	// Scan each row and collect migration names
	for rows.Next() {
		var migrationName string

		if err := rows.Scan(&migrationName); err != nil {
			troubleShootMsg = "scan migration name from the result set failed. " +
				"Check if the 'migration_name' column data type matches string expectations"

			return nil, troubleShootMsg, fmt.Errorf("scan migration name failed: %w", err)
		}

		executedMigrations = append(executedMigrations, migrationName)
	}

	// Check for iteration errors
	if err = rows.Err(); err != nil {
		troubleShootMsg = "an error occurred while iterating through the migration rows. " +
			"Check network stability or potential database cursor timeouts"

		return nil, troubleShootMsg, fmt.Errorf("error iterating migration rows: %w", err)
	}

	return executedMigrations, troubleShootMsg, nil
}

// filterPendingMigrations returns only migrations not yet executed
// Compares all migration files with executed migrations list
// Returns migrations that need to be run.
func filterPendingMigrations(all []MigrationFile, executed []string) []MigrationFile {
	var pending []MigrationFile

	for _, migration := range all {
		// Check if migration is in executed list
		if !slices.Contains(executed, migration.Name) {
			// Migration not executed yet, add to pending
			pending = append(pending, migration)
		}
	}

	return pending
}

// executeMigration executes a single migration SQL file and records it in history
// Executes the SQL, then inserts record in migration_history table
// If any step fails, returns error (entire migration considered failed).
func executeMigration(ctx context.Context, migration MigrationFile) (string, error) {
	troubleShootMsg := ""
	database := GlobalDB

	const (
		queryPostgres = "INSERT INTO migration_history (migration_name, executed_at) VALUES ($1, $2)"
		queryMssql    = "INSERT INTO migration_history (migration_name, executed_at) VALUES (@p1, @p2)"
		queryMysql    = "INSERT INTO migration_history (migration_name, executed_at) VALUES (?, ?)"
		querySqlite   = "INSERT INTO migration_history (migration_name, executed_at) VALUES (?, ?)"
	)

	if troubleShootMsg, err := database.Exec(ctx, migration.Content); err != nil {
		msg := fmt.Sprintf("execute SQL statements inside migration file '%s' failed. "+
			"Check the migration script syntax, schema constraints, or "+
			"potential SQL runtime errors or %s", migration.Name, troubleShootMsg)

		return msg, fmt.Errorf("execute migration SQL failed: %w", err)
	}

	driver := strings.ToLower(config.GlobalConfig.Database.Driver)

	var insertQuery string

	//nolint:revive,nolintlint
	switch driver {
	case POSTGRESQL, POSTGRES, PSQL, PG:
		insertQuery = queryPostgres
	case MSSQL, SQLSERVER, SQLSRV:
		insertQuery = queryMssql
	case MYSQL, MARIADB:
		insertQuery = queryMysql
	case SQLITE, SQLITE3:
		insertQuery = querySqlite
	default:
		troubleShootMsg = fmt.Sprintf("the database driver '%s' is unsupported for recording migration history. "+
			"Check your configuration driver in ./app/config/database.", driver)

		return troubleShootMsg, fmt.Errorf("unsupported database driver: %s", driver)
	}

	// Record migration in migration_history table (mark as executed)
	if troubleShootMsg, err := database.Exec(ctx, insertQuery, migration.Name, time.Now()); err != nil {
		msg := fmt.Sprintf("migration '%s' executed successfully, but failed to record its history into the 'migration_history' table. "+
			"Check database write permissions or unique index constraints or %s", migration.Name, troubleShootMsg)

		return msg, fmt.Errorf("record migration in history failed: %w", err)
	}

	return troubleShootMsg, nil
}

// getMigrationNames returns list of migration filenames (for logging)
// Converts MigrationFile slice to string slice containing just filenames
// Used for logging which migrations were found/executed.
func getMigrationNames(migrations []MigrationFile) []string {
	names := make([]string, 0, len(migrations))

	for _, m := range migrations {
		names = append(names, m.Name)
	}

	return names
} //nolint:revive
