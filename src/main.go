package main

import (
	domain "libs/src/internal/domain/models"
	"libs/src/internal/repositories"
	"libs/src/settings"
	_ "time"

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

	repo := repositories.UserRepository{}
	repo.Model = domain.User{}
	repo.Db = settings.AppVar.DB

	err = repo.ExecuteQuery("DROP TABLE users")

	if err != nil {
		panic(err)
	}

	if err := diCont.Stop(settings.Context.Ctx); err != nil {
		panic(err)
	}
}
