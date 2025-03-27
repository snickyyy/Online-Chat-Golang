package domain

import (
	"libs/src/internal/dto"
)

type User struct {
	BaseModel
	Username    string `gorm:"unique;size:40;not null;"`
	Email       string `gorm:"unique;size:255;not null;"`
	Password    string `gorm:"not null"`
	Description string `gorm:"size:255;"`
	IsActive    bool   `gorm:"not null;default:false;"`
	Role        byte   `gorm:"not null;default:0"`
	Image       string

	OwnerChats []Chat       `gorm:"foreignKey:OwnerID;"`
	Chats      []ChatMember `gorm:"foreignKey:UserID;"`
}

func (u *User) ToDTO() dto.UserDTO {
	return dto.UserDTO{
		ID:          u.ID,
		Username:    u.Username,
		Email:       u.Email,
		Description: u.Description,
		IsActive:    u.IsActive,
		Role:        u.Role,
		Image:       u.Image,
	}
}
