package main

import (
	"fmt"
	"libs/src/settings"
)

func main() {
	cfg := settings.GetBaseConfig()
	fmt.Printf("%+v\n", cfg)
}
