package services

import (
	domain "libs/src/internal/domain/models"
	"libs/src/internal/repositories"
	"libs/src/internal/usecase/utils"
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

func (s *UserService) CreateSuperUser(username string, email string, password string) error {
	passToHash, err := utils.HashPassword(password)
	if err != nil { return err }

	_, err = s.UserRepository.Create(
		&domain.User{
			Username: username,
            Email:    email,
            Password: passToHash,
            IsActive: true,
			Role: domain.ADMIN,
		},
	)
	if err != nil {
        return err
    }
	return nil
}
