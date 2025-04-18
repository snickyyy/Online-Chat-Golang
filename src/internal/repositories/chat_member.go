package repositories

import (
	"fmt"
	"libs/src/internal/domain/enums"
	domain "libs/src/internal/domain/models"
	"libs/src/internal/dto"
	"libs/src/settings"
	"time"
)

//go:generate mockery --name=IChatMemberRepository --dir=. --output=../mocks --with-expecter
type IChatMemberRepository interface {
	IBasePostgresRepository[domain.ChatMember]
	SetNewRole(chatId, userId int64, role byte) error
	GetMemberInfo(memberId, chatId int64) (dto.MemberInfo, error)
	DeleteMember(memberId, chatId int64) error
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

func (r *ChatMemberRepository) GetMemberInfo(memberId, chatId int64) (dto.MemberInfo, error) {
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

func (r *ChatMemberRepository) DeleteMember(memberId, chatId int64) error {
	res := r.Db.Where("user_id = ? AND chat_id = ?", memberId, chatId).Delete(&r.Model)
	if res.Error != nil {
		return parsePgError(res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (r *ChatMemberRepository) GetMembersPreview(chatId int64) ([]dto.MemberPreview, error) {
	members := []struct {
		Username string    `gorm:"column:username"`
		Avatar   string    `gorm:"column:avatar"`
		JoinedAt time.Time `gorm:"column:joined_at"`
		Role     string    `gorm:"column:role"`
	}{}

	buildRoleCase := ""
	for k, v := range enums.ChatRolesToLabels {
		buildRoleCase += fmt.Sprintf("WHEN member_role = %d THEN '%s'\n", k, v)
	}

	res := r.Db.Raw(fmt.Sprintf(
		`SELECT
	   	users.username AS username,
	   	users.image AS avatar,
	   	chat_members.created_at AS joined_at,
	   	CASE
	            %s
	   	END AS role
		FROM chat_members
		JOIN users ON chat_members.user_id = users.id
		WHERE chat_members.chat_id = ?`, buildRoleCase), chatId).Scan(&members)
	if res.Error != nil {
		return nil, parsePgError(res.Error)
	}

	result := make([]dto.MemberPreview, len(members))
	for i, member := range members {
		result[i] = dto.MemberPreview{
			Username: member.Username,
			Avatar:   member.Avatar,
			JoinedAt: member.JoinedAt,
			Role:     member.Role,
		}
	}

	return result, nil
}
