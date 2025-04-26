package unit

import (
	"context"
	"gopkg.in/gomail.v2"
	"libs/src/settings"
	"libs/src/tests/integration"
)

func GetAppMock() *settings.App {
	cfg := integration.GetTestConfig()
	ctx, cancel := context.WithCancel(context.Background())
	logger, err := settings.GetLogger(cfg)
	if err != nil {
		panic(err)
	}
	return &settings.App{
		Ctx:    ctx,
		Cancel: cancel,
		Config: cfg,
		Logger: logger,
		Mail:   &gomail.Dialer{},
	}
}
