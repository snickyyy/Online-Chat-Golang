package settings

import "github.com/redis/go-redis/v9"

func NewRedisSessionClient(config *BaseConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     config.RedisConfig.Host,
		Password: config.RedisConfig.Password,
		DB:       config.RedisConfig.DB.SessionDb,
		Protocol: 2,
	})
}
