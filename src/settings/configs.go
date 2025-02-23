package settings

import (
	"context"
	"fmt"

	"github.com/spf13/viper"
)

var Context *Ctx

type AppConfig struct {
	SecretKey string `mapstructure:"secret_key"`
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	Debug     bool   `mapstructure:"debug"`
	Mode      string `mapstructure:"mode"`
}

type Mail struct {
	Username 	string 	`mapstructure:"username"`
	Password    string  `mapstructure:"password"`
	From 		string 	`mapstructure:"from"`
	Port 		int 	`mapstructure:"port"`
	Server 		string 	`mapstructure:"server"`
}

type Ctx struct {
	Ctx    context.Context
	Cancel context.CancelFunc
}

type RedisDbs struct {
	SessionDb 	int		`mapstructure:"session"`
}

type RedisConfig struct {
	Host     string 	`mapstructure:"host"`
	Port     int    	`mapstructure:"port"`
	Password string 	`mapstructure:"password"`
	DB       RedisDbs   `mapstructure:"db"`
}

type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	Sslmode  string `mapstructure:"sslmode"`
}

type MongoConfig struct {
	Uri string `mapstructure:"uri"`
}

type AuthConfig struct {
	AuthSessionTTL int64 `mapstructure:"session_auth_ttl"`
}

type BaseConfig struct {
	AppConfig      	AppConfig      	`mapstructure:"app"`
	PostgresConfig 	PostgresConfig 	`mapstructure:"db"`
	AuthConfig     	AuthConfig     	`mapstructure:"auth"`
	MongoConfig    	MongoConfig    	`mapstructure:"mongo"`
	RedisConfig    	RedisConfig	  	`mapstructure:"redis"`
	Mail			Mail			`mapstructure:"mail"`
}

func GetBaseConfig() (*BaseConfig, error) {
	viper.AutomaticEnv()
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath("../")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading .env file: %v", err)
	}

	envValues := map[string]interface{}{
		"app.secret_key": viper.GetString("APP_SECRET_KEY"),
		"db.host":        viper.GetString("DB_HOST"),
		"db.user":        viper.GetString("DB_USER"),
		"db.password":    viper.GetString("DB_PASSWORD"),
		"db.sslmode":     viper.GetString("DB_SSLMODE"),
		"mongo.uri":      viper.GetString("MONGO_URI"),
        "mail.username": viper.GetString("MAIL_USERNAME"),
        "mail.password": viper.GetString("MAIL_PASSWORD"),
        "mail.from": viper.GetString("MAIL_FROM"),
        "mail.port": viper.GetInt("MAIL_PORT"),
		"mail.server": viper.GetString("MAIL_SERVER"),
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./settings")

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

func InitContext() {
	ctx, cancel := context.WithCancel(context.Background())
	Context = &Ctx{
		Ctx:    ctx,
		Cancel: cancel,
	}
}
