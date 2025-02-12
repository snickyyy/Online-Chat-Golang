package repositories

import (
	"context"
	domain "libs/src/internal/domain/models"
	"libs/src/settings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const chatCollection = "chats"

type ChatRepository struct {
	Db *mongo.Database
}

func (repo *ChatRepository) Create(obj domain.Chat) (*mongo.InsertOneResult, error) {
	con := repo.Db.Collection(chatCollection)
	return con.InsertOne(settings.Context.Ctx, obj)
}

func (repo *ChatRepository) GetOne(filters interface{}) (domain.Chat, error) {
	con := repo.Db.Collection(chatCollection)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	var chat domain.Chat
	err := con.FindOne(ctx, filters).Decode(&chat)

	return chat, err
}

func (repo *ChatRepository) GetAll(filter interface{}, offset int64, limit int64) ([]domain.Chat, error) {
	result := []domain.Chat{}

	con := repo.Db.Collection(chatCollection)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	options := options.Find().SetSkip(offset).SetLimit(limit)

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

func (repo *ChatRepository) UpdateById(id string, updateFields bson.M) (*mongo.UpdateResult, error) {
	con := repo.Db.Collection(chatCollection)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	res, err := con.UpdateByID(ctx, objID, bson.M{"$set": updateFields})
	return res, err
}

func (repo *ChatRepository) DeleteById(id string) (*mongo.DeleteResult, error) {
	con := repo.Db.Collection(chatCollection)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	res, err := con.DeleteOne(ctx, bson.M{"_id": objID})
	return res, err
}

func (repo *ChatRepository) Count(filters interface{}) (int64, error) {
	con := repo.Db.Collection(chatCollection)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	count, err := con.CountDocuments(ctx, filters)
	return count, err
}
