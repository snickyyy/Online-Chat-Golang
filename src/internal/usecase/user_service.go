package services

import (
	domain "libs/src/internal/domain/models"
	"libs/src/internal/repositories"
	"libs/src/settings"
)

type UserService struct {
	App 				*settings.App
	UserRepository 		*repositories.UserRepository
}

func NewUserService(app *settings.App) *UserService {
	return &UserService{
		App: app,
		UserRepository: &repositories.UserRepository{
			BasePostgresRepository: repositories.BasePostgresRepository[domain.User]{
				Model: domain.User{},
				Db:    app.DB,
			},
		},
	}
}
