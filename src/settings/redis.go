package settings

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

func NewRedisSessionClient(config *BaseConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d",config.RedisConfig.Host, config.RedisConfig.Port),
		Password: config.RedisConfig.Password,
		DB:       config.RedisConfig.DB.SessionDb,
		Protocol: 2,
	})
}
