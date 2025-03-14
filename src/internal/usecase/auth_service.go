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
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AuthService struct {
	RedisBaseRepository *repositories.BaseRedisRepository
	UserRepository      *repositories.UserRepository
	App                 *settings.App
}

func NewAuthService(app *settings.App) *AuthService {
	return &AuthService{
		RedisBaseRepository: &repositories.BaseRedisRepository{
			Client: settings.AppVar.RedisSess,
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

func (s *AuthService) GetUserBySession(session string) (dto.UserDTO, error) {
	res, err := s.RedisBaseRepository.GetByKey(session)
	if err != nil {
		return dto.UserDTO{}, err
	}
	decryptResult, err := utils.Decrypt(s.App.Config.AppConfig.SecretKey, res)
	if err != nil {
		return dto.UserDTO{}, err
	}

	var user dto.AuthSession

	err = json.Unmarshal([]byte(decryptResult), &user)
	if err != nil {
		return dto.UserDTO{}, err
	}

	return user.UserDTO, nil
}

func (s *AuthService) setSession(payload string, ttl time.Duration) (string, error) {
	encrypted, err := utils.Encrypt(s.App.Config.AppConfig.SecretKey, string(payload))
	if err != nil {
		return "", err
	}

	newId := uuid.New().String()
	_, err = s.RedisBaseRepository.Create(
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
	sess_ttl := time.Now().Add(time.Duration(settings.AppVar.Config.AuthConfig.AuthSessionTTL) * time.Second)
	userDto := dto.UserDTO{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		IsActive:  user.IsActive,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}
	payload := dto.AuthSession{
		UserDTO:   userDto,
		TTL:       sess_ttl,
		CreatedAt: time.Now(),
	}

	toJson, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	sessionId, err := s.setSession(string(toJson), time.Duration(settings.AppVar.Config.AuthConfig.AuthSessionTTL)*time.Second)
	if err != nil {
		return "", err
	}

	return sessionId, nil
}

func (s *AuthService) CheckEmailSession(sessionId string) (int64, error) {
	userDto, err := s.GetUserBySession(sessionId)
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
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", api_errors.ErrInvalidToken
		}
		return "", err
	}

	user.IsActive = true
	user.Role = domain.USER

	go func() {
		err = s.App.DB.Save(&user).Error
		if err != nil {
			settings.AppVar.Logger.Error(fmt.Sprintf("Error save user: %v", err))
		}
	}()

	go func() {
		_, err = s.RedisBaseRepository.Delete(sessionId)
		if err != nil {
			settings.AppVar.Logger.Error(fmt.Sprintf("Error delete email confirm session: %v", err))
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
		Role:     domain.ANONYMOUS,
	}

	err = s.App.DB.Create(&user).Error
	if err != nil {
		var pqErr *pgconn.PgError
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return api_errors.ErrUserAlreadyExists
			}
			return err
		}
	}

	sess_ttl := time.Now().Add(time.Duration(settings.AppVar.Config.AuthConfig.EmailConfirmTTL) * time.Second)
	userDto := dto.UserDTO{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		IsActive:  user.IsActive,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}
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

func (s *AuthService) Login(data dto.LoginRequest) (string, error) {
	userRepository := s.UserRepository

	users, err := userRepository.Filter("username = ? OR email = ?", data.UsernameOrEmail, data.UsernameOrEmail)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			settings.AppVar.Logger.Warn(fmt.Sprintf("User not found: %s", data.UsernameOrEmail))
			return "", err
		}
		settings.AppVar.Logger.Error(fmt.Sprintf("Error getting user in login: %v", err))
		return "", err
	}

	if len(users) != 1 {
		return "", api_errors.ErrInvalidCredentials
	}

	user := users[0]
	if !utils.CheckPasswordHash(user.Password, data.Password) || !user.IsActive || user.Role == domain.ANONYMOUS {
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
		_, err := s.RedisBaseRepository.Delete(sessionId)
		if err != nil {
			settings.AppVar.Logger.Error(fmt.Sprintf("Error deleting session: %v", err))
		}
	}()
}
