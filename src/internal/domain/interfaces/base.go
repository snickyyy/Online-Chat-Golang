package domain

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BaseModelsInterface interface {
	String() string
}

type BaseMongoRepository interface {
	UpdateById(id string, updateFields bson.M) (*mongo.UpdateResult, error)
	DeleteById(id string) (*mongo.DeleteResult, error)
	Count(filters interface{}) (int64, error)
}

type BaseMongoInterface interface {
	Mapping() 	map[string]interface{}
}
