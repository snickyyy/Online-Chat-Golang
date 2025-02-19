package settings

import (
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var AppVar *App

type App struct {
	DB 			*gorm.DB
	Logger 		*zap.Logger
	Config 		*BaseConfig
	MongoDB 	*mongo.Database
	RedisSess	*redis.Client
}

func NewApp(db *gorm.DB, logger *zap.Logger, config *BaseConfig, mongodb *mongo.Database, redisSess *redis.Client) *App {
	AppVar = &App{
        DB:     	db,
        Logger: 	logger,
        Config: 	config,
		MongoDB: 	mongodb,
		RedisSess:	redisSess,
    }
    return AppVar
}
