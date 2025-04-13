package dto

import "time"

type ChatDTO struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	OwnerID     int64  `json:"owner_id"`
	Description string `json:"description"`
}

type ChatMemberDTO struct {
	ID         int64 `json:"id"`
	ChatID     int64 `json:"chat_id"`
	UserID     int64 `json:"user_id"`
	MemberRole byte  `json:"member_role"`
}

type CreateChatRequest struct {
	Title       string `json:"title" binding:"required,min=1,max=38"`
	Description string `json:"description" binding:"required,min=1,max=254"`
}

type ChatPreview struct {
	Title       string `json:"title"`
	Owner       string `json:"owner"`
	Description string `json:"description"`
}

type FilterChatsResponse struct {
	Chats []ChatPreview `json:"chats"`
}

type MemberInfo struct {
	ChatID     int64     `json:"chat_id" gorm:"column:chat_id"`
	ChatTitle  string    `json:"chat_title" gorm:"column:chat_title"`
	MemberID   int64     `json:"member_id" gorm:"column:member_id"`
	MemberRole byte      `json:"member_role" gorm:"column:member_role"`
	DateJoined time.Time `json:"date_joined" gorm:"column:date_joined"`
	UpdateAt   time.Time `json:"updated_at" gorm:"column:updated_at"`
}

type ChangeChatRequest struct {
	NewTitle       *string `json:"new_title" binding:"omitempty,min=1,max=38"`
	NewDescription *string `json:"new_description" binding:"omitempty,min=1,max=254"`
}
