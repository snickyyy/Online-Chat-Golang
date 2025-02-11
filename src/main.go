package main

import (
	"fmt"
	"libs/src/settings"
)

func init() {
	settings.InitContext()
}

func main() {
	diCont := settings.GetDI()
	diCont.Run()
	fmt.Println(&settings.Context)
}
