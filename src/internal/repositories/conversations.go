package repositories

import (
	"libs/src/internal/domain/models"
	"libs/src/settings"

	"go.mongodb.org/mongo-driver/mongo"
)

const chatCollection = "chats"

type ChatRepository struct {
	Db *mongo.Database
}


func (repo *ChatRepository) Create(obj domain.Chat) (*mongo.InsertOneResult, error) {
	con := repo.Db.Collection(chatCollection)
	return con.InsertOne(settings.Context.Ctx, obj)
}
