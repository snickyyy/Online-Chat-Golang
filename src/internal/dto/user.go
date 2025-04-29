package dto

import (
	"mime/multipart"
	"time"
)

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
	IsOnline    bool      `json:"is_online"`
	CreatedAt   time.Time `json:"created_at"`
}

type ChangeUserProfileRequest struct {
	NewUsername    *string               `form:"new_username" binding:"omitempty,username"`
	NewDescription *string               `form:"new_description" binding:"omitempty,max=254"`
	NewImage       *multipart.FileHeader `form:"new_image" binding:"omitempty,image"`
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

type ChangePasswordRequest struct {
	OldPassword        string `json:"old_password" binding:"required,password"`
	NewPassword        string `json:"new_password" binding:"required,password"`
	ConfirmNewPassword string `json:"confirm_new_password" binding:"required,password"`
}
