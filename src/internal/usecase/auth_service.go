package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"libs/src/internal/domain/enums"
	domain "libs/src/internal/domain/models"
	"libs/src/internal/dto"
	"libs/src/internal/repositories"
	usecase_errors "libs/src/internal/usecase/errors"
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
		return dto.UserDTO{}, usecase_errors.BadRequestError{Msg: "Invalid email token"}
	}

	if userDto.Role != enums.ANONYMOUS || userDto.IsActive {
		return dto.UserDTO{}, usecase_errors.BadRequestError{Msg: "Invalid email token"}
	}
	return userDto, nil
}

func (s *AuthService) ConfirmAccount(ctx context.Context, caller dto.UserDTO, sessionId string) (string, error) {
	if caller.Role != enums.ANONYMOUS || caller.IsActive {
		return "", usecase_errors.BadRequestError{Msg: "User is already authenticated"}
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

	err = userRepository.UpdateById(ctx, userDto.ID, changeFields)
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

func (s *AuthService) RegisterUser(ctx context.Context, caller dto.UserDTO, data dto.RegisterRequest) error {
	if caller.Role != enums.ANONYMOUS || caller.IsActive {
		return usecase_errors.BadRequestError{Msg: "User is already authenticated"}
	}
	if data.Password != data.ConfirmPassword {
		return usecase_errors.BadRequestError{Msg: "Passwords don't match"}
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
	err = s.UserRepository.Create(ctx, &user)

	if err != nil {
		if errors.Is(err, repositories.ErrDuplicate) {
			return usecase_errors.AlreadyExistsError{Msg: "User with this email or username already exists"}
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

	go s.EmailService.SendRegisterEmail(user.Email, sessionId)

	return nil
}

func (s *AuthService) Login(ctx context.Context, caller dto.UserDTO, data dto.LoginRequest) (string, error) {
	if caller.Role != enums.ANONYMOUS || caller.IsActive {
		return "", usecase_errors.BadRequestError{Msg: "User is already authorized"}
	}
	userRepository := s.UserRepository

	users, err := userRepository.Filter(ctx, "username = ? OR email = ?", data.UsernameOrEmail, data.UsernameOrEmail)
	if err != nil {
		return "", usecase_errors.BadRequestError{Msg: "Invalid credentials"}
	}

	if len(users) != 1 {
		return "", usecase_errors.BadRequestError{Msg: "Invalid credentials"}
	}

	user := users[0]
	if !utils.CheckPasswordHash(user.Password, data.Password) || !user.IsActive || user.Role == enums.ANONYMOUS {
		return "", usecase_errors.BadRequestError{Msg: "Invalid credentials"}
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
