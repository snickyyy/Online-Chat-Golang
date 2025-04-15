package repositories

import (
	"libs/src/internal/domain/enums"
	domain "libs/src/internal/domain/models"
	"libs/src/settings"
)

//go:generate mockery --name=IChatRepository --dir=. --output=../mocks --with-expecter
type IChatRepository interface {
	IBasePostgresRepository[domain.Chat]
	GetListForUser(userId int64, limit int, offset int) ([]domain.Chat, error)
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

func (r *ChatRepository) GetListForUser(userId int64, limit int, offset int) ([]domain.Chat, error) {
	var chats []domain.Chat
	err := r.Db.Table("chats").
		Select("chats.*").
		Joins("JOIN chat_members ON chat_members.chat_id = chats.id").
		Where("chat_members.user_id = ?", userId).
		Limit(limit).
		Offset(offset).
		Find(&chats).Error

	if err != nil {
		return nil, parsePgError(err)
	}

	return chats, nil
}
