package services

import (
	"encoding/json"
	"errors"
	"fmt"
	domain "libs/src/internal/domain/models"
	"libs/src/internal/dto"
	"libs/src/internal/repositories"
	api_errors "libs/src/internal/usecase/errors"
	"libs/src/internal/usecase/utils"
	"libs/src/settings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AuthService struct {
	RedisBaseRepository repositories.BaseRedisRepository
	App                 *settings.App
}

func (s *AuthService) GetUserBySession(session string) (dto.UserDTO, error) {
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

func (s *AuthService) setSession(payload string, ttl time.Duration) (string, error) {
	encrypted, err := utils.Encrypt(s.App.Config.AppConfig.SecretKey, string(payload))
	if err != nil {
		return "", err
	}

	newId := uuid.New().String()
	_, err = s.RedisBaseRepository.Create(newId,
		encrypted,
		ttl,
	)
	if err != nil {
		return "", err
	}

	return newId, nil
}

func (s *AuthService) CheckSession(sessionId string) (int64, error) {
	userDto, err := s.GetUserBySession(sessionId)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, api_errors.ErrInvalidToken
		}
		return 0, err
	}

	userRepository := repositories.UserRepository{
		BasePostgresRepository: repositories.BasePostgresRepository[domain.User]{
			Model: domain.User{},
			Db:    s.App.DB,
		},
	}

	user, err := userRepository.GetById(userDto.ID)
	if err != nil {
		return 0, api_errors.ErrInvalidToken
	}

	if user.IsActive {
		return 0, api_errors.ErrInvalidToken
	}

	if user.Email != userDto.Email {
		return 0, api_errors.ErrInvalidToken
	}

	return userDto.ID, nil
}

func (s *AuthService) ConfirmAccount(sessionId string) error {
	userId, err := s.CheckSession(sessionId)
	if err != nil {
		return err
	}

	userRepository := repositories.UserRepository{
		BasePostgresRepository: repositories.BasePostgresRepository[domain.User]{
			Model: domain.User{},
			Db:    s.App.DB,
		},
	}

	user, err := userRepository.GetById(userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return api_errors.ErrInvalidToken
		}
		return err
	}

	user.IsActive = true
	user.Role = domain.USER

	err = s.App.DB.Save(&user).Error
	if err != nil {
		return err
	}

	_, err = s.RedisBaseRepository.Delete(sessionId)
	if err != nil {
		return err
	}

	return nil

}

func (s *AuthService) RegisterUser(data dto.RegisterRequest) error {
	fmt.Println(data)
	if data.Password != data.ConfirmPassword {
		return api_errors.ErrPasswordsDontMatch
	}

	user := domain.User{
		Username: data.Username,
		Email:    data.Email,
		Password: data.Password,
		IsActive: false,
	}

	err := s.App.DB.Create(&user).Error
	if err != nil {
		return err
	}

	go func() {
		toJson, err := json.Marshal(
			dto.UserDTO{
				ID:        user.ID,
				Username:  user.Username,
				Email:     user.Email,
				IsActive:  user.IsActive,
				Role:      user.Role,
				CreatedAt: user.CreatedAt,
			},
		)
		if err != nil {
			settings.AppVar.Logger.Error(fmt.Sprintf("Error creating session: %v", err))
			return
		}

		sessionId, err := s.setSession(string(toJson), time.Duration(settings.AppVar.Config.AuthConfig.EmailConfirmTTL)*time.Second)
		if err != nil {
			settings.AppVar.Logger.Error(fmt.Sprintf("Error creating session: %v", err))
			return
		}

		url := fmt.Sprintf(
			"Thank you for choosing our service, to confirm your registration, follow the url below\nhttp://%s:%d/accounts/auth/confirm-account/%s",
			settings.AppVar.Config.AppConfig.DomainName,
			settings.AppVar.Config.AppConfig.Port,
			sessionId,
		)

		err = utils.SendMail(settings.AppVar.Mail,
			settings.AppVar.Config.Mail.From,
			user.Email,
			"Online-Chat-Golang || Confirm registration",
			url,
		)

		if err != nil {
			settings.AppVar.Logger.Error(fmt.Sprintf("Error registering user: %v", err))
		}
	}()

	return nil
}
