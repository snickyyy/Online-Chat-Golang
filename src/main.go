package main

import (
	"libs/src/settings"
)

func main() {
	settings.InitBaseConfig()
	settings.InitLogger()
	settings.InitDb()

	logger := settings.GetLogger()
	logger.Info("Initialized main function")
}
