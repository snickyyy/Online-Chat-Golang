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
	"libs/src/pkg/utils"
	"libs/src/settings"
	"mime/multipart"
	"path/filepath"
	"time"
)

type UserService struct {
	App            *settings.App
	UserRepository repositories.IUserRepository
	SessionService ISessionService
	EmailService   IEmailService
}

func NewUserService(app *settings.App) *UserService {
	return &UserService{
		App:            app,
		UserRepository: repositories.NewUserRepository(app),
		SessionService: NewSessionService(app),
		EmailService:   NewEmailService(app),
	}
}

func (s *UserService) saveUserAvatar(oldImage string, newImage *multipart.FileHeader) (string, error) {
	if oldImage != "" {
		utils.DeleteIfExist(filepath.Join(s.App.Config.AppConfig.UploadDir, oldImage))
	}

	fileName := uuid.New().String() + filepath.Ext(newImage.Filename)
	filePath := filepath.Join(s.App.Config.AppConfig.UploadDir, fileName)
	err := utils.UploadFile(newImage, filePath)
	return fileName, err
}

func (s *UserService) CreateSuperUser(username string, email string, password string) error {
	passToHash, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	err = s.UserRepository.Create(
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

func (s *UserService) ChangeUserProfile(caller dto.UserDTO, data dto.ChangeUserProfileRequest) error {
	if caller.Role == enums.ANONYMOUS || !caller.IsActive {
		return api_errors.ErrNeedLoginForChangeProfile
	}

	filterData := map[string]*string{
		"username":    data.NewUsername,
		"description": data.NewDescription,
	}

	if data.NewImage != nil {
		fileName, err := s.saveUserAvatar(caller.Image, data.NewImage)
		if err != nil {
			return err
		}
		filterData["image"] = &fileName
	}

	updateData := make(map[string]any, len(filterData))

	for k, v := range filterData {
		if v != nil {
			updateData[k] = v
		}
	}

	err := s.UserRepository.UpdateById(caller.ID, updateData)
	if err != nil {
		if errors.Is(err, repositories.ErrDuplicate) {
			return api_errors.ErrUserAlreadyExists
		}
		return err
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
		s.EmailService.SendResetPasswordEmail(user.Email, sessionBody.SessionID)
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

func (s *UserService) ChangePassword(caller dto.UserDTO, request dto.ChangePasswordRequest) error {
	user, err := s.UserRepository.GetById(caller.ID)
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			return api_errors.ErrInvalidSession
		}
		return err
	}

	if !user.IsActive || user.Role == enums.ANONYMOUS {
		return api_errors.ErrUnauthorized
	} else if !utils.CheckPasswordHash(user.Password, request.OldPassword) {
		return api_errors.ErrInvalidPassword
	}

	if request.NewPassword == request.OldPassword {
		return api_errors.ErrSamePassword
	}

	if request.NewPassword != request.ConfirmNewPassword {
		return api_errors.ErrPasswordsDontMatch
	}

	passToHash, err := utils.HashPassword(request.NewPassword)
	if err != nil {
		return err
	}

	err = s.UserRepository.UpdateById(user.ID, map[string]any{"password": passToHash})
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			return api_errors.ErrInvalidSession
		}
	}
	return nil
}
