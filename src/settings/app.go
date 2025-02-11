package settings

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"go.mongodb.org/mongo-driver/mongo"
)

var AppVar *App

type App struct {
	DB 			*gorm.DB
	Logger 		*zap.Logger
	Config 		*BaseConfig
	MongoDB 	*mongo.Database
}

func NewApp(db *gorm.DB, logger *zap.Logger, config *BaseConfig, mongodb *mongo.Database) *App {
	AppVar = &App{
        DB:     		db,
        Logger: 		logger,
        Config: 		config,
		MongoDB: 	mongodb,
    }
    return AppVar
}
