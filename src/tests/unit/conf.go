package unit

import (
	"libs/src/settings"
	"libs/src/tests/integration"
)

func GetAppMock() *settings.App {
	cfg := integration.GetTestConfig()
	logger, err := settings.GetLogger(cfg)
	if err != nil {
		panic(err)
	}
	return &settings.App{
		Config: cfg,
		Logger: logger,
	}
}
