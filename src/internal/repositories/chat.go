package repositories

import (
	"context"
	"libs/src/internal/domain/enums"
	domain "libs/src/internal/domain/models"
	"libs/src/settings"
	"time"
)

//go:generate mockery --name=IChatRepository --dir=. --output=../mocks --with-expecter
type IChatRepository interface {
	IBasePostgresRepository[domain.Chat]
	GetListForUser(userId int64, limit int, offset int) ([]domain.Chat, error)
	SearchForUser(userId int64, name string, limit, offset int) ([]domain.Chat, error)
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
	ctx, cancel := context.WithTimeout(settings.Context.Ctx, time.Duration(settings.AppVar.Config.Timeout.Postgres.Large)*time.Millisecond)
	defer cancel()

	tx := r.Db.WithContext(ctx).Begin()
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
	ctx, cancel := context.WithTimeout(settings.Context.Ctx, time.Duration(settings.AppVar.Config.Timeout.Postgres.Medium)*time.Millisecond)
	defer cancel()

	var chats []domain.Chat
	err := r.Db.WithContext(ctx).Table("chats").
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

func (r *ChatRepository) SearchForUser(userId int64, name string, limit, offset int) ([]domain.Chat, error) {
	ctx, cancel := context.WithTimeout(settings.Context.Ctx, time.Duration(settings.AppVar.Config.Timeout.Postgres.Medium)*time.Millisecond)
	defer cancel()

	var chats []domain.Chat
	err := r.Db.WithContext(ctx).Table("chats").
		Select("chats.*").
		Joins("JOIN chat_members ON chat_members.chat_id = chats.id").
		Where("chat_members.user_id = ? AND chats.title LIKE ?", userId, "%"+name+"%").
		Limit(limit).
		Offset(offset).
		Find(&chats).Error

	if err != nil {
		return nil, parsePgError(err)
	}

	return chats, nil
}
