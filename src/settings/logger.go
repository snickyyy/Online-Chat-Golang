package settings

import (
	"fmt"
	"go.uber.org/zap"
)

func GetLogger(config *BaseConfig) (*zap.Logger, error) {

	mode := config.AppConfig.Mode

	var cfg zap.Config
	if mode == "" { //nolint:staticcheck
		return nil, fmt.Errorf("mode is not set")
	} else if mode == "prod" {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.CallerKey = "caller"
	}

	logger, err := cfg.Build()

	if err != nil {
		return nil, fmt.Errorf("error initializing logger: %v", err)
	}

	return logger, nil
}
