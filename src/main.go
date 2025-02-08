package main

import (
	"libs/src/settings"
)

func main() {
	diCont := settings.GetDI()
	diCont.Run()
}
