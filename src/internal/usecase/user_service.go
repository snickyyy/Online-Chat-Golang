package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"libs/src/internal/domain/enums"
	domain "libs/src/internal/domain/models"
	"libs/src/internal/dto"
	"libs/src/internal/repositories"
	api_errors "libs/src/internal/usecase/errors"
	"libs/src/internal/usecase/utils"
	"libs/src/settings"
	"time"
)

type UserService struct {
	App            *settings.App
	UserRepository *repositories.UserRepository
	SessionService *SessionService
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
		SessionService: NewSessionService(app),
	}
}

func (s *UserService) CreateSuperUser(username string, email string, password string) error {
	passToHash, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	_, err = s.UserRepository.Create(
		&domain.User{
			Username: username,
			Email:    email,
			Password: passToHash,
			IsActive: true,
			Role:     enums.ADMIN,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) GetUserProfile(username string) (*dto.UserProfile, error) {
	user, err := s.UserRepository.Filter("username = ?", username)
	if err != nil {
		return nil, err
	}

	if len(user) != 1 {
		return nil, api_errors.ErrProfileNotFound
	}

	oneUser := user[0]

	if !oneUser.IsActive || oneUser.Role == enums.ANONYMOUS {
		return nil, api_errors.ErrProfileNotFound
	}

	profile := &dto.UserProfile{
		Username:    oneUser.Username,
		Description: oneUser.Description,
		Role:        enums.RolesToLabels[int(oneUser.Role)],
		Image:       oneUser.Image,
		CreatedAt:   oneUser.CreatedAt,
	}

	return profile, nil
}

func (s *UserService) ChangeUserProfile(data dto.ChangeUserProfileRequest, sessionId string) error {
	user, err := s.SessionService.GetUserByAuthSession(sessionId)
	if err != nil {
		return err
	}

	if user.Role == enums.ANONYMOUS || !user.IsActive {
		return api_errors.ErrNeedLoginForChangeProfile
	}

	filterData := map[string]*string{
		"username":    data.NewUsername,
		"description": data.NewDescription,
		"image":       data.NewImage,
	}

	updateData := make(map[string]any, len(filterData))

	for k, v := range filterData {
		if v != nil {
			updateData[k] = v
		}
	}

	err = s.UserRepository.UpdateById(user.ID, updateData)
	if err != nil {
		if errors.Is(err, repositories.ErrDuplicate) {
			return api_errors.ErrUserAlreadyExists
		}
	}
	return nil
}

func (s *UserService) ResetPassword(request dto.ResetPasswordRequest) (int, error) {
	users, err := s.UserRepository.Filter("email = ? OR username = ?", request.UsernameOrEmail, request.UsernameOrEmail)

	if err != nil {
		return -1, err
	}
	if len(users) != 1 {
		return -1, api_errors.ErrUserNotFound
	}

	user := users[0]
	if !user.IsActive || user.Role == enums.ANONYMOUS {
		return -1, api_errors.ErrUserNotFound
	}
	userDto := user.ToDTO()

	secretCode, err := utils.GenerateSecureCode(1000, 9999)
	if err != nil {
		return -1, err
	}

	resetPasswordDto := dto.ResetPasswordSession{
		UserDTO: userDto,
		Code:    secretCode,
	}
	toJson, _ := json.Marshal(&resetPasswordDto)
	encrypt, err := utils.Encrypt(s.App.Config.AppConfig.SecretKey, string(toJson))
	if err != nil {
		return -1, err
	}

	sessionBody := dto.SessionDTO{
		SessionID: uuid.New().String(),
		Expire:    time.Now().Add(time.Duration(s.App.Config.AuthConfig.ResetPasswordTTL) * time.Second),
		Prefix:    s.App.Config.RedisConfig.Prefixes.ConfirmResetPassword,
		Payload:   encrypt,
	}
	_, err = s.SessionService.SetSession(sessionBody)
	if err != nil {
		s.App.Logger.Error(fmt.Sprintf("Error while set session to redis: %s", err.Error()))
		return -1, err
	}

	go func() {
		msg := fmt.Sprintf(
			"For To confirm password reset, follow the link\nhttp://%s:%d/accounts/profile/reset-password/%s",
			s.App.Config.AppConfig.DomainName,
			s.App.Config.AppConfig.Port,
			sessionBody.SessionID,
		)

		for i := 1; i <= 3; i++ {
			err = utils.SendMail(s.App.Mail,
				s.App.Config.Mail.From,
				user.Email,
				"Online-Chat-Golang || Reset-password",
				msg,
			)

			if err == nil {
				break
			}
			s.App.Logger.Error(fmt.Sprintf("Error registering user: %v || try: %d", err, i))
			time.Sleep(time.Duration(i) * time.Second)
		}
	}()
	return secretCode, nil
}

func (s *UserService) ConfirmResetPassword(token string, request dto.ConfirmResetPasswordRequest) error {
	if request.NewPassword != request.ConfirmNewPassword {
		return api_errors.ErrPasswordsDontMatch
	}

	session, err := s.SessionService.GetSession(s.App.Config.RedisConfig.Prefixes.ConfirmResetPassword, token)
	if err != nil {
		return api_errors.ErrInvalidToken
	}

	var sessionBody dto.ResetPasswordSession
	err = s.SessionService.DecryptAndParsePayload(session, &sessionBody)
	if err != nil {
		return api_errors.ErrInvalidToken
	}

	if sessionBody.Code != request.Code {
		return api_errors.ErrInvalidCode
	}

	passToHash, err := utils.HashPassword(request.NewPassword)
	if err != nil {
		return err
	}

	err = s.UserRepository.UpdateById(sessionBody.UserDTO.ID, map[string]any{"password": passToHash})
	if err != nil {
		return api_errors.ErrInvalidToken
	}

	err = s.SessionService.DeleteSession(s.App.Config.RedisConfig.Prefixes.ConfirmResetPassword, token)
	if err != nil {
		return api_errors.ErrInvalidToken
	}

	return nil
}
