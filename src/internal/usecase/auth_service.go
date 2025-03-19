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

	var sessionBody dto.SessionDTO

	err = json.Unmarshal([]byte(res), &sessionBody)
	if err != nil {
		return dto.UserDTO{}, api_errors.ErrInvalidSession
	}

	decryptResult, err := utils.Decrypt(s.App.Config.AppConfig.SecretKey, sessionBody.Payload)
	if err != nil {
		return dto.UserDTO{}, api_errors.ErrInvalidSession
	}

	var authSessionBody dto.AuthSession

	err = json.Unmarshal([]byte(decryptResult), &authSessionBody)
	if err != nil {
		return dto.UserDTO{}, err
	}

	return authSessionBody.UserDTO, nil
}

func (s *AuthService) setSession(session dto.SessionDTO) (string, error) {
	encoding, _ := json.Marshal(&session)

	_, err := s.RedisBaseRepository.Create(
		session.Prefix,
		session.SessionID,
		string(encoding),
		session.Expire.Sub(time.Now()),
	)
	if err != nil {
		return "", err
	}

	return session.SessionID, nil
}

func (s *AuthService) setAuthSession(user domain.User) (string, error) {
	sess_ttl := time.Now().Add(time.Duration(s.App.Config.AuthConfig.AuthSessionTTL) * time.Second)

	userDto := user.ToDTO()

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
		Payload:   string(encrypt),
	}

	sessionId, err := s.setSession(session)
	if err != nil {
		return "", err
	}

	return sessionId, nil
}

func (s *AuthService) CheckEmailSession(sessionId string) (int64, error) {
	res, err := s.RedisBaseRepository.GetByKey(s.App.Config.RedisConfig.Prefixes.ConfirmEmail, sessionId)
	if err != nil {
		return 0, api_errors.ErrInvalidSession
	}

	var sessionBody dto.SessionDTO

	err = json.Unmarshal([]byte(res), &sessionBody)
	if err != nil {
		return 0, api_errors.ErrInvalidToken
	}

	decryptResult, err := utils.Decrypt(s.App.Config.AppConfig.SecretKey, sessionBody.Payload)
	if err != nil {
		return 0, api_errors.ErrInvalidToken
	}

	var emailSessionBody dto.EmailSession
	err = json.Unmarshal([]byte(decryptResult), &emailSessionBody)
	if err != nil {
		return 0, api_errors.ErrInvalidToken
	}

	userDto := emailSessionBody.UserDTO

	userRepository := s.UserRepository

	user, err := userRepository.GetById(userDto.ID)
	if err != nil {
		return 0, api_errors.ErrInvalidToken
	}

	if user.Role != enums.ANONYMOUS {
		return 0, api_errors.ErrInvalidToken
	}

	if user.IsActive {
		return 0, api_errors.ErrInvalidToken
	}

	if user.Email != userDto.Email {
		return 0, api_errors.ErrInvalidToken
	}

	return user.ID, nil
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
		UserDTO: userDto,
	}
	go func() {
		toJson, err := json.Marshal(payload)
		if err != nil {
			s.App.Logger.Error(fmt.Sprintf("Error creating session: %v", err))
			return
		}

		encrypt, err := utils.Encrypt(s.App.Config.AppConfig.SecretKey, string(toJson))
		if err != nil {
			s.App.Logger.Error(fmt.Sprintf("Error creating session: %v", err))
			return
		}

		session := dto.SessionDTO{
			SessionID: uuid.New().String(),
			Expire:    sess_ttl,
			Prefix:    s.App.Config.RedisConfig.Prefixes.ConfirmEmail,
			Payload:   string(encrypt),
		}

		sessionId, err := s.setSession(session)
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
