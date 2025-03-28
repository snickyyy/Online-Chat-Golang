package domain

import "libs/src/internal/dto"

type Chat struct {
	BaseModel
	Title       string `gorm:"unique;size:40;not null;"`
	Description string `gorm:"size:255;"`
	OwnerID     int64  `gorm:"not null;"`

	Owner User `gorm:"foreignKey:OwnerID;references:ID;constraint:OnDelete:CASCADE;"`

	Members []ChatMember `gorm:"foreignKey:ChatID;"`
}

func (c *Chat) ToDTO() dto.ChatDTO {
	return dto.ChatDTO{
		ID:          c.ID,
		Title:       c.Title,
		OwnerID:     c.OwnerID,
		Description: c.Description,
	}
}

type ChatMember struct {
	BaseModel
	ChatID     int64 `gorm:"not null;"`
	UserID     int64 `gorm:"not null;"`
	MemberRole byte  `gorm:"not null;"`

	Chat Chat `gorm:"foreignKey:ChatID;references:ID;constraint:OnDelete:CASCADE;"`
	User User `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE;"`
}

func (cm *ChatMember) ToDTO() dto.ChatMemberDTO {
	return dto.ChatMemberDTO{
		ID:         cm.ID,
		ChatID:     cm.ChatID,
		UserID:     cm.UserID,
		MemberRole: cm.MemberRole,
	}
}
