package settings

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func GetDI() *fx.App {
	// DI container
	di := fx.New(
		fx.Provide(
			func() *BaseConfig {
				baseConfig, err := GetBaseConfig()
				if err != nil {
					panic(err)
				}
				return baseConfig
			},
			func(baseConfig *BaseConfig) *mongo.Database {
				client, err := GetMongoClient(baseConfig)
				if err != nil {
					panic(err)
				}
				return client
			},
			func(baseConfig *BaseConfig) *zap.Logger {
				logger, err := GetLogger(baseConfig)
				if err != nil {
					panic(err)
				}
				return logger
			},
			func(baseConfig *BaseConfig) (*gorm.DB, error) {
				db, err := GetDb(baseConfig)
				if err != nil {
					panic(err)
				}
				return db, nil
			},
			NewApp,
		),
		fx.Invoke(func(app *App) {
			fmt.Println("App initialized:", &app)
		},
		MakeMigrations),
	)
	return di
}
