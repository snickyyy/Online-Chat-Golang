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
	"libs/src/pkg/utils"
	"libs/src/settings"
	"time"

	"github.com/google/uuid"
)

type AuthService struct {
	App            *settings.App
	UserRepository repositories.IUserRepository
	SessionService ISessionService
	EmailService   IEmailService
}

func NewAuthService(app *settings.App) *AuthService {
	return &AuthService{
		UserRepository: repositories.NewUserRepository(app),
		SessionService: NewSessionService(app),
		EmailService:   NewEmailService(app),
		App:            app,
	}
}

func (s *AuthService) setAuthCookie(userDto dto.UserDTO) (string, error) {
	sess_ttl := time.Now().Add(time.Duration(s.App.Config.AuthConfig.AuthSessionTTL) * time.Second)

	encoding, _ := json.Marshal(
		dto.AuthSession{
			UserDTO: userDto,
		},
	)

	encrypt, err := utils.Encrypt(s.App.Config.AppConfig.SecretKey, string(encoding))
	if err != nil {
		return "", err
	}

	session := dto.SessionDTO{
		SessionID: uuid.New().String(),
		Expire:    sess_ttl,
		Prefix:    s.App.Config.RedisConfig.Prefixes.SessionPrefix,
		Payload:   encrypt,
	}

	sessionId, err := s.SessionService.SetSession(session)
	if err != nil {
		return "", err
	}

	return sessionId, nil
}

func (s *AuthService) CheckEmailToken(sessionId string) (dto.UserDTO, error) {
	userDto, err := s.SessionService.GetUserByEmailSession(sessionId)
	if err != nil {
		return dto.UserDTO{}, api_errors.ErrInvalidToken
	}

	if userDto.Role != enums.ANONYMOUS || userDto.IsActive {
		return dto.UserDTO{}, api_errors.ErrInvalidToken
	}
	return userDto, nil
}

func (s *AuthService) ConfirmAccount(caller dto.UserDTO, sessionId string) (string, error) {
	if caller.Role != enums.ANONYMOUS || caller.IsActive {
		return "", api_errors.ErrAlreadyLoggedIn
	}
	userDto, err := s.CheckEmailToken(sessionId)
	if err != nil {
		return "", err
	}

	userRepository := s.UserRepository

	changeFields := map[string]any{
		"IsActive": true,
		"Role":     enums.USER,
	}

	err = userRepository.UpdateById(userDto.ID, changeFields)
	if err != nil {
		return "", err
	}

	go func() {
		err = s.SessionService.DeleteSession(s.App.Config.RedisConfig.Prefixes.ConfirmEmail, sessionId)
		if err != nil {
			s.App.Logger.Error(fmt.Sprintf("Error delete email confirm session: %v", err))
		}
	}()

	userDto.Role = enums.USER
	userDto.IsActive = true
	session, err := s.setAuthCookie(userDto)
	if err != nil {
		return "", err
	}

	return session, nil
}

func (s *AuthService) RegisterUser(caller dto.UserDTO, data dto.RegisterRequest) error {
	if caller.Role != enums.ANONYMOUS || caller.IsActive {
		return api_errors.ErrAlreadyLoggedIn
	}
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
	err = s.UserRepository.Create(&user)

	if err != nil {
		if errors.Is(err, repositories.ErrDuplicate) {
			return api_errors.ErrUserAlreadyExists
		}
		return err
	}

	payloadJson, _ := json.Marshal(dto.EmailSession{UserDTO: user.ToDTO()})
	encrypt, err := utils.Encrypt(s.App.Config.AppConfig.SecretKey, string(payloadJson))
	if err != nil {
		return err
	}

	session := dto.SessionDTO{
		SessionID: uuid.New().String(),
		Expire:    time.Now().Add(time.Duration(s.App.Config.AuthConfig.EmailConfirmTTL) * time.Second),
		Prefix:    s.App.Config.RedisConfig.Prefixes.ConfirmEmail,
		Payload:   encrypt,
	}

	sessionId, err := s.SessionService.SetSession(session)
	if err != nil {
		return err
	}

	go func() {
		s.EmailService.SendRegisterEmail(user.Email, sessionId)
	}()

	return nil
}

func (s *AuthService) Login(caller dto.UserDTO, data dto.LoginRequest) (string, error) {
	if caller.Role != enums.ANONYMOUS || caller.IsActive {
		return "", api_errors.ErrAlreadyLoggedIn
	}
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

	session, err := s.setAuthCookie(user.ToDTO())
	if err != nil {
		return "", err
	}

	return session, nil
}

func (s *AuthService) Logout(sessionId string) {
	go func() {
		err := s.SessionService.DeleteSession(s.App.Config.RedisConfig.Prefixes.SessionPrefix, sessionId)
		if err != nil {
			s.App.Logger.Error(fmt.Sprintf("Error deleting session: %v", err))
		}
	}()
}
