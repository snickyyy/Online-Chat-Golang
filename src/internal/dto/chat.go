package dto

type ChatDTO struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	OwnerID     int64  `json:"owner_id"`
	Description string `json:"description"`
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

type ChangeChatRequest struct {
	NewTitle       *string `json:"new_title" binding:"omitempty,min=1,max=38"`
	NewDescription *string `json:"new_description" binding:"omitempty,min=1,max=254"`
}

type ChatsForUserResponse struct {
	Chats []ChatDTO `json:"chats"`
}
