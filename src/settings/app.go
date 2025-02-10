package settings

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	DB 			*gorm.DB
	Logger 		*zap.Logger
	Config 		*BaseConfig
	MongoClient *mongo.Client
}

func NewApp(db *gorm.DB, logger *zap.Logger, config *BaseConfig, mongodb *mongo.Client) *App {
    return &App{
        DB:     db,
        Logger: logger,
        Config: config,
		MongoClient: mongodb,
    }
}
