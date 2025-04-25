package repositories

import (
	"context"
	models "libs/src/internal/domain/models"
	"libs/src/settings"
	"maps"
	"slices"
	"time"

	"gorm.io/gorm"
)

type IBasePostgresRepository[T models.User | models.Chat | models.ChatMember] interface {
	Create(Ctx context.Context, obj *T) error
	GetById(Ctx context.Context, id int64) (T, error)
	GetAll(Ctx context.Context) ([]T, error)
	Filter(Ctx context.Context, query string, args ...interface{}) ([]T, error)
	DeleteById(Ctx context.Context, id int64) error
	UpdateById(Ctx context.Context, id int64, updateFields map[string]interface{}) error
	Count(Ctx context.Context, filter string, args ...interface{}) (int64, error)
	ExecuteQuery(Ctx context.Context, query string, args ...interface{}) error
}

type BasePostgresRepository[T models.User | models.Chat | models.ChatMember] struct {
	Model T
	Db    *gorm.DB
}

func (r *BasePostgresRepository[T]) Create(Ctx context.Context, obj *T) error {
	ctx, cancel := context.WithTimeout(Ctx, time.Duration(settings.AppVar.Config.Timeout.Postgres.Large)*time.Millisecond)
	defer cancel()

	result := r.Db.
		WithContext(ctx).
		Create(obj)
	if result.Error != nil {
		return parsePgError(result.Error)
	}

	return nil
}

func (r *BasePostgresRepository[T]) GetById(Ctx context.Context, id int64) (T, error) {
	var obj T

	ctx, cancel := context.WithTimeout(Ctx, time.Duration(settings.AppVar.Config.Timeout.Postgres.Small)*time.Millisecond)
	defer cancel()

	result := r.Db.
		WithContext(ctx).
		First(&obj, id)
	if result.Error != nil {
		return obj, parsePgError(result.Error)
	}
	return obj, nil
}

func (repo *BasePostgresRepository[T]) GetAll(Ctx context.Context) ([]T, error) {
	var result []T

	ctx, cancel := context.WithTimeout(Ctx, time.Duration(settings.AppVar.Config.Timeout.Postgres.Medium)*time.Millisecond)
	defer cancel()

	stmt := repo.Db.
		WithContext(ctx).
		Find(&result)
	if stmt.Error != nil {
		return result, parsePgError(stmt.Error)
	}
	return result, nil
}

func (repo *BasePostgresRepository[T]) Filter(Ctx context.Context, query string, args ...interface{}) ([]T, error) {
	ctx, cancel := context.WithTimeout(Ctx, time.Duration(settings.AppVar.Config.Timeout.Postgres.Medium)*time.Millisecond)
	defer cancel()

	var result []T

	stmt := repo.Db.
		WithContext(ctx).
		Where(query, args...).
		Find(&result)
	if stmt.Error != nil {
		return result, parsePgError(stmt.Error)
	}
	return result, nil
}

func (repo *BasePostgresRepository[T]) DeleteById(Ctx context.Context, id int64) error {
	var obj T

	ctx, cancel := context.WithTimeout(Ctx, time.Duration(settings.AppVar.Config.Timeout.Postgres.Small)*time.Millisecond)
	defer cancel()

	result := repo.Db.
		WithContext(ctx).
		Where("id = ?", id).
		Delete(&obj)
	if result.Error != nil {
		return parsePgError(result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (repo *BasePostgresRepository[T]) UpdateById(Ctx context.Context, id int64, updateFields map[string]interface{}) error {
	var obj T

	ctx, cancel := context.WithTimeout(Ctx, time.Duration(settings.AppVar.Config.Timeout.Postgres.Large)*time.Millisecond)
	defer cancel()

	result := repo.Db.
		WithContext(ctx).
		Model(&obj).
		Where("id = ?", id).Select(slices.Collect(maps.Keys(updateFields))).
		Updates(updateFields)

	if result.Error != nil {
		return parsePgError(result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (repo *BasePostgresRepository[T]) Count(Ctx context.Context, filter string, args ...interface{}) (int64, error) {
	var obj T
	var count int64

	ctx, cancel := context.WithTimeout(Ctx, time.Duration(settings.AppVar.Config.Timeout.Postgres.Medium)*time.Millisecond)
	defer cancel()

	result := repo.Db.
		WithContext(ctx).
		Model(&obj)
	if filter != "" {
		result = result.Where(filter, args...)
	}
	result.Count(&count)
	if result.Error != nil {
		return 0, parsePgError(result.Error)
	}
	return count, nil
}

func (repo *BasePostgresRepository[T]) ExecuteQuery(Ctx context.Context, query string, args ...interface{}) error {
	ctx, cancel := context.WithTimeout(Ctx, time.Duration(settings.AppVar.Config.Timeout.Postgres.Medium)*time.Millisecond)
	defer cancel()

	stmt := repo.Db.
		WithContext(ctx).
		Exec(query, args...)
	if stmt.Error != nil {
		return parsePgError(stmt.Error)
	}
	return nil
}
