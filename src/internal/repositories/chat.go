package repositories

import (
	domain "libs/src/internal/domain/models"
	"libs/src/settings"
)

func NewChatRepository(app *settings.App) *ChatRepository {
	return &ChatRepository{
		BasePostgresRepository: BasePostgresRepository[domain.Chat]{
			Model: domain.Chat{},
			Db:    app.DB,
		},
	}
}

func NewChatMemberRepository(app *settings.App) *ChatMemberRepository {
	return &ChatMemberRepository{
		BasePostgresRepository: BasePostgresRepository[domain.ChatMember]{
			Model: domain.ChatMember{},
			Db:    app.DB,
		},
	}
}

type ChatMemberRepository struct {
	BasePostgresRepository[domain.ChatMember]
}

type ChatRepository struct {
	BasePostgresRepository[domain.Chat]
}
