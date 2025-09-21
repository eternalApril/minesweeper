package main

import (
	"main/field"
	"main/game"
	"main/settings"
)

func main() {
	mode, err := settings.ParseSettings()
	if err != nil {
		panic(err)
	}

	field.InitField(mode)
	if err = game.Start(mode); err != nil {
		panic(err)
	}
}
