package domain

import (
	"time"
)

type BaseModel struct {
	ID        int64 		`gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
