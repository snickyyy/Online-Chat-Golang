package repositories

import domain "libs/src/internal/domain/models"

type ChatRepository struct {
	BaseMongoRepository[domain.Chat]
}

type MessageRepository struct {
	BaseMongoRepository[domain.Message]
}
