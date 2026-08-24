//nolint:dupl
package database

import (
	"fmt"

	"github.com/elsyahtech/go-rest/config"
	"github.com/elsyahtech/go-rest/logger"

	"go.uber.org/zap"
)

// ======================================================================================
// MongoDBMigration orchestrates the execution of registered MongoDB schema migrations.
// It handles initialization of the migration registry, checks migration history,
// filters out already-executed migrations, and executes pending migrations sequentially.
// ======================================================================================.
func MongoDBMigration(registrars ...config.MongoDBMigrationRegistrar) error {
	// Log the initiation of the MongoDB migration process
	logger.Logger(&logger.LoggerPayload{}).Info("starting mongodb migrations")

	// Return early if no migration registrars were provided in the configuration
	if len(registrars) == 0 {
		logger.Logger(&logger.LoggerPayload{}).Info("no mongodb migrations registered")

		return nil
	}

	// Initialize a new migration registry to collect all structured migrations
	registry := config.NewMigrationRegistry()

	// Execute each registration function to populate the registry collection
	for _, register := range registrars {
		register(registry)
	}

	// Retrieve the complete slice of collected MongoDB migration objects
	migrations := registry.GetMigrations()

	// Log the summary of successfully registered migrations along with their names/identifiers
	logger.Logger(&logger.LoggerPayload{
		Fields: map[string]any{
			LogFieldKeyCount:         len(migrations),
			LogFieldKeyMigrationName: getMongoDBMigrationNames(migrations),
		},
	}).Info("registered mongodb migrations")

	// Create the migration history tracking structure or collection in the database if it doesn't exist
	if troubleShootMsg, err := createMongoDBMigrationHistory(); err != nil {
		logger.Logger(&logger.LoggerPayload{
			Fields: map[string]any{
				LogFieldKeyError:        zap.Error(err),
				LogFieldKeyTroubleshoot: troubleShootMsg,
			},
		}).Error("create mongodb migration history failed")

		return fmt.Errorf("%w", err)
	}

	// Fetch the list of migrations that have already been executed from the database history
	executedMigrations, troubleShootMsg, err := getMongoDBExecutedMigrations()
	if err != nil {
		logger.Logger(&logger.LoggerPayload{
			Fields: map[string]any{
				LogFieldKeyError:        zap.Error(err),
				LogFieldKeyTroubleshoot: troubleShootMsg,
			},
		}).Error("execute mongodb migration failed")

		return fmt.Errorf("%w", err)
	}

	// Filter out already executed migrations, leaving only the pending ones to be run
	pendingMigrations := filterMongoDBPendingMigrations(migrations, executedMigrations)

	// Skip the migration step if no new pending migrations are found
	if len(pendingMigrations) == 0 {
		logger.Logger(&logger.LoggerPayload{}).Info("database migration ignored! all mongoDB migrations already executed")

		return nil
	}

	// Iterate through and execute each pending migration in sequential order
	for i, migration := range pendingMigrations {
		// Log the current step and details of the migration about to be executed
		logger.Logger(&logger.LoggerPayload{
			Fields: map[string]any{
				LogFieldKeyStep:          i + 1,
				LogFieldKeyTotal:         len(pendingMigrations),
				LogFieldKeyMigrationName: migration.Name(),
			},
		}).Info("executing mongodb migration")

		// Perform the migration operation (Up method and history recording)
		if troubleShootMsg, err := executeMongoDBMigration(migration); err != nil {
			logger.Logger(&logger.LoggerPayload{
				Fields: map[string]any{
					LogFieldKeyMigrationName: migration.Name(),
					LogFieldKeyError:         zap.Error(err),
					LogFieldKeyTroubleshoot:  troubleShootMsg,
				},
			}).Error("migration failed")

			return fmt.Errorf(
				"migration %s failed: %w",
				migration.Name(),
				err,
			)
		}

		// Log confirmation upon successful completion of the individual migration
		logger.Logger(&logger.LoggerPayload{
			Fields: map[string]any{
				LogFieldKeyMigrationName: migration.Name(),
			},
		}).Info("mongodb migration completed")
	}

	// Log the final success message indicating all pending migrations have finished executing
	logger.Logger(&logger.LoggerPayload{
		Fields: map[string]any{
			LogFieldKeyMigrationName: migrations,
		},
	}).Info("all mongodb migrations completed")

	return nil
}
