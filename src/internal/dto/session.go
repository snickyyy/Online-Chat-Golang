package dto

import (
	"time"
)

type SessionDTO struct {
	SessionID string    `json:"id" binding:"required"`
	Expire    time.Time `json:"exp" binding:"required"`
	Prefix    string    `json:"prefix" binding:"required"`
	Payload   string    `json:"payload" binding:"required"`
}

type ResetPasswordSession struct {
	UserDTO UserDTO `json:"user"`
	Code    int     `json:"code"`
}

type AuthSession struct {
	UserDTO UserDTO `json:"user_dto" binding:"required"`
}

type EmailSession struct {
	UserDTO UserDTO `json:"user_dto" binding:"required"`
}
