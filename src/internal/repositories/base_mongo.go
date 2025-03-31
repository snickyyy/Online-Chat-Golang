package repositories

import (
	domain_interface "libs/src/internal/domain/interfaces"

	"context"
	"libs/src/settings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BaseMongoRepository[T domain_interface.MongoObjectTypes] struct {
	Db             *mongo.Database
	Schema         T
	CollectionName string
}

func (repo *BaseMongoRepository[T]) Create(obj T) (*mongo.InsertOneResult, error) {
	con := repo.Db.Collection(repo.CollectionName)
	return con.InsertOne(settings.Context.Ctx, obj)
}

func (repo *BaseMongoRepository[T]) GetOne(filters interface{}) (T, error) {
	con := repo.Db.Collection(repo.CollectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	var chat T
	err := con.FindOne(ctx, filters).Decode(&chat)

	return chat, err
}

func (repo *BaseMongoRepository[T]) GetAll(filter interface{}, offset int64, limit int64, sortOption ...bson.D) ([]T, error) {
	result := []T{}

	con := repo.Db.Collection(repo.CollectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

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

func (repo *BaseMongoRepository[T]) UpdateById(id string, updateFields bson.M) (*mongo.UpdateResult, error) {
	con := repo.Db.Collection(repo.CollectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	res, err := con.UpdateByID(ctx, objID, bson.M{"$set": updateFields})
	return res, err
}

func (repo *BaseMongoRepository[T]) DeleteById(id string) (*mongo.DeleteResult, error) {
	con := repo.Db.Collection(repo.CollectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	res, err := con.DeleteOne(ctx, bson.M{"_id": objID})
	return res, err
}

func (repo *BaseMongoRepository[T]) Count(filters interface{}) (int64, error) {
	con := repo.Db.Collection(repo.CollectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	count, err := con.CountDocuments(ctx, filters)
	return count, err
}
