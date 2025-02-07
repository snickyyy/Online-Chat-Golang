package settings

import (
	"fmt"
	
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
	AppConfig     AppConfig     `mapstructure:"app"`
	DatabaseConfig DatabaseConfig `mapstructure:"db"`
	AuthConfig    AuthConfig    `mapstructure:"auth"`
}

func GetBaseConfig() (*BaseConfig, error){
	viper.AutomaticEnv()
	viper.SetConfigFile(".env")
	
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading .env file: %v", err)
	}

	envValues := map[string]interface{} {
		"app.secret_key": viper.GetString("APP_SECRET_KEY"),
		"db.host": viper.GetString("DB_HOST"),
		"db.user": viper.GetString("DB_USER"),
		"db.password": viper.GetString("DB_PASSWORD"),
		"db.sslmode": viper.GetString("DB_SSLMODE"),
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./src/settings")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading .yaml file: %v", err)
	}

	for key, value := range envValues {
		viper.Set(key, value)
	}

	cfg := new(BaseConfig)
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshal config %v", err)
	}
	return cfg, nil
}

