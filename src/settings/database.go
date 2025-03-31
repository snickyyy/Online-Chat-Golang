package settings

import (
	"fmt"
	"libs/src/internal/domain/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Models = []interface{}{
	&domain.User{},
	&domain.Chat{},
	&domain.ChatMember{},
}

func GetDb(baseConfig *BaseConfig) (*gorm.DB, error) {

	dbConfig := baseConfig.PostgresConfig
	if dbConfig.Host == "" {
		return nil, fmt.Errorf("db host is not set")
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		dbConfig.Host, dbConfig.User, dbConfig.Password, dbConfig.Database, dbConfig.Port, dbConfig.Sslmode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		return nil, fmt.Errorf("db connection error: %v", err)
	}
	return db, nil
}

func MakeMigrations(app *App) {
	app.Logger.Info("Migrating models")
	app.DB.AutoMigrate(Models...)
}
