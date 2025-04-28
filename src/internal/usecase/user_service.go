package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"libs/src/internal/domain/enums"
	domain "libs/src/internal/domain/models"
	"libs/src/internal/dto"
	"libs/src/internal/repositories"
	usecase_errors "libs/src/internal/usecase/errors"
	"libs/src/pkg/utils"
	"libs/src/settings"
	"mime/multipart"
	"path/filepath"
	"strconv"
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

func (s *UserService) CreateSuperUser(ctx context.Context, username string, email string, password string) error {
	passToHash, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	err = s.UserRepository.Create(
		ctx,
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

func (s *UserService) GetUserProfile(ctx context.Context, username string) (*dto.UserProfile, error) {
	user, err := s.UserRepository.Filter(ctx, "username = ?", username)
	if err != nil {
		return nil, err
	}

	if len(user) != 1 {
		return nil, usecase_errors.NotFoundError{Msg: "Profile not found"}
	}

	oneUser := user[0]

	if !oneUser.IsActive || oneUser.Role == enums.ANONYMOUS {
		return nil, usecase_errors.NotFoundError{Msg: "Profile not found"}
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

func (s *UserService) ChangeUserProfile(ctx context.Context, caller dto.UserDTO, data dto.ChangeUserProfileRequest) error {
	if caller.Role == enums.ANONYMOUS || !caller.IsActive {
		return usecase_errors.UnauthorizedError{Msg: "You must be logged in to change your profile"}
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

	err := s.UserRepository.UpdateById(ctx, caller.ID, updateData)
	if err != nil {
		if errors.Is(err, repositories.ErrDuplicate) {
			return usecase_errors.AlreadyExistsError{Msg: "User with this username already exists"}
		}
		return err
	}
	return nil
}

func (s *UserService) ResetPassword(ctx context.Context, request dto.ResetPasswordRequest) (int, error) {
	users, err := s.UserRepository.Filter(ctx, "email = ? OR username = ?", request.UsernameOrEmail, request.UsernameOrEmail)

	if err != nil {
		return -1, err
	}
	if len(users) != 1 {
		return -1, usecase_errors.NotFoundError{Msg: "User not found"}
	}

	user := users[0]
	if !user.IsActive || user.Role == enums.ANONYMOUS {
		return -1, usecase_errors.NotFoundError{Msg: "User not found"}
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
	_, err = s.SessionService.SetSession(ctx, sessionBody)
	if err != nil {
		s.App.Logger.Error(fmt.Sprintf("Error while set session to redis: %s", err.Error()))
		return -1, err
	}

	go s.EmailService.SendResetPasswordEmail(user.Email, sessionBody.SessionID)

	return secretCode, nil
}

func (s *UserService) ConfirmResetPassword(ctx context.Context, token string, request dto.ConfirmResetPasswordRequest) error {
	if request.NewPassword != request.ConfirmNewPassword {
		return usecase_errors.BadRequestError{Msg: "Password does not match"}
	}

	session, err := s.SessionService.GetSession(ctx, s.App.Config.RedisConfig.Prefixes.ConfirmResetPassword, token)
	if err != nil {
		return usecase_errors.BadRequestError{Msg: "Invalid token"}
	}

	var sessionBody dto.ResetPasswordSession
	err = s.SessionService.DecryptAndParsePayload(session, &sessionBody)
	if err != nil {
		return usecase_errors.BadRequestError{Msg: "Invalid token"}
	}

	if sessionBody.Code != request.Code {
		return usecase_errors.BadRequestError{Msg: "Invalid code"}
	}

	passToHash, err := utils.HashPassword(request.NewPassword)
	if err != nil {
		return err
	}

	err = s.SessionService.DeleteSession(ctx, s.App.Config.RedisConfig.Prefixes.ConfirmResetPassword, token)
	if err != nil {
		return usecase_errors.BadRequestError{Msg: "Invalid token"}
	}

	err = s.UserRepository.UpdateById(ctx, sessionBody.UserDTO.ID, map[string]any{"password": passToHash})
	if err != nil {
		return usecase_errors.BadRequestError{Msg: "Invalid token"}
	}

	return nil
}

func (s *UserService) ChangePassword(ctx context.Context, caller dto.UserDTO, request dto.ChangePasswordRequest) error {
	user, err := s.UserRepository.GetById(ctx, caller.ID)
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			return usecase_errors.BadRequestError{Msg: "Invalid session, login again"}
		}
		return err
	}

	if !user.IsActive || user.Role == enums.ANONYMOUS {
		return usecase_errors.UnauthorizedError{Msg: "You must be logged in to change"}
	} else if !utils.CheckPasswordHash(user.Password, request.OldPassword) {
		return usecase_errors.BadRequestError{Msg: "Invalid old password"}
	}

	if request.NewPassword == request.OldPassword {
		return usecase_errors.BadRequestError{Msg: "Old password cannot be the same one"}
	}

	if request.NewPassword != request.ConfirmNewPassword {
		return usecase_errors.BadRequestError{Msg: "Password does not match"}
	}

	passToHash, err := utils.HashPassword(request.NewPassword)
	if err != nil {
		return err
	}

	err = s.UserRepository.UpdateById(ctx, user.ID, map[string]any{"password": passToHash})
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			return usecase_errors.BadRequestError{Msg: "Invalid session, login again"}
		}
	}
	return nil
}

func (s *UserService) SetOnline(ctx context.Context, user dto.UserDTO) error {
	if user.Role == enums.ANONYMOUS || !user.IsActive {
		return usecase_errors.UnauthorizedError{Msg: "You must be logged"}
	}

	session := dto.SessionDTO{
		SessionID: strconv.Itoa(int(user.ID)),
		Expire:    time.Now().Add(time.Duration(s.App.Config.AuthConfig.IsOnlineTTL) * time.Second),
		Prefix:    s.App.Config.RedisConfig.Prefixes.InOnline,
		Payload:   "",
	}

	_, err := s.SessionService.SetSession(ctx, session)
	return err
}
