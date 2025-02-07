package main

import (
	"fmt"
)

func main() {
	app, err := InitializeApp()
	if err != nil {
		panic(err)
	}
	fmt.Println(app)
}
