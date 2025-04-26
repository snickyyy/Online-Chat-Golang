package main

import (
	"context"
	"libs/src/internal/repositories"
	"libs/src/settings"
	"libs/src/settings/server"
)

func init() {
	settings.AppVar = &settings.App{}
	ctx, cancel := context.WithCancel(context.Background())
	settings.AppVar.Ctx = ctx
	settings.AppVar.Cancel = cancel
}

func main() {
	defer settings.AppVar.Cancel()

	diCont := settings.GetDI()
	err := diCont.Start(settings.AppVar.Ctx)
	if err != nil {
		panic(err)
	}

	repositories.CreateIndexes(settings.AppVar)

	server.RunServer()

	if err := diCont.Stop(settings.AppVar.Ctx); err != nil {
		panic(err)
	}
}
