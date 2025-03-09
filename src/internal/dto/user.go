package dto

import "time"

type UserDTO struct {
	ID          int64  		`json:"id"`
	Username    string 		`json:"username"`
	Email       string 		`json:"email"`
	Description string 		`json:"description"`
	IsActive    bool   		`json:"is_active"`
	Role        byte 		`json:"role"`
	Image       string 		`json:"image"`
	CreatedAt   time.Time 	`json:"created_at"`
}

type UserProfile struct {
	Username    string 		`json:"username"`
	Description string 		`json:"description"`
	Role		string 		`json:"role"`
	Image       string 		`json:"image"`
	CreatedAt   time.Time 	`json:"created_at"`
}
