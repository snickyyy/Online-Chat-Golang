package services

import (
	domain "libs/src/internal/domain/models"
	"libs/src/internal/dto"
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

func (s *UserService) CreateSuperUser(username string, email string, password string) (dto.UserDTO, error) {
	user, err := s.UserRepository.Create(
		&domain.User{
			Username: username,
            Email:    email,
            Password: password,
            IsActive: true,
			Role: domain.ADMIN,
		},
	)
	if err != nil {
        return dto.UserDTO{}, err
    }
	return user.(*domain.User).ToDTO(), nil
}
