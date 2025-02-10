package settings

import (
	"fmt"

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
			GetContext,
			NewApp,
		),
		fx.Invoke(func(app *App) {
			fmt.Println("App initialized:", app)
		},
		MakeMigrations),
	)
	return di
}
