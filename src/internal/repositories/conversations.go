package repositories

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	domain "libs/src/internal/domain/models"
	"libs/src/settings"
)

type ChatRepository struct {
	BasePostgresRepository[domain.Chat]
}

type MessageRepository struct {
	BaseMongoRepository[domain.Message]
}

func NewMessageRepository(app *settings.App) *MessageRepository {
	return &MessageRepository{
		BaseMongoRepository: BaseMongoRepository[domain.Message]{
			Db:             app.MongoDB,
			Schema:         domain.Message{},
			CollectionName: "messages",
		},
	}
}

func (r *MessageRepository) CreateIndex() error {
	compoundIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "chat_id", Value: 1},
			{Key: "created_at", Value: -1},
		},
		Options: options.Index().SetName("chat_id_created_at_index"),
	}
	_, err := r.Db.Collection(r.CollectionName).Indexes().CreateOne(settings.Context.Ctx, compoundIndex)
	return err
}
