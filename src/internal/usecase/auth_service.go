package services

import (
	"encoding/json"
	"errors"
	"fmt"
	domain "libs/src/internal/domain/models"
	"libs/src/internal/dto"
	"libs/src/internal/repositories"
	"libs/src/internal/usecase/utils"
	"libs/src/settings"
	"time"

	"github.com/google/uuid"
)

var ErrPasswordsDontMatch = errors.New("passwords don't match")

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

func (s *AuthService) RegisterUser(data dto.RegisterRequest) error {
	fmt.Println(data)
	if data.Password != data.ConfirmPassword {
		return ErrPasswordsDontMatch
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
				ID:        	user.ID,
				Username:  	user.Username,
				Email:     	user.Email,
				IsActive:  	user.IsActive,
				Role: 		user.Role,
				CreatedAt: 	user.CreatedAt,
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
