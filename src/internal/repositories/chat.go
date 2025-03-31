package repositories

import (
	"libs/src/internal/domain/enums"
	domain "libs/src/internal/domain/models"
	"libs/src/internal/dto"
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

func (r *ChatRepository) GetMemberInfo(memberId, chatId int64) (dto.MemberInfo, error) {
	var memberInfo dto.MemberInfo
	res := r.Db.Raw(`
				SELECT chats.id AS chat_id,
					   chats.title AS chat_title,
					   user_id AS member_id,
					   member_role,
					   chat_members.created_at AS date_joined,
					   chat_members.updated_at
				FROM chats
				JOIN chat_members ON chats.id = chat_members.chat_id
				WHERE chats.id = ? AND chat_members.user_id = ?;
			`, chatId, memberId).Scan(&memberInfo)

	if res.Error != nil {
		return dto.MemberInfo{}, parsePgError(res.Error)
	}
	if res.RowsAffected == 0 {
		return dto.MemberInfo{}, ErrRecordNotFound
	}
	return memberInfo, nil
}
