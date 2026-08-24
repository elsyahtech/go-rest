package database

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/elsyahtech/go-rest/config"
	"github.com/elsyahtech/go-rest/logger"

	"go.uber.org/zap"
)

// MigrationFile represents a single migration SQL file.
type MigrationFile struct {
	CreatedAt time.Time
	Name      string
	Path      string
	Content   string
	Number    int
}

// MigrationHistory tracks executed migrations.
type MigrationHistory struct {
	ExecutedAt    time.Time
	MigrationName string
}

func SQLMigration(ctx context.Context) error {
	// Get database configuration
	cfg := config.GlobalConfig

	// Normalize driver name to folder name
	driver := normalizeDatabaseDriver(cfg.Database.Driver)

	// Build migration directory path for SQL databases
	migrationDir := filepath.Join(".", "app", "database", "migrations", driver)

	logger.Logger(&logger.LoggerPayload{
		Fields: map[string]any{
			LogFieldKeyDriver:       driver,
			LogFieldKeyMigrationDir: migrationDir,
		},
	}).Info("starting database migration")

	// Step 1: Read all migration files from directory
	migrationFiles, troubleShootMsg, err := readMigrationFiles(migrationDir)
	if err != nil {
		logger.Logger(&logger.LoggerPayload{
			Fields: map[string]any{
				LogFieldKeyError:        zap.Error(err),
				LogFieldKeyTroubleshoot: troubleShootMsg,
			},
		}).Error("read migration files")

		return fmt.Errorf("read migration files failed: %w", err)
	}

	// No migrations found, that's OK (first run or no migrations needed)
	if len(migrationFiles) == 0 {
		logger.Logger(&logger.LoggerPayload{
			Fields: map[string]any{
				LogFieldKeyDriver:       driver,
				LogFieldKeyMigrationDir: migrationDir,
			},
		}).Info("no migration files found")

		return nil
	}

	// Log found migrations
	logger.Logger(&logger.LoggerPayload{
		Fields: map[string]any{
			LogFieldKeyCount:  len(migrationFiles),
			LogFieldKeyDriver: driver,
			LogFieldKeyFile:   getMigrationNames(migrationFiles),
		},
	}).Info("found migration files")

	// Step 2: Create migration history table (if not exists)
	if troubleShootMsg, err := createMigrationHistoryTable(ctx); err != nil {
		logger.Logger(&logger.LoggerPayload{
			Fields: map[string]any{
				LogFieldKeyError:        zap.Error(err),
				LogFieldKeyTroubleshoot: troubleShootMsg,
			},
		}).Error("create migration history table failed")

		return fmt.Errorf("create migration history table failed: %w", err)
	}

	// Step 3: Get list of already-executed migrations
	executedMigrations, troubleShootMsg, err := getExecutedMigrations(ctx)
	if err != nil {
		logger.Logger(&logger.LoggerPayload{
			Fields: map[string]any{
				LogFieldKeyError:        zap.Error(err),
				LogFieldKeyTroubleshoot: troubleShootMsg,
			},
		}).Error("get executed migrations failed")

		return fmt.Errorf("get executed migrations failed: %w", err)
	}

	if len(executedMigrations) > 0 {
		logger.Logger(&logger.LoggerPayload{
			Fields: map[string]any{
				LogFieldKeyCount: len(executedMigrations),
				LogFieldKeyFile:  executedMigrations,
			},
		}).Info("found executed migrations")
	}

	// Step 4: Filter out already-executed migrations (find pending)
	pendingMigrations := filterPendingMigrations(migrationFiles, executedMigrations)

	// All migrations already executed, nothing to do
	if len(pendingMigrations) == 0 {
		logger.Logger(&logger.LoggerPayload{
			Fields: map[string]any{
				LogFieldKeyDriver: driver,
			},
		}).Info("database migration ignored! all migrations were already executed")

		return nil
	}

	// Log pending migrations
	logger.Logger(&logger.LoggerPayload{
		Fields: map[string]any{
			LogFieldKeyCount: len(pendingMigrations),
			LogFieldKeyFile:  getMigrationNames(pendingMigrations),
		},
	}).Info("found pending migrations")

	// Step 5: Execute each pending migration in order
	for i, migration := range pendingMigrations {
		logger.Logger(&logger.LoggerPayload{
			Fields: map[string]any{
				LogFieldKeyStep:          i + 1,
				LogFieldKeyTotal:         len(pendingMigrations),
				LogFieldKeyMigrationName: migration.Name,
			},
		}).Info("executing migration")

		if troubleShootMsg, err := executeMigration(ctx, migration); err != nil {
			logger.Logger(&logger.LoggerPayload{
				Fields: map[string]any{
					LogFieldKeyMigrationName: migration.Name,
					LogFieldKeyError:         zap.Error(err),
					LogFieldKeyTroubleshoot:  troubleShootMsg,
				},
			}).Error("migration failed")

			return fmt.Errorf("migration %s failed: %w", migration.Name, err)
		}

		logger.Logger(&logger.LoggerPayload{
			Fields: map[string]any{
				LogFieldKeyMigrationName: migration.Name,
			},
		}).Info("migration executed successfully")
	}

	// All pending migrations executed successfully
	logger.Logger(&logger.LoggerPayload{
		Fields: map[string]any{
			LogFieldKeyCount:  len(pendingMigrations),
			LogFieldKeyDriver: driver,
		},
	}).Info("all migrations completed successfully")

	return nil
}
