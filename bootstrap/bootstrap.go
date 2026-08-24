package bootstrap

import (
	"context"
	"fmt"
	"log"

	"github.com/elsyahtech/go-rest/database"
	"github.com/elsyahtech/go-rest/logger"
	"go.uber.org/zap"
)

func AppRun() error {
	// ==================================(2)==================================
	// =============================INIT TIMEZONE=============================
	// Set the application's default timezone.
	// This loads and sets the application's timezone from configuration.
	//
	// The timezone should be configured in the application config file.
	// Example: "Asia/Jakarta" or "UTC"
	if err := InitTimezone(); err != nil {
		return fmt.Errorf("load timezone failed, %w", err)
	}

	// =====================================(3)==================================
	// ============================LOAD RUN MODE TESTING=========================
	// Run commands that are specific to the current application mode.
	// Different application modes (test, development, production) may require different initialization steps.
	// This section executes mode-specific commands BEFORE the main application starts.
	if err := RunTesting(); err != nil {
		return fmt.Errorf("run testing failed, %w", err)
	}

	// =====================================(4)==================================
	// =================================INIT LOGGER==============================
	// Initialize the logger
	// This sets up the global logger instance that will be used throughout the application.
	// The logger is configured based on APP_MODE environment variable:
	// - "production": logs to files (./writable/logs/info.log and error.log)
	// - "test": logs to console with DebugLevel
	// - default: logs to console with InfoLevel
	//
	// NOTE: Do NOT call Sync() in individual handlers or services.
	// The Sync() call below ensures all buffered logs are flushed before the application exits.
	InitLogger()

	defer func() {
		if err := logger.GlobalLogger.Sync(); err != nil {
			log.Printf("close logger sync failed, %v", err)
		}
	}()

	// =======================================(5)==================================
	// =================================INIT DATABASE==============================
	// Open the database connection for dependency injection.
	// Establishes connections to the configured databases based on ./app/config/database.
	// This application supports multiple database connections:
	// - MySQL: Primary relational database
	// - PostgreSQL: Secondary/alternative relational database
	// - SQLite: Secondary/alternative relational database
	// - MSSQL: Enterprise SQL Server support
	// - MongoDB: Secondary/alternative no-SQL database
	if troubleShootMsg, err := InitDatabase(); err != nil {
		logger.Logger(&logger.LoggerPayload{
			Fields: map[string]any{
				LogFieldKeyError:        zap.Error(err),
				LogFieldKeyTroubleshoot: troubleShootMsg,
			},
		}).Error("database initialization failed")

		return fmt.Errorf("database initialization failed, %w", err)
	}

	defer func() {
		if troubleShootMsg, err := database.GlobalDB.Close(); err != nil {
			logger.Logger(&logger.LoggerPayload{
				Fields: map[string]any{
					LogFieldKeyError:        zap.Error(err),
					LogFieldKeyTroubleshoot: troubleShootMsg,
				},
			}).Error("close database connection failed")
		}
	}()

	// =====================================(6)==================================
	// =============================LOAD DB MIGRATION============================
	// Execute database migrations
	if err := database.DatabaseMigration(context.Background()); err != nil {
		return fmt.Errorf("db migration failed, %w", err)
	}

	// =====================================(7)==================================
	// =================================INIT ROUTER==============================
	// Initialize and load all HTTP routes for the application.
	// This step registers all route handlers from the modules defined in GlobalConfig.
	// The router is configured with CORS middleware and all application routes.
	//
	// Routes are loaded from: ./app/config/router.go
	// Router configuration (CORS, allowed origins, etc) comes from: GlobalConfig.Router
	if troubleShootMsg, err := InitRouter(); err != nil {
		logger.Logger(&logger.LoggerPayload{
			Fields: map[string]any{
				LogFieldKeyError:        zap.Error(err),
				LogFieldKeyTroubleshoot: troubleShootMsg,
			},
		}).Error("initialize HTTP router failed")

		return fmt.Errorf("initialize HTTP router failed, %w", err)
	}

	// =====================================(8)==================================
	// ================================START SERVER==============================
	// Start the HTTP server on the configured port.
	// This is the FINAL step before the application becomes fully operational
	if troubleShootMsg, err := StartServer(); err != nil {
		logger.Logger(&logger.LoggerPayload{
			Fields: map[string]any{
				LogFieldKeyError:        zap.Error(err),
				LogFieldKeyTroubleshoot: troubleShootMsg,
			},
		}).Error("start server failed")

		return fmt.Errorf("start server failed, %w", err)
	}

	return nil
}
