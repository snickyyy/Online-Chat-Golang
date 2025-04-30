package repositories

import (
	models "libs/src/internal/domain/models"

	"context"
	"libs/src/settings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IBaseMongoRepository[T models.Message] interface {
	Create(Ctx context.Context, obj T) (*mongo.InsertOneResult, error)
	GetOne(Ctx context.Context, filters interface{}) (T, error)
	GetAll(Ctx context.Context, filter interface{}, offset int64, limit int64, sortOption ...bson.D) ([]T, error)
	UpdateById(Ctx context.Context, id string, updateFields bson.M) (*mongo.UpdateResult, error)
	DeleteById(Ctx context.Context, id string) (*mongo.DeleteResult, error)
	Count(Ctx context.Context, filters interface{}) (int64, error)
}

type BaseMongoRepository[T models.Message] struct {
	Db             *mongo.Database
	Schema         T
	CollectionName string
}

func (repo *BaseMongoRepository[T]) Create(Ctx context.Context, obj T) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(Ctx, time.Duration(settings.AppVar.Config.Timeout.Mongo.Large)*time.Millisecond)
	defer cancel()

	con := repo.Db.Collection(repo.CollectionName)
	return con.InsertOne(ctx, obj)
}

func (repo *BaseMongoRepository[T]) GetOne(Ctx context.Context, filters interface{}) (T, error) {
	ctx, cancel := context.WithTimeout(Ctx, time.Duration(settings.AppVar.Config.Timeout.Mongo.Small)*time.Millisecond)
	defer cancel()

	con := repo.Db.Collection(repo.CollectionName)

	var chat T
	err := con.FindOne(ctx, filters).Decode(&chat)

	return chat, err
}

func (repo *BaseMongoRepository[T]) GetAll(Ctx context.Context, filter interface{}, offset int64, limit int64, sortOption ...bson.D) ([]T, error) {
	ctx, cancel := context.WithTimeout(Ctx, time.Duration(settings.AppVar.Config.Timeout.Mongo.Medium)*time.Millisecond)
	defer cancel()

	result := []T{}

	con := repo.Db.Collection(repo.CollectionName)

	options := options.Find().SetSkip(offset).SetLimit(limit)

	if len(sortOption) != 0 {
		options.SetSort(sortOption[0])
	}

	cur, err := con.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	err = cur.All(ctx, &result)
	if err != nil {
		return nil, err
	}
	return result, err
}

func (repo *BaseMongoRepository[T]) UpdateById(Ctx context.Context, id string, updateFields bson.M) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(Ctx, time.Duration(settings.AppVar.Config.Timeout.Mongo.Large)*time.Millisecond)
	defer cancel()

	con := repo.Db.Collection(repo.CollectionName)

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	res, err := con.UpdateByID(ctx, objID, bson.M{"$set": updateFields})
	return res, err
}

func (repo *BaseMongoRepository[T]) DeleteById(Ctx context.Context, id string) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(Ctx, time.Duration(settings.AppVar.Config.Timeout.Mongo.Large)*time.Millisecond)
	defer cancel()

	con := repo.Db.Collection(repo.CollectionName)

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	res, err := con.DeleteOne(ctx, bson.M{"_id": objID})
	return res, err
}

func (repo *BaseMongoRepository[T]) Count(Ctx context.Context, filters interface{}) (int64, error) {
	ctx, cancel := context.WithTimeout(Ctx, time.Duration(settings.AppVar.Config.Timeout.Mongo.Small)*time.Millisecond)
	defer cancel()

	con := repo.Db.Collection(repo.CollectionName)

	count, err := con.CountDocuments(ctx, filters)
	return count, err
}
