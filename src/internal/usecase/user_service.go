package services

import (
	domain "libs/src/internal/domain/models"
	"libs/src/internal/dto"
	"libs/src/internal/repositories"
	api_errors "libs/src/internal/usecase/errors"
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

func (s *UserService) GetUserProfile(username string) (*dto.UserProfile, error) {
	user, err := s.UserRepository.Filter("username = ?", username)
	if err != nil { return nil, err }

	if len(user) != 1 { return nil, api_errors.ErrProfileNotFound }

	oneUser := user[0]

	if !oneUser.IsActive || oneUser.Role == domain.ANONYMOUS { return nil, api_errors.ErrProfileNotFound }

	profile := &dto.UserProfile{
		Username: 		oneUser.Username,
		Description: 	oneUser.Description,
		Role: 			domain.RolesToLabels[int(oneUser.Role)],	
		Image: 			oneUser.Image,
		CreatedAt: 		oneUser.CreatedAt,	
	}

	return profile, nil
}
