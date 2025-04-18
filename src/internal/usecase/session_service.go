package services

import (
	"encoding/json"
	"libs/src/internal/dto"
	"libs/src/internal/repositories"
	api_errors "libs/src/internal/usecase/errors"
	"libs/src/pkg/utils"
	"libs/src/settings"
	"time"
)

//go:generate mockery --name=ISessionService --dir=. --output=../mocks --with-expecter
type ISessionService interface {
	GetSession(prefix string, session string) (dto.SessionDTO, error)
	SetSession(session dto.SessionDTO) (string, error)
	DeleteSession(prefix string, session string) error
	DecryptAndParsePayload(session dto.SessionDTO, parseTo any) error
	GetUserByAuthSession(session string) (dto.UserDTO, error)
	GetUserByEmailSession(session string) (dto.UserDTO, error)
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

func (s *SessionService) GetSession(prefix string, session string) (dto.SessionDTO, error) {
	res, err := s.RedisBaseRepository.GetByKey(prefix, session)
	if err != nil {
		return dto.SessionDTO{}, api_errors.ErrInvalidSession
	}

	var sessionBody dto.SessionDTO

	err = json.Unmarshal([]byte(res), &sessionBody)
	if err != nil {
		return dto.SessionDTO{}, api_errors.ErrInvalidSession
	}

	return sessionBody, nil
}

func (s *SessionService) SetSession(session dto.SessionDTO) (string, error) {
	encoding, _ := json.Marshal(&session)

	_, err := s.RedisBaseRepository.Create(
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

func (s *SessionService) DeleteSession(prefix string, session string) error {
	_, err := s.RedisBaseRepository.Delete(prefix, session)
	if err != nil {
		return api_errors.ErrInvalidSession
	}

	return nil
}

func (s *SessionService) DecryptAndParsePayload(session dto.SessionDTO, parseTo any) error {
	decryptResult, err := utils.Decrypt(s.App.Config.AppConfig.SecretKey, session.Payload)
	if err != nil {
		return api_errors.ErrInvalidSession
	}

	err = json.Unmarshal([]byte(decryptResult), &parseTo)
	if err != nil {
		return api_errors.ErrInvalidSession
	}

	return nil
}

func (s *SessionService) GetUserByAuthSession(session string) (dto.UserDTO, error) {
	sessionBody, err := s.GetSession(s.App.Config.RedisConfig.Prefixes.SessionPrefix, session)
	if err != nil {
		return dto.UserDTO{}, api_errors.ErrInvalidSession
	}
	var authSessionBody dto.AuthSession
	err = s.DecryptAndParsePayload(sessionBody, &authSessionBody)
	if err != nil {
		return dto.UserDTO{}, api_errors.ErrInvalidSession
	}

	return authSessionBody.UserDTO, nil
}

func (s *SessionService) GetUserByEmailSession(session string) (dto.UserDTO, error) {
	sessionBody, err := s.GetSession(s.App.Config.RedisConfig.Prefixes.ConfirmEmail, session)
	if err != nil {
		return dto.UserDTO{}, api_errors.ErrInvalidSession
	}

	var emailSessionBody dto.EmailSession
	err = s.DecryptAndParsePayload(sessionBody, &emailSessionBody)
	if err != nil {
		return dto.UserDTO{}, api_errors.ErrInvalidSession
	}

	return emailSessionBody.UserDTO, nil
}
