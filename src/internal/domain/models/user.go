package domain

import (
	"reflect"
	"fmt"
)

type User struct {
	BaseModel
	Username    string `gorm:"unique;size:40;not null;"`
	Email       string `gorm:"unique;size:255;not null;"`
	Password    string `gorm:"not null"`
	Description string `gorm:"size:255;"`
	Role        string `gorm:"size:50;not null;default:'anonymous'"`
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

func NewUser(username, email, password, description, role, image string) *User {
	return &User{
		Username:    username,
		Email:       email,
		Password:    password,
		Description: description,
		Role:        role,
		Image:       image,
	}
}
