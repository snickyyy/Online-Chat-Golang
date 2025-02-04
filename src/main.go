package main

import (
	"libs/src/settings"
)

func main() {
	settings.InitBaseConfig()
	settings.InitLogger()

	logger := settings.GetLogger()
	logger.Info("Initialized main function")
}
