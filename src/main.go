package main

import "libs/src/settings"

func main() {
	settings.InitBaseConfig()
	settings.InitLogger()

	logger := settings.GetLogger()
	defer logger.Sync()

	logger.Info("Initializing database...")
	settings.InitDb()
	
	logger.Info("Initializing models...")
	settings.MakeMigrations()


}
