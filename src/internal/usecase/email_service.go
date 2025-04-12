package services

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"libs/src/settings"
	"time"
)

//go:generate mockery --name=IEmailService --dir=. --output=../mocks --with-expecter
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

	var err error
	for i := 0; i <= 3; i++ {
		err = s.App.Mail.DialAndSend(message)
		if err == nil {
			break
		}
		settings.AppVar.Logger.Error(fmt.Sprintf("Error sending email: %v || try: %d", err, i))
		time.Sleep(time.Second * time.Duration(i))
	}

	settings.AppVar.Logger.Info(fmt.Sprintf("Email sent to %s", to))
	return err
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
	subject := "Online-Chat-Golang || Reset-password"
	body := fmt.Sprintf(
		"For To confirm password reset, follow the link\nhttp://%s:%d/accounts/profile/reset-password/%s",
		s.App.Config.AppConfig.DomainName,
		s.App.Config.AppConfig.Port,
		token,
	)
	return s.sendEmail(s.App.Config.Mail.From, to, subject, body)
}
