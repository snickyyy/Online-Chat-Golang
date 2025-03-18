package dto

import "time"

type ErrorResponse struct {
	Error string `json:"error" binding:"required"`
	Code  int    `json:"-"`
}

type SessionDTO struct {
	SessionID string    `json:"id" binding:"required"`
	Expire    time.Time `json:"exp" binding:"required"`
	Prefix    string    `json:"prefix" binding:"required"`
	Payload   string    `json:"payload" binding:"required"`
}
