//nolint:dupl
package database

import (
	"fmt"

	"github.com/elsyahtech/go-rest/config"
	"github.com/elsyahtech/go-rest/logger"
	"go.uber.org/zap"
)

func MongoDBSeeder(registrars ...config.MongoDBSeederRegistrar) error {
	logger.Logger(&logger.LoggerPayload{}).Info("starting mongodb seeders")

	if len(registrars) == 0 {
		logger.Logger(&logger.LoggerPayload{}).Info("no mongodb seeders registered")
		return nil
	}

	registry := config.NewMongoDBSeederRegistry()

	for _, register := range registrars {
		register(registry)
	}

	seeders := registry.GetMongoDBSeeders()

	logger.Logger(&logger.LoggerPayload{
		Fields: map[string]any{
			LogFieldKeyCount:      len(seeders),
			LogFieldKeySeederName: getMongoDBSeederNames(seeders),
		},
	}).Info("registered mongodb seeders")

	// Buat koleksi pencatatan riwayat seeder (seeder_history)
	if troubleShootMsg, err := createMongoDBSeederHistory(); err != nil {
		logger.Logger(&logger.LoggerPayload{
			Fields: map[string]any{
				LogFieldKeyError:        zap.Error(err),
				LogFieldKeyTroubleshoot: troubleShootMsg,
			},
		}).Error("create mongodb seeder history failed")

		return fmt.Errorf("%w", err)
	}

	// Ambil daftar seeder yang sudah pernah dieksekusi sebelumnya
	executedSeeders, troubleShootMsg, err := getMongoDBExecutedSeeders()
	if err != nil {
		logger.Logger(&logger.LoggerPayload{
			Fields: map[string]any{
				LogFieldKeyError:        zap.Error(err),
				LogFieldKeyTroubleshoot: troubleShootMsg,
			},
		}).Error("fetch executed mongodb seeders failed")

		return fmt.Errorf("%w", err)
	}

	// Filter seeder yang belum pernah dijalankan saja
	pendingSeeders := filterMongoDBPendingSeeders(seeders, executedSeeders)

	if len(pendingSeeders) == 0 {
		logger.Logger(&logger.LoggerPayload{}).Info("database seeder ignored! all mongoDB seeders already executed")
		return nil
	}

	// Eksekusi setiap pending seeder secara berurutan
	for i, seeder := range pendingSeeders {
		logger.Logger(&logger.LoggerPayload{
			Fields: map[string]any{
				LogFieldKeyStep:       i + 1,
				LogFieldKeyTotal:      len(pendingSeeders),
				LogFieldKeySeederName: seeder.Name(),
			},
		}).Info("executing mongodb seeder")

		if troubleShootMsg, err := executeMongoDBSeeder(seeder); err != nil {
			logger.Logger(&logger.LoggerPayload{
				Fields: map[string]any{
					LogFieldKeySeederName:   seeder.Name(),
					LogFieldKeyError:        zap.Error(err),
					LogFieldKeyTroubleshoot: troubleShootMsg,
				},
			}).Error("seeder failed")

			return fmt.Errorf(
				"seeder %s failed: %w",
				seeder.Name(),
				err,
			)
		}

		logger.Logger(&logger.LoggerPayload{
			Fields: map[string]any{
				LogFieldKeySeederName: seeder.Name(),
			},
		}).Info("mongodb seeder completed")
	}

	logger.Logger(&logger.LoggerPayload{}).Info("all mongodb seeders completed")

	return nil
}
