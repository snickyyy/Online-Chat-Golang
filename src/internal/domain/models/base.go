package domain

import (
	"time"
)

type BaseModel struct {
	ID        int64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type BaseMongo struct {
	Id        string    `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at,omitempty"`
	DeleteAt  time.Time `bson:"delete_at" json:"delete_at,omitempty"`
}
