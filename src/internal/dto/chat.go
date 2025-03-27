package dto

type ChatDTO struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	OwnerID     int64  `json:"owner_id"`
	Description string `json:"description"`
}

type ChatMemberDTO struct {
	ID         int64 `json:"id"`
	ChatID     int64 `json:"chat_id"`
	UserID     int64 `json:"user_id"`
	MemberRole byte  `json:"member_role"`
}
