package domain

import (
	domain "libs/src/internal/domain/models"
)

type PostgresModelsTypes interface {
	domain.User | domain.Chat | domain.ChatMember
}

type MongoObjectTypes interface {
	domain.Message
}
