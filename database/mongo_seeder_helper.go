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
	mongooptions "go.mongodb.org/mongo-driver/v2/mongo/options"
)

func createMongoDBSeederHistory() (string, error) {
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
		if name == "seeder_history" {
			return troubleShootMsg, nil
		}
	}

	if err := database.CreateCollection(ctx, "seeder_history"); err != nil {
		troubleShootMsg = "check database user privileges or ensure the database name is correct"
		return troubleShootMsg, fmt.Errorf("%w", err)
	}

	collection := database.Collection("seeder_history")

	_, err = collection.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys: bson.D{
				{
					Key:   "seeder_name",
					Value: 1,
				},
			},
			Options: mongooptions.Index().
				SetName("uk_seeder_name").
				SetUnique(true),
		},
	)
	if err != nil {
		troubleShootMsg = "check if an conflicting index with the same name already exists"
		return troubleShootMsg, fmt.Errorf("create unique index 'uk_seeder_name' on the 'seeder_history' collection failed: %w", err)
	}

	return troubleShootMsg, nil
}

func getMongoDBExecutedSeeders() ([]string, string, error) {
	troubleShootMsg := ""
	cfg := config.GlobalConfig
	globalDB := GlobalDB

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Database.Timeout)
	defer cancel()

	// Pastikan helper globalDB.Collection atau akses collection seeder_history tersedia
	coll, troubleShootMsg, err := globalDB.Collection("seeder_history")
	if err != nil {
		return nil, troubleShootMsg, fmt.Errorf("%w", err)
	}

	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		troubleShootMsg = "ensure the collection exists, the database connection is active, and you have read permissions"
		return nil, troubleShootMsg, fmt.Errorf("query the 'seeder_history' collection failed: %w", err)
	}

	defer func() {
		if err := cursor.Close(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "close rows failed: %v\n", err)
		}
	}()

	var rows []struct {
		SeederName string `bson:"seeder_name"`
	}

	if err := cursor.All(ctx, &rows); err != nil {
		troubleShootMsg = "check if the document schema matches the expected 'seeder_name' structure"
		return nil, troubleShootMsg, fmt.Errorf("decode/parse MongoDB seeder history documents failed: %w", err)
	}

	var executed []string
	for _, row := range rows {
		executed = append(executed, row.SeederName)
	}

	return executed, troubleShootMsg, nil
}

func filterMongoDBPendingSeeders(all []config.MongoDBSeeder, executed []string) []config.MongoDBSeeder {
	executedMap := make(map[string]struct{})

	for _, seeder := range executed {
		executedMap[seeder] = struct{}{}
	}

	var pending []config.MongoDBSeeder

	for _, seeder := range all {
		if _, exists := executedMap[seeder.Name()]; exists {
			continue
		}

		pending = append(pending, seeder)
	}

	return pending
}

func executeMongoDBSeeder(seeder config.MongoDBSeeder) (string, error) {
	troubleShootMsg := ""
	cfg := config.GlobalConfig
	globalDB := GlobalDB

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Database.Timeout)
	defer cancel()

	database := globalDB.MongoDB.Database(config.GlobalConfig.Database.Name)

	// Jalankan fungsi Up pada seeder yang dibuat sebelumnya
	if err := seeder.Up(ctx, database); err != nil {
		troubleShootMsg = "check the seeder script logic, database permissions, or syntax errors."
		return troubleShootMsg, fmt.Errorf("failed to execute the 'Up' method for MongoDB seeder: %w", err)
	}

	coll, troubleShootMsg, err := globalDB.Collection("seeder_history")
	if err != nil {
		return troubleShootMsg, fmt.Errorf("%w", err)
	}

	// Catat ke history bahwa seeder ini sudah sukses dijalankan
	_, err = coll.InsertOne(
		ctx,
		bson.M{
			"seeder_name": seeder.Name(),
			"executed_at": time.Now(),
		},
	)
	if err != nil {
		troubleShootMsg = fmt.Sprintf("seeder '%s' executed successfully, "+
			"but failed to record its history into the 'seeder_history' collection. "+
			"Check write permissions or unique index constraints.", seeder.Name())

		return troubleShootMsg, fmt.Errorf("%w", err)
	}

	return troubleShootMsg, nil
}

func getMongoDBSeederNames(seeders []config.MongoDBSeeder) []string {
	names := make([]string, 0, len(seeders))

	for _, seeder := range seeders {
		names = append(names, seeder.Name())
	}

	return names
}
