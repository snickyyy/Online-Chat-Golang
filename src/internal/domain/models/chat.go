package domain

type Chat struct {
	BaseModel
	Title       string `gorm:"unique;size:40;not null;"`
	Description string `gorm:"size:255;"`
	OwnerID     int64  `gorm:"not null;"`

	Owner User `gorm:"foreignKey:OwnerID;references:ID;constraint:OnDelete:CASCADE;"`

	Members []ChatMember `gorm:"foreignKey:ChatID;"`
}

type ChatMember struct {
	BaseModel
	ChatID     int64 `gorm:"not null;"`
	UserID     int64 `gorm:"not null;"`
	MemberRole byte  `gorm:"not null;"`

	Chat Chat `gorm:"foreignKey:ChatID;references:ID;constraint:OnDelete:CASCADE;"`
	User User `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE;"`
}
