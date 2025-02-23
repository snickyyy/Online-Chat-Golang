package services

import (
	"encoding/json"
	"errors"
	domain "libs/src/internal/domain/models"
	"libs/src/internal/dto"
	"libs/src/internal/repositories"
	"libs/src/internal/usecase/utils"
	"libs/src/settings"
	"time"

	"github.com/google/uuid"
)

const ErrPasswordsDontMatch = errors.New("passwords don't match")

type AuthService struct {
	RedisBaseRepository repositories.BaseRedisRepository
	App                 *settings.App
}

func (s *AuthService) GetUserByAuthSession(session string) (dto.UserDTO, error) {
	res, err := s.RedisBaseRepository.GetByKey(session)
	if err != nil {
		return dto.UserDTO{}, err
	}
	decryptResult, err := utils.Decrypt(s.App.Config.AppConfig.SecretKey, res)
	if err != nil {
		return dto.UserDTO{}, err
	}

	var user dto.UserDTO

	err = json.Unmarshal([]byte(decryptResult), &user)
	if err != nil {
		return dto.UserDTO{}, err
	}

	return user, nil
}

func (s *AuthService) setSession(user dto.UserDTO) (string, error) {
	toJson, err := json.Marshal(user)
	if err != nil {
		return "", err
	}

	encrypted, err := utils.Encrypt(s.App.Config.AppConfig.SecretKey, string(toJson))
	if err != nil {
		return "", err
	}

	newId := uuid.New().String()
	_, err = s.RedisBaseRepository.Create(newId,
		encrypted,
		time.Duration(s.App.Config.AuthConfig.AuthSessionTTL)*time.Second,
	)
	if err != nil {
		return "", err
	}

	return newId, nil
}

func (s *AuthService) RegisterUser(data dto.RegisterRequest) error {
	if data.Password != data.ConfirmPassword {
		return ErrPasswordsDontMatch
	}

	user := domain.User{
		Username:   data.Username,
        Email:      data.Email,
        Password:   data.Password,
		IsActive: 	false,
	}

	err := s.App.DB.Create(&user).Error
	if err != nil {
        return err
    }
	return nil
}
