package main

import (
	"fmt"
	"libs/src/internal/repositories"
	"libs/src/settings"

	"go.mongodb.org/mongo-driver/bson"
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
	
	res, _ := repo.Count(bson.M{"owner_id": 1})
	fmt.Printf("Total chats: %d\n", res)

	resf, _ := repo.GetOne(bson.M{"owner_id": 1})
	fmt.Println(resf)

	resa, _ := repo.GetAll(bson.M{"owner_id": 1}, 0,10)
	fmt.Println(resa)

	if err := diCont.Stop(settings.Context.Ctx); err != nil {
		panic(err)
	}
}
