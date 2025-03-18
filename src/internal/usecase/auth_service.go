package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"libs/src/internal/domain/enums"
	domain "libs/src/internal/domain/models"
	"libs/src/internal/dto"
	"libs/src/internal/repositories"
	api_errors "libs/src/internal/usecase/errors"
	"libs/src/internal/usecase/utils"
	"libs/src/settings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type AuthService struct {
	RedisBaseRepository *repositories.BaseRedisRepository
	UserRepository      *repositories.UserRepository
	App                 *settings.App
}

func NewAuthService(app *settings.App) *AuthService {
	return &AuthService{
		RedisBaseRepository: &repositories.BaseRedisRepository{
			Client: app.RedisClient,
			Ctx:    settings.Context.Ctx,
		},
		UserRepository: &repositories.UserRepository{
			BasePostgresRepository: repositories.BasePostgresRepository[domain.User]{
				Model: domain.User{},
				Db:    app.DB,
			},
		},
		App: app,
	}
}

func (s *AuthService) GetUserBySession(prefix string, session string) (dto.UserDTO, error) {
	res, err := s.RedisBaseRepository.GetByKey(prefix, session)
	if err != nil {
		return dto.UserDTO{}, api_errors.ErrInvalidSession
	}
	decryptResult, err := utils.Decrypt(s.App.Config.AppConfig.SecretKey, res)
	if err != nil {
		return dto.UserDTO{}, api_errors.ErrInvalidSession
	}

	var user dto.AuthSession

	err = json.Unmarshal([]byte(decryptResult), &user)
	if err != nil {
		return dto.UserDTO{}, err
	}

	return user.UserDTO, nil
}

func (s *AuthService) setSession(prefix string, payload string, ttl time.Duration) (string, error) {
	encrypted, err := utils.Encrypt(s.App.Config.AppConfig.SecretKey, string(payload))
	if err != nil {
		return "", err
	}

	newId := uuid.New().String()
	_, err = s.RedisBaseRepository.Create(
		prefix,
		newId,
		encrypted,
		ttl,
	)
	if err != nil {
		return "", err
	}

	return newId, nil
}

func (s *AuthService) setAuthSession(user domain.User) (string, error) {
	sess_ttl := time.Now().Add(time.Duration(s.App.Config.AuthConfig.AuthSessionTTL) * time.Second)
	userDto := user.ToDTO()
	payload := dto.AuthSession{
		UserDTO:   userDto,
		TTL:       sess_ttl,
		CreatedAt: time.Now(),
	}

	toJson, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	sessionId, err := s.setSession(
		s.App.Config.RedisConfig.Prefixes.SessionPrefix,
		string(toJson),
		time.Duration(s.App.Config.AuthConfig.AuthSessionTTL)*time.Second,
	)
	if err != nil {
		return "", err
	}

	return sessionId, nil
}

func (s *AuthService) CheckEmailSession(sessionId string) (int64, error) {
	userDto, err := s.GetUserBySession(s.App.Config.RedisConfig.Prefixes.ConfirmEmail, sessionId)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, api_errors.ErrInvalidToken
		}
		return 0, err
	}

	userRepository := s.UserRepository

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

func (s *AuthService) ConfirmAccount(sessionId string) (string, error) {
	userId, err := s.CheckEmailSession(sessionId)
	if err != nil {
		return "", err
	}

	userRepository := s.UserRepository

	user, err := userRepository.GetById(userId)
	// TODO: Оптимизировать тут что бы был только 1 запрос(юзер возвращался с CheckEmailSession
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			return "", api_errors.ErrInvalidToken
		}
		return "", err
	}

	changeFields := map[string]any{
		"IsActive": true,
		"Role":     enums.USER,
	}

	go func() {
		err = userRepository.UpdateById(user.ID, changeFields)
		if err != nil {
			s.App.Logger.Error(fmt.Sprintf("Error save user: %v", err))
		}
	}()

	go func() {
		_, err = s.RedisBaseRepository.Delete(s.App.Config.RedisConfig.Prefixes.ConfirmEmail, sessionId)
		if err != nil {
			s.App.Logger.Error(fmt.Sprintf("Error delete email confirm session: %v", err))
		}
	}()

	session, err := s.setAuthSession(user)
	if err != nil {
		return "", err
	}

	return session, nil
}

func (s *AuthService) RegisterUser(data dto.RegisterRequest) error {
	if data.Password != data.ConfirmPassword {
		return api_errors.ErrPasswordsDontMatch
	}

	hashedPassword, err := utils.HashPassword(data.Password)
	if err != nil {
		return err
	}

	user := domain.User{
		Username: data.Username,
		Email:    data.Email,
		Password: hashedPassword,
		IsActive: false,
		Role:     enums.ANONYMOUS,
	}

	_, err = s.UserRepository.Create(&user)
	if err != nil {
		if errors.Is(err, repositories.ErrDuplicate) {
			return api_errors.ErrUserAlreadyExists
		}
	}

	sess_ttl := time.Now().Add(time.Duration(s.App.Config.AuthConfig.EmailConfirmTTL) * time.Second)
	userDto := user.ToDTO()
	payload := dto.EmailSession{
		AuthSession: dto.AuthSession{
			UserDTO:   userDto,
			TTL:       sess_ttl,
			CreatedAt: time.Now(),
		},
	}
	go func() {
		toJson, err := json.Marshal(payload)
		if err != nil {
			s.App.Logger.Error(fmt.Sprintf("Error creating session: %v", err))
			return
		}

		sessionId, err := s.setSession(s.App.Config.RedisConfig.Prefixes.ConfirmEmail, string(toJson), time.Duration(s.App.Config.AuthConfig.EmailConfirmTTL)*time.Second)
		if err != nil {
			s.App.Logger.Error(fmt.Sprintf("Error creating session: %v", err))
			return
		}

		msg := fmt.Sprintf(
			"Thank you for choosing our service, to confirm your registration, follow the url below\nhttp://%s:%d/accounts/auth/confirm-account/%s",
			s.App.Config.AppConfig.DomainName,
			s.App.Config.AppConfig.Port,
			sessionId,
		)

		err = utils.SendMail(s.App.Mail,
			s.App.Config.Mail.From,
			user.Email,
			"Online-Chat-Golang || Confirm registration",
			msg,
		)

		if err != nil {
			s.App.Logger.Error(fmt.Sprintf("Error registering user: %v", err))
		}
	}()

	return nil
}

func (s *AuthService) Login(data dto.LoginRequest) (string, error) {
	userRepository := s.UserRepository

	users, err := userRepository.Filter("username = ? OR email = ?", data.UsernameOrEmail, data.UsernameOrEmail)
	if err != nil {
		s.App.Logger.Error(fmt.Sprintf("Error getting user in login: %v", err))
		return "", api_errors.ErrInvalidCredentials
	}

	if len(users) != 1 {
		return "", api_errors.ErrInvalidCredentials
	}

	user := users[0]
	if !utils.CheckPasswordHash(user.Password, data.Password) || !user.IsActive || user.Role == enums.ANONYMOUS {
		return "", api_errors.ErrInvalidCredentials
	}

	session, err := s.setAuthSession(user)
	if err != nil {
		return "", err
	}

	return session, nil
}

func (s *AuthService) Logout(sessionId string) {
	go func() {
		_, err := s.RedisBaseRepository.Delete(s.App.Config.RedisConfig.Prefixes.SessionPrefix, sessionId)
		if err != nil {
			s.App.Logger.Error(fmt.Sprintf("Error deleting session: %v", err))
		}
	}()
}
