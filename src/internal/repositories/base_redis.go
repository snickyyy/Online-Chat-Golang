package repositories

import (
	"context"
	"encoding/json"
	"libs/src/internal/dto"
	"libs/src/settings"
	"time"

	"github.com/redis/go-redis/v9"
)

//go:generate mockery --name=IBaseRedisRepository --dir=. --output=../mocks --with-expecter
type IBaseRedisRepository interface {
	SetDTO(Ctx context.Context, prefix string, obj dto.SessionDTO) (string, error)
	GetByKey(Ctx context.Context, prefix string, key string) (string, error)
	Create(Ctx context.Context, prefix string, key string, value any, expiration time.Duration) (string, error)
	Delete(Ctx context.Context, prefix string, key string) (int64, error)
	CountAll(Ctx context.Context) (int64, error)
	IsExist(Ctx context.Context, prefix string, key string) (bool, error)
}

type BaseRedisRepository struct {
	Client *redis.Client
}

func NewBaseRedisRepository(app *settings.App) *BaseRedisRepository {
	return &BaseRedisRepository{
		Client: app.RedisClient,
	}
}

func (repo *BaseRedisRepository) SetDTO(Ctx context.Context, prefix string, obj dto.SessionDTO) (string, error) {
	ctx, cancel := context.WithTimeout(Ctx, time.Duration(settings.AppVar.Config.Timeout.Redis.Large)*time.Millisecond)
	defer cancel()

	toJson, _ := json.Marshal(obj)
	obj.SessionID = prefix + obj.SessionID
	result, err := repo.Client.Set(
		ctx,
		obj.SessionID,
		toJson,
		time.Until(obj.Expire)).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

func (repo *BaseRedisRepository) GetByKey(Ctx context.Context, prefix string, key string) (string, error) {
	ctx, cancel := context.WithTimeout(Ctx, time.Duration(settings.AppVar.Config.Timeout.Redis.Small)*time.Millisecond)
	defer cancel()

	result, err := repo.Client.Get(ctx, prefix+key).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

func (repo *BaseRedisRepository) Create(Ctx context.Context, prefix string, key string, value any, expiration time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(Ctx, time.Duration(settings.AppVar.Config.Timeout.Redis.Large)*time.Millisecond)
	defer cancel()

	result, err := repo.Client.Set(ctx, prefix+key, value, expiration).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

func (repo *BaseRedisRepository) Delete(Ctx context.Context, prefix string, key string) (int64, error) {
	ctx, cancel := context.WithTimeout(Ctx, time.Duration(settings.AppVar.Config.Timeout.Redis.Small)*time.Millisecond)
	defer cancel()

	res, err := repo.Client.Del(ctx, prefix+key).Result()
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (repo *BaseRedisRepository) CountAll(Ctx context.Context) (int64, error) {
	ctx, cancel := context.WithTimeout(Ctx, time.Duration(settings.AppVar.Config.Timeout.Redis.Medium)*time.Millisecond)
	defer cancel()

	count, err := repo.Client.DBSize(ctx).Result()
	return count, err
}

func (repo *BaseRedisRepository) IsExist(Ctx context.Context, prefix string, key string) (bool, error) {
	ctx, cancel := context.WithTimeout(Ctx, time.Duration(settings.AppVar.Config.Timeout.Redis.Small)*time.Millisecond)
	defer cancel()

	res, err := repo.Client.Exists(ctx, prefix+key).Result()
	if err != nil {
		return false, err
	}
	return res > 0, nil
}
