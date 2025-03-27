package domain

import (
	domain "libs/src/internal/domain/models"
)

type PostgresModelsTypes interface {
	domain.User
}

type MongoObjectTypes interface {
	domain.Chat | domain.Message
}
