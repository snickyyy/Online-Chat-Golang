package services

import (
	"context"
	"encoding/json"
	"libs/src/internal/dto"
	"libs/src/internal/repositories"
	usecase_errors "libs/src/internal/usecase/errors"
	"libs/src/pkg/utils"
	"libs/src/settings"
	"time"
)

//go:generate mockery --name=ISessionService --dir=. --output=../mocks --with-expecter
type ISessionService interface {
	GetSession(ctx context.Context, prefix string, session string) (dto.SessionDTO, error)
	SetSession(ctx context.Context, session dto.SessionDTO) (string, error)
	DeleteSession(ctx context.Context, prefix string, session string) error
	DecryptAndParsePayload(session dto.SessionDTO, parseTo any) error
	GetUserByAuthSession(ctx context.Context, session string) (dto.UserDTO, error)
	GetUserByEmailSession(ctx context.Context, session string) (dto.UserDTO, error)
}

type SessionService struct {
	App                 *settings.App
	RedisBaseRepository repositories.IBaseRedisRepository
}

func NewSessionService(app *settings.App) *SessionService {

	return &SessionService{
		App:                 app,
		RedisBaseRepository: repositories.NewBaseRedisRepository(app),
	}
}

func (s *SessionService) GetSession(ctx context.Context, prefix string, session string) (dto.SessionDTO, error) {
	res, err := s.RedisBaseRepository.GetByKey(ctx, prefix, session)
	if err != nil {
		return dto.SessionDTO{}, usecase_errors.BadRequestError{Msg: "Invalid session"}
	}

	var sessionBody dto.SessionDTO

	err = json.Unmarshal([]byte(res), &sessionBody)
	if err != nil {
		return dto.SessionDTO{}, usecase_errors.BadRequestError{Msg: "Invalid session"}
	}

	return sessionBody, nil
}

func (s *SessionService) SetSession(ctx context.Context, session dto.SessionDTO) (string, error) {
	encoding, _ := json.Marshal(&session)

	_, err := s.RedisBaseRepository.Create(
		ctx,
		session.Prefix,
		session.SessionID,
		string(encoding),
		time.Until(session.Expire),
	)
	if err != nil {
		return "", err
	}

	return session.SessionID, nil
}

func (s *SessionService) DeleteSession(ctx context.Context, prefix string, session string) error {
	_, err := s.RedisBaseRepository.Delete(ctx, prefix, session)
	if err != nil {
		return usecase_errors.BadRequestError{Msg: "Invalid session"}
	}

	return nil
}

func (s *SessionService) DecryptAndParsePayload(session dto.SessionDTO, parseTo any) error {
	decryptResult, err := utils.Decrypt(s.App.Config.AppConfig.SecretKey, session.Payload)
	if err != nil {
		return usecase_errors.BadRequestError{Msg: "Invalid session"}
	}

	err = json.Unmarshal([]byte(decryptResult), &parseTo)
	if err != nil {
		return usecase_errors.BadRequestError{Msg: "Invalid session"}
	}

	return nil
}

func (s *SessionService) GetUserByAuthSession(ctx context.Context, session string) (dto.UserDTO, error) {
	sessionBody, err := s.GetSession(ctx, s.App.Config.RedisConfig.Prefixes.SessionPrefix, session)
	if err != nil {
		return dto.UserDTO{}, usecase_errors.BadRequestError{Msg: "Invalid session"}
	}
	var authSessionBody dto.AuthSession
	err = s.DecryptAndParsePayload(sessionBody, &authSessionBody)
	if err != nil {
		return dto.UserDTO{}, usecase_errors.BadRequestError{Msg: "Invalid session"}
	}

	return authSessionBody.UserDTO, nil
}

func (s *SessionService) GetUserByEmailSession(ctx context.Context, session string) (dto.UserDTO, error) {
	sessionBody, err := s.GetSession(ctx, s.App.Config.RedisConfig.Prefixes.ConfirmEmail, session)
	if err != nil {
		return dto.UserDTO{}, usecase_errors.BadRequestError{Msg: "Invalid session"}
	}

	var emailSessionBody dto.EmailSession
	err = s.DecryptAndParsePayload(sessionBody, &emailSessionBody)
	if err != nil {
		return dto.UserDTO{}, usecase_errors.BadRequestError{Msg: "Invalid session"}
	}

	return emailSessionBody.UserDTO, nil
}
