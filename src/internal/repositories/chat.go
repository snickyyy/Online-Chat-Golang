package repositories

import (
	"libs/src/internal/domain/enums"
	domain "libs/src/internal/domain/models"
	"libs/src/settings"
)

//go:generate mockery --name=IChatRepository --dir=. --output=../mocks --with-expecter
type IChatRepository interface {
	IBasePostgresRepository[domain.Chat]
}

func NewChatRepository(app *settings.App) *ChatRepository {
	return &ChatRepository{
		BasePostgresRepository: BasePostgresRepository[domain.Chat]{
			Model: domain.Chat{},
			Db:    app.DB,
		},
	}
}

type ChatRepository struct {
	BasePostgresRepository[domain.Chat]
}

func (r *ChatRepository) Create(chat *domain.Chat) error {
	tx := r.Db.Begin()
	defer tx.Commit()

	if err := tx.Create(chat).Error; err != nil {
		tx.Rollback()
		return parsePgError(err)
	}

	owner := domain.ChatMember{
		ChatID:     chat.ID,
		UserID:     chat.OwnerID,
		MemberRole: enums.OWNER,
	}

	if err := tx.Create(&owner).Error; err != nil {
		tx.Rollback()
		return parsePgError(err)
	}
	return nil
}
