package integration

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"libs/src/settings"
)

func SetupTestDatabase(dbConfig settings.PostgresConfig, testDbConfig settings.PostgresConfig) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		dbConfig.Host, dbConfig.User, dbConfig.Password, dbConfig.Database, dbConfig.Port, dbConfig.Sslmode,
	)

	baseDb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		panic(err)
	}

	baseDb.Exec("CREATE DATABASE test_db")
	sqlDB, err := baseDb.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.Close()

	testDsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		testDbConfig.Host, testDbConfig.User, testDbConfig.Password, testDbConfig.Database, testDbConfig.Port, testDbConfig.Sslmode,
	)
	testDB, err := gorm.Open(postgres.Open(testDsn), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		panic(err)
	}

	return testDB
}

func DropTestDatabase(dbConfig settings.PostgresConfig) error {
	oldDB, err := settings.AppVar.DB.DB()
	if err != nil {
		panic(err)
	}
	oldDB.Close()

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		dbConfig.Host, dbConfig.User, dbConfig.Password, dbConfig.Database, dbConfig.Port, dbConfig.Sslmode,
	)

	baseDb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		panic(err)
	}

	return baseDb.Exec("DROP DATABASE test_db").Error
}
