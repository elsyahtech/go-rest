package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/elsyahtech/go-rest/config"
	"github.com/elsyahtech/go-rest/logger"
	"go.uber.org/zap"
)

func DatabaseMigration(ctx context.Context) error {
	// Get Global Configuration
	cfg := config.GlobalConfig

	if !cfg.Database.Migration {
		return nil
	}

	if cfg.Database.Driver == "" {
		logger.Logger(&logger.LoggerPayload{
			Fields: map[string]any{
				LogFieldKeyDriver: cfg.Database.Driver,
				LogFieldKeyTroubleshoot: "check your configuration file in ./app/config/database and " +
					"ensure a database driver (e.g., postgres, mysql, mssql, mongodb or sqlite) is specified under Database.Driver",
			},
		}).Error("database driver is not configured")

		return errors.New("database driver is not configured")
	}

	// Normalize driver name to folder name
	driver := normalizeDatabaseDriver(cfg.Database.Driver)

	// MongoDB doesn't use SQL migrations (it's NoSQL)
	switch driver {
	case MONGODB, MONGO:
		logger.Logger(&logger.LoggerPayload{}).Info("running MongoDB migrations")

		// Run MongoDB migrations
		if err := MongoDBMigration(cfg.MongoDBMigrations...); err != nil {
			logger.Logger(&logger.LoggerPayload{
				Fields: map[string]any{
					LogFieldKeyAdapter:      cfg.Database.Driver,
					LogFieldKeyError:        zap.Error(err),
					LogFieldKeyTroubleshoot: "verify MongoDB connection state, user permissions, and collection initialization scripts",
				},
			}).Error("MongoDB migration failed")

			return fmt.Errorf("MongoDB migration failed: %w", err)
		}

		return nil

	default:
		logger.Logger(&logger.LoggerPayload{}).Info("running SQL migrations")

		// Run SQL migrations
		if err := SQLMigration(ctx); err != nil {
			logger.Logger(&logger.LoggerPayload{
				Fields: map[string]any{
					LogFieldKeyAdapter:      cfg.Database.Driver,
					LogFieldKeyError:        zap.Error(err),
					LogFieldKeyTroubleshoot: "verify SQL connection state, user permissions, etc initialization scripts",
				},
			}).Error("SQL migration failed")

			return fmt.Errorf("SQL migration failed: %w", err)
		}

		return nil
	}
}
