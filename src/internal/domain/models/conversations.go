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

type Message struct {
	BaseMongo
	SenderId  int64  `bson:"sender_id" json:"sender_id"`
	ChatId    string `bson:"chat_id" json:"chat_id"`
	Content   string `bson:"content" json:"content"`
	IsRead    bool   `bson:"is_read" json:"is_read"`
	IsUpdated bool   `bson:"is_updated" json:"is_updated"`
	IsDeleted bool   `bson:"is_deleted" json:"is_deleted"`
}
