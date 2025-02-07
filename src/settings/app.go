package settings

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type App struct {
	DB 		*gorm.DB
	Logger 	*zap.Logger
	Config 	*BaseConfig
}

func NewApp(db *gorm.DB, logger *zap.Logger, config *BaseConfig) *App {
    return &App{
        DB:     db,
        Logger: logger,
        Config: config,
    }
}
