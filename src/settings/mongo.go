package settings

import (
	"context"
	"time"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMongoClient(config *BaseConfig) (*mongo.Client, error) {
	client, err := mongo.Connect(Context.Ctx, options.Client().ApplyURI(config.MongoConfig.Uri))
	if err != nil {
		return nil, err
	}
	return client, nil
}

func DisconnectMongoClient(app *App) {
	ctx, cancel := context.WithTimeout(Context.Ctx, time.Second*5)
	defer cancel()
	if err := app.MongoClient.Disconnect(ctx); err != nil {
		app.Logger.Error(fmt.Sprintf("error disconnecting from MongoDB: %v", err))
	}
	app.Logger.Info("disconnected from MongoDB")
}
