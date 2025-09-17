package main

import (
	"fmt"
	"main/settings"
)

func main() {
	mode, err := settings.ParseSettings()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Mode: %d\n", mode)
}
