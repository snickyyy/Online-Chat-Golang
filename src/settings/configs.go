package settings

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type AppConfig struct {
	SecretKey  string `mapstructure:"secret_key"`
	Host       string `mapstructure:"host"`
	Port       int    `mapstructure:"port"`
	Debug      bool   `mapstructure:"debug"`
	Mode       string `mapstructure:"mode"`
	DomainName string `mapstructure:"domain_name"`
	UploadDir  string `mapstructure:"upload_dir"`
}

type PostgresTimeout struct {
	Small  int `mapstructure:"small"`
	Medium int `mapstructure:"medium"`
	Large  int `mapstructure:"large"`
}

type MongoTimeout struct {
	Small  int `mapstructure:"small"`
	Medium int `mapstructure:"medium"`
	Large  int `mapstructure:"large"`
}

type RedisTimeout struct {
	Small  int `mapstructure:"small"`
	Medium int `mapstructure:"medium"`
	Large  int `mapstructure:"large"`
}

type Timeout struct {
	Postgres PostgresTimeout `mapstructure:"postgres"`
	Mongo    MongoTimeout    `mapstructure:"mongo"`
	Redis    RedisTimeout    `mapstructure:"redis"`
}

type Mail struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	From     string `mapstructure:"from"`
	Port     int    `mapstructure:"port"`
	Server   string `mapstructure:"server"`
}

type RedisPrefixes struct {
	SessionPrefix        string `mapstructure:"session"`
	ConfirmEmail         string `mapstructure:"confirm_email"`
	Message              string `mapstructure:"message"`
	ResetPassword        string `mapstructure:"reset_password"`
	ConfirmResetPassword string `mapstructure:"confirm_reset_password"`
	InOnline             string `mapstructure:"in_online"`
}

type RedisConfig struct {
	Host     string        `mapstructure:"host"`
	Port     int           `mapstructure:"port"`
	Password string        `mapstructure:"password"`
	Prefixes RedisPrefixes `mapstructure:"prefixes"`
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
	AuthSessionTTL       int64 `mapstructure:"session_auth_ttl"`
	EmailConfirmTTL      int64 `mapstructure:"confirm_email_ttl"`
	ResetPasswordTTL     int64 `mapstructure:"reset_password_ttl"`
	IsOnlineTTL          int64 `mapstructure:"is_online_ttl"`
	TimeToChangePassword int64 `mapstructure:"time_to_change_password"`
}

type Pagination struct {
	ChatList        int `mapstructure:"chat_list"`
	GlobalChatList  int `mapstructure:"global_chat_list"`
	MessagesList    int `mapstructure:"messages_list"`
	UsersInChatList int `mapstructure:"users_in_chat_list"`
	SearchUsersList int `mapstructure:"search_users_list"`
}

type WsActions struct {
	Subscribe     string `mapstructure:"subscribe"`
	Unsubscribe   string `mapstructure:"unsubscribe"`
	SendMessage   string `mapstructure:"send_message"`
	EditMessage   string `mapstructure:"edit_message"`
	DeleteMessage string `mapstructure:"delete_message"`
}

type WsConfig struct {
	Actions WsActions `mapstructure:"actions"`
}

type BaseConfig struct {
	AppConfig      AppConfig      `mapstructure:"app"`
	WsConfig       WsConfig       `mapstructure:"ws"`
	Timeout        Timeout        `mapstructure:"context_timeout_ms"`
	Pagination     Pagination     `mapstructure:"pagination"`
	PostgresConfig PostgresConfig `mapstructure:"db"`
	AuthConfig     AuthConfig     `mapstructure:"auth"`
	MongoConfig    MongoConfig    `mapstructure:"mongo"`
	RedisConfig    RedisConfig    `mapstructure:"redis"`
	Mail           Mail           `mapstructure:"mail"`
}

func GetBaseConfig() (*BaseConfig, error) {
	basePath, _ := os.Getwd()
	viper.AutomaticEnv()
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(filepath.Dir(basePath))

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
		"mail.username":  viper.GetString("MAIL_USERNAME"),
		"mail.password":  viper.GetString("MAIL_PASSWORD"),
		"mail.from":      viper.GetString("MAIL_FROM"),
		"mail.port":      viper.GetInt("MAIL_PORT"),
		"mail.server":    viper.GetString("MAIL_SERVER"),
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(filepath.Join(basePath, "settings"))

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

	// UPLOAD_DIR
	cfg.AppConfig.UploadDir = filepath.Join(basePath, "assets")

	return cfg, nil
}
