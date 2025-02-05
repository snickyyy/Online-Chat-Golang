package settings

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"fmt"
)

var Db *gorm.DB
var Models = []interface{}{}

func InitDb() {
	dbConfig := GetBaseConfig().DatabaseConfig
	if dbConfig.Host == "" {
		log.Fatal("Database config is not set")
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		dbConfig.Host, dbConfig.User, dbConfig.Password, dbConfig.Database, dbConfig.Port, dbConfig.Sslmode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	Db = db
}

func GetDb() *gorm.DB {
	return Db
}

func MakeMigrations() {
	GetDb().AutoMigrate(&Models)
}
