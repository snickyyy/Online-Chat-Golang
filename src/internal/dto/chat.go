package dto

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
