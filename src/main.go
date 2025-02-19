package main

import (
	"libs/src/settings"

	_ "go.mongodb.org/mongo-driver/bson"
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

	settings.RunServer()

	if err := diCont.Stop(settings.Context.Ctx); err != nil {
		panic(err)
    }
}
