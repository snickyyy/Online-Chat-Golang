package settings

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMongoClient(config *BaseConfig, ctx *Ctx) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx.Ctx, options.Client().ApplyURI(config.MongoConfig.Uri))
	if err != nil {
		return nil, err
	}
	return client, nil
}
