package services

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"libs/src/settings"
)

type IEmailService interface {
	SendRegisterEmail(to string, token string) error
	SendResetPasswordEmail(to string, token string) error
}

type EmailService struct {
	App *settings.App
}

func NewEmailService(app *settings.App) *EmailService {
	return &EmailService{
		App: app,
	}
}

func (s *EmailService) sendEmail(from string, to string, subject string, body string) error {
	message := gomail.NewMessage()
	message.SetHeader("From", from)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/plain", body)

	if err := s.App.Mail.DialAndSend(message); err != nil {
		return err
	}
	settings.AppVar.Logger.Info(fmt.Sprintf("Email sent to %s", to))
	return nil
}

func (s *EmailService) SendRegisterEmail(to string, token string) error {
	subject := "Online-Chat-Golang || Confirm registration"
	body := fmt.Sprintf(
		"Thank you for choosing our service, to confirm your registration, follow the url below\nhttp://%s:%d/accounts/auth/confirm-account/%s",
		s.App.Config.AppConfig.DomainName,
		s.App.Config.AppConfig.Port,
		token,
	)
	return s.sendEmail(s.App.Config.Mail.From, to, subject, body)
}

func (s *EmailService) SendResetPasswordEmail(to string, token string) error {
	subject := "Reset your password"
	body := fmt.Sprintf(
		"For To confirm password reset, follow the link\nhttp://%s:%d/accounts/profile/reset-password/%s",
		s.App.Config.AppConfig.DomainName,
		s.App.Config.AppConfig.Port,
		token,
	)
	return s.sendEmail(s.App.Config.Mail.From, to, subject, body)
}
