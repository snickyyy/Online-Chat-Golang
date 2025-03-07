package tests

import (
	"fmt"
	"libs/src/settings"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func GetTestConfig() *settings.BaseConfig {
	basePath, _ := os.Getwd()
	if err := godotenv.Load(filepath.Join(filepath.Dir(filepath.Dir(basePath)), ".env")); err != nil {
		fmt.Println(err)
		panic(err)
	}

	return &settings.BaseConfig{
		AppConfig: settings.AppConfig{
			SecretKey:  "testtesttesttesttesttesttesttest",
			Host:       "localhost",
			Port:       8080,
			Debug:      true,
			Mode:       "dev",
			DomainName: "127.0.0.1",
		},
		PostgresConfig: settings.PostgresConfig{
			Host:     "localhost",
			User:     "admin",
			Password: "admin",
			Database: "test_db",
			Port:     5432,
			Sslmode:  "disable",
		},
		AuthConfig: settings.AuthConfig{
			AuthSessionTTL:  86400,
			EmailConfirmTTL: 3600,
		},
		MongoConfig: settings.MongoConfig{
			Uri: "mongodb://localhost:27017",
		},
		RedisConfig: settings.RedisConfig{
			Host: "localhost",
			Port: 6379,
			DB: settings.RedisDbs{
				SessionDb: 2,
			},
		},
		Mail: settings.Mail{},
	}
}
