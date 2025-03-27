package repositories

import (
	"libs/src/settings"
)

func CreateIndexes(app *settings.App) {
	app.Logger.Info("Creating indexes")
	if err := NewMessageRepository(app).CreateIndex(); err != nil {
		panic(err)
	}
}
