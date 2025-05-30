package settings

import (
	"context"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

var AppVar *App

type App struct {
	DB          *gorm.DB
	Logger      *zap.Logger
	Config      *BaseConfig
	MongoDB     *mongo.Database
	RedisClient *redis.Client
	Mail        *gomail.Dialer
	Ctx         context.Context
	Cancel      context.CancelFunc
}

func NewApp(db *gorm.DB, logger *zap.Logger, config *BaseConfig, mongodb *mongo.Database, redis *redis.Client, mail *gomail.Dialer) *App {
	AppVar = &App{
		Ctx:         AppVar.Ctx,
		Cancel:      AppVar.Cancel,
		DB:          db,
		Logger:      logger,
		Config:      config,
		MongoDB:     mongodb,
		RedisClient: redis,
		Mail:        mail,
	}

	return AppVar
}
