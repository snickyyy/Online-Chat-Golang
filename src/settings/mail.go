package settings

import "gopkg.in/gomail.v2"

func NewMail(config *BaseConfig) *gomail.Dialer {
	return gomail.NewDialer(
		config.Mail.Server, config.Mail.Port, config.Mail.Username, config.Mail.Password,
	)
}
