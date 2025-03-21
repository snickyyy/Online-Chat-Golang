package repositories

import (
	domain "libs/src/internal/domain/interfaces"
	"maps"
	"slices"

	"errors"
	"reflect"

	"gorm.io/gorm"
)

type BasePostgresRepository[T domain.BaseModelsInterface] struct {
	Model T
	Db    *gorm.DB
}

func (r *BasePostgresRepository[T]) Create(obj *T) (interface{}, error) {
	result := r.Db.Create(obj)
	if result.Error != nil {
		return 0, result.Error
	}
	v := reflect.ValueOf(obj).Elem()

	field := v.FieldByName("ID")

	if !field.IsValid() {
		return 0, errors.New("not fount field ID")
	}

	return field.Interface(), nil
}

func (r *BasePostgresRepository[T]) GetById(id int64) (T, error) {
	var obj T
    result := r.Db.First(&obj, id)
    if result.Error != nil {
        return obj, result.Error
    }
    return obj, nil
}

func (repo *BasePostgresRepository[T]) GetAll() ([]T, error) {
	var result []T
    stmt := repo.Db.Find(&result)
    if stmt.Error != nil {
        return result, stmt.Error
    }
    return result, nil
}

func (repo *BasePostgresRepository[T]) Filter(query string, args... interface{}) ([]T, error) {
	var result []T
    stmt := repo.Db.Where(query, args...).Find(&result)
    if stmt.Error != nil {
        return result, stmt.Error
    }
    return result, nil
}

func (repo *BasePostgresRepository[T]) DeleteById(id int64) error {
	var obj T
    result := repo.Db.Where("id = ?", id).Delete(&obj)
	if result.Error != nil {
        return result.Error
    }
    return nil
}

func (repo *BasePostgresRepository[T]) UpdateById(id int64, updateFields map[string]interface{}) error {
	var obj T

	repo.Db.First(&obj, id)

	result := repo.Db.Model(&obj).Select(slices.Collect(maps.Keys(updateFields))).Updates(updateFields)
    if result.Error != nil {
        return result.Error
    }
    return nil
}

func (repo *BasePostgresRepository[T]) Count(filter string, args...interface{}) (int64, error) {
	var obj T
	var count int64
    result := repo.Db.Model(&obj)
	if filter != "" {
        result = result.Where(filter, args...)
    }
	result.Count(&count)
    if result.Error != nil {
        return 0, result.Error
    }
    return count, nil
}

func (repo *BasePostgresRepository[T]) ExecuteQuery(query string, args... interface{}) error {
    stmt := repo.Db.Exec(query, args...)
    if stmt.Error != nil {
        return stmt.Error
    }
    return nil
}
