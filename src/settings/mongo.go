package settings

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMongoClient(config *BaseConfig, dbName string) (*mongo.Database, error) {
	client, err := mongo.Connect(AppVar.Ctx, options.Client().ApplyURI(config.MongoConfig.Uri))
	if err != nil {
		return nil, err
	}

	if err = client.Ping(AppVar.Ctx, nil); err != nil {
		return nil, err
	}

	db := client.Database(dbName)
	return db, nil
}
