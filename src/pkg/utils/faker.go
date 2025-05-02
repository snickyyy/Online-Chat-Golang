package utils

import "github.com/brianvoe/gofakeit/v7"

type FakeUser struct {
	Username    string
	Email       string
	Password    string
	Description string
	IsActive    bool
	Role        byte
	Image       string
}

func NewFakeUser() *FakeUser {
	return &FakeUser{
		Username:    gofakeit.Username(),
		Email:       gofakeit.Email(),
		Password:    gofakeit.Password(true, true, true, false, false, 10),
		Description: gofakeit.Sentence(10),
		IsActive:    gofakeit.Bool(),
		Role:        byte(gofakeit.Number(0, 2)),
		Image:       gofakeit.Sentence(50),
	}
}
