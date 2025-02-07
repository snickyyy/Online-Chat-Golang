package settings

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB
var Models = []interface{}{}

func InitDb(baseConfig *BaseConfig) (*gorm.DB, error){
	
	dbConfig := baseConfig.DatabaseConfig
	if dbConfig.Host == "" {
		return nil, fmt.Errorf("db host is not set")
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		dbConfig.Host, dbConfig.User, dbConfig.Password, dbConfig.Database, dbConfig.Port, dbConfig.Sslmode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		return nil, fmt.Errorf("db connection error: %v", err)
	}
	Db = db
	return db, nil
}

func GetDb() *gorm.DB {
	return Db
}

func MakeMigrations() {
	GetDb().AutoMigrate(&Models)
}
