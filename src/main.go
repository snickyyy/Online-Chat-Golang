package main

import (
	"libs/src/settings"
	"libs/src/settings/server"
)

func init() {
	settings.InitContext()
}

func main() {
	defer settings.Context.Cancel()

	diCont := settings.GetDI()
	err := diCont.Start(settings.Context.Ctx)
	if err != nil {
		panic(err)
	}

	server.RunServer()

	if err := diCont.Stop(settings.Context.Ctx); err != nil {
		panic(err)
	}
}
