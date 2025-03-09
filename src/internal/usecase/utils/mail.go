package utils

import (
	"libs/src/settings"
	"fmt"

	"gopkg.in/gomail.v2"
)

func SendMail(mail *gomail.Dialer, from string, to string, subject string, body string) error {
	message := gomail.NewMessage()
	message.SetHeader("From", from)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/plain", body)

	if err := mail.DialAndSend(message); err != nil {
		return err
	}
	settings.AppVar.Logger.Info(fmt.Sprintf("Email sent to %s", to))
	return nil
}