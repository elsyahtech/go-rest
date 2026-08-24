//nolint:dupl
package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/elsyahtech/go-rest/config"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func createMongoDBMigrationHistory() (string, error) {
	troubleShootMsg := ""
	cfg := config.GlobalConfig
	globalDB := GlobalDB

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Database.Timeout)

	defer cancel()

	database := globalDB.MongoDB.Database(config.GlobalConfig.Database.Name)

	collectionNames, err := database.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		troubleShootMsg = "verify that the database connection is active, network is reachable, and the user has 'listCollections' permissions"

		return troubleShootMsg, fmt.Errorf("%w", err)
	}

	for _, name := range collectionNames {
		if name == "migration_history" {
			return troubleShootMsg, nil
		}
	}

	if err := database.CreateCollection(ctx, "migration_history"); err != nil {
		troubleShootMsg = "check database user privileges or ensure the database name is correct"

		return troubleShootMsg, fmt.Errorf("%w", err)
	}

	collection := database.Collection("migration_history")

	_, err = collection.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys: bson.D{
				{
					Key:   "migration_name",
					Value: 1,
				},
			},
			Options: options.Index().
				SetName(
					"uk_migration_name",
				).
				SetUnique(true),
		},
	)
	if err != nil {
		troubleShootMsg = "check if an conflicting index with the same name already exists"

		return troubleShootMsg, fmt.Errorf("create unique index 'uk_migration_name' on the 'migration_history' collection failed: %w", err)
	}

	return troubleShootMsg, nil
}

func getMongoDBExecutedMigrations() ([]string, string, error) {
	troubleShootMsg := ""
	cfg := config.GlobalConfig
	globalDB := GlobalDB

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Database.Timeout)

	defer cancel()

	coll, troubleShootMsg, err := globalDB.Collection("migration_history")
	if err != nil {
		return nil, troubleShootMsg, fmt.Errorf("%w", err)
	}

	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		troubleShootMsg := "ensure the collection exists, the database connection is active, and you have read permissions"

		return nil, troubleShootMsg, fmt.Errorf("query the 'migration_history' collection failed: %w", err)
	}

	defer func() {
		if err := cursor.Close(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "close rows failed: %v\n", err)
		}
	}()

	var rows []struct {
		MigrationName string `bson:"migration_name"`
	}

	if err := cursor.All(ctx, &rows); err != nil {
		troubleShootMsg := "check if the document schema matches the expected 'migration_name' structure"

		return nil, troubleShootMsg, fmt.Errorf("decode/parse MongoDB migration history documents failed: %w", err)
	}

	var executed []string

	for _, row := range rows {
		executed = append(
			executed,
			row.MigrationName,
		)
	}

	return executed, troubleShootMsg, nil
}

func filterMongoDBPendingMigrations(all []config.MongoDBMigration, executed []string) []config.MongoDBMigration {
	executedMap := make(
		map[string]struct{},
	)

	for _, migration := range executed {
		executedMap[migration] = struct{}{}
	}

	var pending []config.MongoDBMigration

	for _, migration := range all {
		if _, exists :=
			executedMap[migration.Name()]; exists {
			continue
		}

		pending = append(
			pending,
			migration,
		)
	}

	return pending
}

func executeMongoDBMigration(migration config.MongoDBMigration) (string, error) {
	troubleShootMsg := ""
	cfg := config.GlobalConfig
	globalDB := GlobalDB

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Database.Timeout)

	defer cancel()

	database := globalDB.MongoDB.Database(config.GlobalConfig.Database.Name)

	if err := migration.Up(ctx, database); err != nil {
		troubleShootMsg = "check the migration script logic, database permissions, or syntax errors."

		return troubleShootMsg, fmt.Errorf("failed to execute the 'Up' method for MongoDB migration: %w", err)
	}

	coll, troubleShootMsg, err := globalDB.Collection("migration_history")
	if err != nil {
		return troubleShootMsg, fmt.Errorf("%w", err)
	}

	_, err = coll.InsertOne(
		ctx,
		bson.M{
			"migration_name": migration.Name(),
			"executed_at":    time.Now(),
		},
	)
	if err != nil {
		troubleShootMsg = fmt.Sprintf("migration '%s' executed successfully, "+
			"but failed to record its history into the 'migration_history' collection. "+
			"Check write permissions or unique index constraints.", migration.Name())

		return troubleShootMsg, fmt.Errorf("%w", err)
	}

	return troubleShootMsg, nil
}

func getMongoDBMigrationNames(migrations []config.MongoDBMigration) []string {
	names := make(
		[]string,
		0,
		len(migrations),
	)

	for _, migration := range migrations {
		names = append(
			names,
			migration.Name(),
		)
	}

	return names
}
