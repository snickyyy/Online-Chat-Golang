package settings

import (
	"log"
	"github.com/spf13/viper"
)

type AppConfig struct {
	SecretKey string `mapstructure:"secret_key"`
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	Debug     bool   `mapstructure:"debug"`
	Mode      string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	Sslmode  string `mapstructure:"sslmode"`
}

type AuthConfig struct {
	AccessTokenTTL  int `mapstructure:"access_token_ttl"`
	RefreshTokenTTL int `mapstructure:"refresh_token_ttl"`
}

type BaseConfig struct {
	AppConfig
	DatabaseConfig
	AuthConfig
}

func GetBaseConfig() BaseConfig {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading .env file: %v", err)
	}

	viper.AutomaticEnv()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./src/settings")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config.yaml file: %v", err)
	}

	var cfg BaseConfig
	if err := viper.Unmarshal(&cfg); err != nil {
        log.Fatalf("Error unmarshalling config into struct: %v", err)
    }

	return cfg
}
