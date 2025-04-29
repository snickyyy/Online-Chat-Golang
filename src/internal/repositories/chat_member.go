package repositories

import (
	"context"
	"fmt"
	"libs/src/internal/domain/enums"
	domain "libs/src/internal/domain/models"
	"libs/src/internal/dto"
	"libs/src/settings"
	"strconv"
	"time"
)

//go:generate mockery --name=IChatMemberRepository --dir=. --output=../mocks --with-expecter
type IChatMemberRepository interface {
	IBasePostgresRepository[domain.ChatMember]
	SetNewRole(Ctx context.Context, chatId, userId int64, role byte) error
	GetMemberInfo(Ctx context.Context, memberId, chatId int64) (dto.MemberInfo, error)
	DeleteMember(Ctx context.Context, memberId, chatId int64) error
	GetMembersPreview(Ctx context.Context, chatId int64, limit, offset int, searchUsername string) ([]dto.MemberPreview, error)
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

func (r *ChatMemberRepository) SetNewRole(Ctx context.Context, chatId, userId int64, role byte) error {
	ctx, cancel := context.WithTimeout(Ctx, time.Duration(settings.AppVar.Config.Timeout.Postgres.Small)*time.Millisecond)
	defer cancel()

	res := r.Db.WithContext(ctx).Model(&r.Model).
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

func (r *ChatMemberRepository) GetMemberInfo(Ctx context.Context, memberId, chatId int64) (dto.MemberInfo, error) {
	ctx, cancel := context.WithTimeout(Ctx, time.Duration(settings.AppVar.Config.Timeout.Postgres.Medium)*time.Millisecond)
	defer cancel()

	var memberInfo dto.MemberInfo
	res := r.Db.WithContext(ctx).Raw(`
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

func (r *ChatMemberRepository) DeleteMember(Ctx context.Context, memberId, chatId int64) error {
	ctx, cancel := context.WithTimeout(Ctx, time.Duration(settings.AppVar.Config.Timeout.Postgres.Small)*time.Millisecond)
	defer cancel()

	res := r.Db.WithContext(ctx).Where("user_id = ? AND chat_id = ?", memberId, chatId).Delete(&r.Model)
	if res.Error != nil {
		return parsePgError(res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (r *ChatMemberRepository) GetMembersPreview(Ctx context.Context, chatId int64, limit, offset int, searchUsername string) ([]dto.MemberPreview, error) {
	members := []struct {
		Id       int64     `gorm:"column:user_id"`
		Username string    `gorm:"column:username"`
		Avatar   string    `gorm:"column:avatar"`
		JoinedAt time.Time `gorm:"column:joined_at"`
		Role     string    `gorm:"column:role"`
	}{}

	ctx, cancel := context.WithTimeout(Ctx, time.Duration(settings.AppVar.Config.Timeout.Postgres.Medium)*time.Millisecond)
	defer cancel()

	query := r.Db.WithContext(ctx).Table("chat_members").
		Select(fmt.Sprintf(`
		users.id AS user_id,
		users.username AS username,
		users.image AS avatar,
		chat_members.created_at AS joined_at,
		CASE
			%s
		END AS role
		`, r.buildCaseByRole(enums.ChatRolesToLabels))).
		Joins("JOIN users ON chat_members.user_id = users.id").
		Where("chat_members.chat_id = ?", chatId)

	if searchUsername != "" {
		query.Where("users.username LIKE ?", "%"+searchUsername+"%")
	}

	res := query.Limit(limit).Offset(offset).Scan(&members)

	if res.Error != nil {
		return nil, parsePgError(res.Error)
	}

	ids := make([]string, len(members))

	for i, member := range members {
		ids[i] = settings.AppVar.Config.RedisConfig.Prefixes.InOnline + strconv.Itoa(int(member.Id))
	}

	redisRepository := NewBaseRedisRepository(settings.AppVar)
	usersOnline, err := redisRepository.ManyToGet(ctx, ids)
	if err != nil {
		return nil, err
	}

	result := make([]dto.MemberPreview, len(members))
	for i, member := range members {
		result[i] = dto.MemberPreview{
			Username: member.Username,
			Avatar:   member.Avatar,
			IsOnline: usersOnline[i] != nil,
			JoinedAt: member.JoinedAt,
			Role:     member.Role,
		}
	}

	return result, nil
}
