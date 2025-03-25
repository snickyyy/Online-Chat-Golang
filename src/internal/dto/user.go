package dto

import "time"

type UserDTO struct {
	ID          int64     `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	Role        byte      `json:"role"`
	Image       string    `json:"image"`
	CreatedAt   time.Time `json:"created_at"`
}

type UserProfile struct {
	Username    string    `json:"username"`
	Description string    `json:"description"`
	Role        string    `json:"role"`
	Image       string    `json:"image"`
	CreatedAt   time.Time `json:"created_at"`
}

type ChangeUserProfileRequest struct {
	NewUsername    *string `json:"new_username" binding:"username"`
	NewDescription *string `json:"new_description" binding:"max=254"`
	NewImage       *string `json:"new_image"` // TODO: сделать типа что бы изображение загружалось а не путь к нему
}

type ChangeUserProfileResponse struct {
	ChangedFields ChangeUserProfileRequest `json:"changed_fields"`
	Message       string                   `json:"message"`
}

type ResetPasswordRequest struct {
	UsernameOrEmail string `json:"username_or_email" binding:"required"`
}

type ResetPasswordResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type ConfirmResetPasswordRequest struct {
	NewPassword        string `json:"new_password" binding:"required,password"`
	ConfirmNewPassword string `json:"confirm_new_password" binding:"required,password"`
	Code               int    `json:"code" binding:"required,numeric"`
}
