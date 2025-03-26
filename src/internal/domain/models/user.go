package domain

import (
	"fmt"
	"libs/src/internal/dto"
	"reflect"
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
}

func (u User) String() string {
	result := ""

	type_ := reflect.TypeOf(u)
	value := reflect.ValueOf(u)

	for i := 0; i < type_.NumField(); i++ {
		result += "[" + type_.Field(i).Name + ": " + fmt.Sprintf("%v ", value.Field(i)) + "]  ||  "
	}

	return result
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
