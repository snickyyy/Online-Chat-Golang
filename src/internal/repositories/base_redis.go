package repositories

import (
	"libs/src/settings"

	"github.com/redis/go-redis/v9"
)

type BaseRedisRepository struct {
	Client 	*redis.Client
}

func (repo *BaseRedisRepository)GetByID(key string) (string, error) {
	result, err := repo.Client.Get(settings.Context.Ctx, key).Result()
	if err != nil {
        return "", err
    }
	return result, nil
}

func (repo *BaseRedisRepository) Create(key string, value interface{}) (string, error) {
	result, err := repo.Client.Set(settings.Context.Ctx, key, value, 0).Result()
    if err != nil {
        return "", err
    }
    return result, nil
}

func (repo *BaseRedisRepository) Delete(key string) (int64, error) {
	res, err := repo.Client.Del(settings.Context.Ctx, key).Result()
    if err != nil {
        return 0,err
    }
    return res, nil
}

func (repo *BaseRedisRepository) CountAll() (int64, error) {
	count, err := repo.Client.DBSize(settings.Context.Ctx).Result()
    return count, err
}

func (repo *BaseRedisRepository) IsExist(key string) (bool, error) {
	res, err := repo.Client.Exists(settings.Context.Ctx, key).Result()
	if err != nil {
		return false, err
	}
    return res > 0, nil
}
