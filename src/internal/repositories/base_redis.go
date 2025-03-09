package repositories

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type BaseRedisRepository struct {
	Client 	*redis.Client
	Ctx		context.Context
}

func (repo *BaseRedisRepository)GetByKey(key string) (string, error) {
	result, err := repo.Client.Get(repo.Ctx, key).Result()
	if err != nil {
        return "", err
    }
	return result, nil
}

func (repo *BaseRedisRepository) Create(key string, value any, expiration time.Duration) (string, error) {
	result, err := repo.Client.Set(repo.Ctx, key, value, expiration).Result()
    if err != nil {
        return "", err
    }
    return result, nil
}

func (repo *BaseRedisRepository) Delete(key string) (int64, error) {
	res, err := repo.Client.Del(repo.Ctx, key).Result()
    if err != nil {
        return 0,err
    }
    return res, nil
}

func (repo *BaseRedisRepository) CountAll() (int64, error) {
	count, err := repo.Client.DBSize(repo.Ctx).Result()
    return count, err
}

func (repo *BaseRedisRepository) IsExist(key string) (bool, error) {
	res, err := repo.Client.Exists(repo.Ctx, key).Result()
	if err != nil {
		return false, err
	}
    return res > 0, nil
}
