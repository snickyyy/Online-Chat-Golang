package services

import (
	"encoding/json"
	"libs/src/internal/dto"
	"libs/src/internal/repositories"
	"libs/src/internal/usecase/utils"
	"libs/src/settings"
	"time"

	"github.com/google/uuid"
)

type AuthSessionService struct {
	RedisBaseRepository repositories.BaseRedisRepository
	App 				settings.App
}

func (s *AuthSessionService) GetUserByAuthSession(session string) (dto.UserDTO, error) {
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

func (s *AuthSessionService) SetSession(user dto.UserDTO) (string, error) {
	toJson, err := json.Marshal(user)
	if err != nil {
		return "", err
	}

	encrypted, err := utils.Encrypt(s.App.Config.AppConfig.SecretKey, string(toJson))
	if err != nil {
		return "", err
	}

	newId := uuid.New().String()
	_, err = s.RedisBaseRepository.Create(newId, encrypted, time.Duration(s.App.Config.AuthConfig.AuthSessionTTL)*time.Second)
	if err != nil {
		return "", err
	}

	return newId, nil
}
