package settings

import (

	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB
var Models = []interface{}{}

func InitDb() {
	dbConfig := GetBaseConfig().DatabaseConfig
	logger := GetLogger()
	defer logger.Sync()
	if dbConfig.Host == "" {
		logger.Error("Database host is not set")
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		dbConfig.Host, dbConfig.User, dbConfig.Password, dbConfig.Database, dbConfig.Port, dbConfig.Sslmode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		logger.Error("Database host is not set")
	}
	Db = db
}

func GetDb() *gorm.DB {
	return Db
}

func MakeMigrations() {
	GetDb().AutoMigrate(&Models)
}
