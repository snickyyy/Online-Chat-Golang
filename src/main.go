package main

import (
	"fmt"
	"libs/src/internal/domain/models"
	"libs/src/internal/repositories"
	"libs/src/settings"
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

	repo := repositories.ChatRepository{
		Db: settings.AppVar.MongoDB,
	}
	res, err := repo.Create(domain.Chat{
		OwnerId: 1,
		Title: "blabla",
		Description: "blablablabla4",
		Members: []int64{1,2,3,4,5,9},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(res)

	if err := diCont.Stop(settings.Context.Ctx); err != nil {
		panic(err)
	}
}
