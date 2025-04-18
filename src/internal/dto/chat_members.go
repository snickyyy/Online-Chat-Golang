package dto

import "time"

type ChangeMemberRoleRequest struct {
	NewRole string `json:"new_role"`
}

type ChatMemberDTO struct {
	ID         int64 `json:"id"`
	ChatID     int64 `json:"chat_id"`
	UserID     int64 `json:"user_id"`
	MemberRole byte  `json:"member_role"`
}

type MemberInfo struct {
	ChatID     int64     `json:"chat_id" gorm:"column:chat_id"`
	ChatTitle  string    `json:"chat_title" gorm:"column:chat_title"`
	MemberID   int64     `json:"member_id" gorm:"column:member_id"`
	MemberRole byte      `json:"member_role" gorm:"column:member_role"`
	DateJoined time.Time `json:"date_joined" gorm:"column:date_joined"`
	UpdateAt   time.Time `json:"updated_at" gorm:"column:updated_at"`
}

type MemberPreview struct { //TODO: добавить isOnline и сделать систему IsOnline
	Username string    `json:"username"`
	Avatar   string    `json:"avatar"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
}

type MemberListPreview struct {
	Members []MemberPreview `json:"members"`
}
