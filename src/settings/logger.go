package settings

import (
	"go.uber.org/zap"
	"fmt"
)


func GetLogger(config *BaseConfig) (*zap.Logger, error) {

	mode := config.AppConfig.Mode

	var cfg zap.Config
	if mode == "" {
		return nil, fmt.Errorf("mode is not set")
	}else if mode == "prod" {
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