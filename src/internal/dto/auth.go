package dto

import "time"

type RegisterRequest struct {
	Username        string `json:"username" binding:"required,min=4,max=28"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

type RegisterResponse struct {
	Message string `json:"message" binding:"required"`
	Status  bool   `json:"status" binding:"required"`
}

type AuthSession struct {
	UserDTO 	UserDTO 	`json:"user_dto" binding:"required"`
	TTL     	time.Time 	`json:"ttl" binding:"required"`
	CreatedAt 	time.Time 	`json:"created_at" binding:"required"`
}
