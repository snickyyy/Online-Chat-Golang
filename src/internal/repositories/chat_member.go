package repositories

import (
	domain "libs/src/internal/domain/models"
	"libs/src/settings"
)

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

func (r *ChatMemberRepository) SetNewRole(chatId, userId int64, role byte) error {
	res := r.Db.Model(&r.Model).
		Where("chat_id = ? AND user_id = ?", chatId, userId).
		Update("member_role", role)

	if res.Error != nil {
		return parsePgError(res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
