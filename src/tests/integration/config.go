package integration

import (
	"fmt"
	"libs/src/settings"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func GetTestConfig() *settings.BaseConfig {
	basePath, _ := os.Getwd()
	if err := godotenv.Load(filepath.Join(filepath.Dir(filepath.Dir(filepath.Dir(basePath))), ".env.test")); err != nil {
		fmt.Println(err)
		panic(err)
	}

	return &settings.BaseConfig{
		AppConfig: settings.AppConfig{
			SecretKey:  os.Getenv("APP_SECRET_KEY"),
			Host:       "localhost",
			Port:       8000,
			Debug:      true,
			Mode:       "dev",
			DomainName: "127.0.0.1",
		},
		PostgresConfig: settings.PostgresConfig{
			Host:     os.Getenv("DB_HOST"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Database: os.Getenv("DB_DATABASE"),
			Port:     5432,
			Sslmode:  os.Getenv("DB_SSL_MODE"),
		},
		AuthConfig: settings.AuthConfig{
			AuthSessionTTL:  86400,
			EmailConfirmTTL: 3600,
		},
		MongoConfig: settings.MongoConfig{
			Uri: os.Getenv("MONGO_URI"),
		},
		RedisConfig: settings.RedisConfig{
			Host: "redis",
			Port: 6379,
			Prefixes: settings.RedisPrefixes{
				SessionPrefix:        "session:",
				ConfirmEmail:         "confirm_email:",
				Message:              "message:",
				ConfirmResetPassword: "confirm_reset_password:",
			},
		},
		Mail: settings.Mail{},
	}
}
