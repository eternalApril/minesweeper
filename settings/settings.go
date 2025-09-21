package settings

import (
	"errors"
	"flag"
	"image/color"
)

type GameMod int

var (
	White    = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	Gray     = color.RGBA{R: 189, G: 189, B: 189, A: 255}
	Blue     = color.RGBA{R: 0, G: 0, B: 255, A: 255}
	Green    = color.RGBA{R: 2, G: 124, B: 2, A: 255}
	Red      = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	DarkBlue = color.RGBA{R: 0, G: 0, B: 123, A: 255}
	DarkRed  = color.RGBA{R: 123, G: 0, B: 0, A: 255}
	Navy     = color.RGBA{R: 0, G: 123, B: 123, A: 255}
)

const (
	BlockSize = 20 // block is a square with a side of 20

	Beginner GameMod = iota + 1
	Intermediate
	Expert
)

type Settings struct {
	Mode   GameMod
	Height int
	Width  int
	Mines  int
}

func ParseSettings() (Settings, error) {
	var settings Settings

	var beginner bool
	var intermediate bool
	var expert bool

	flag.BoolVar(&beginner, "beginner", false, "Set game mode to Beginner")
	flag.BoolVar(&beginner, "b", false, "Alias for --beginner")

	flag.BoolVar(&intermediate, "intermediate", false, "Set game mode to Intermediate")
	flag.BoolVar(&intermediate, "i", false, "Alias for --intermediate")

	flag.BoolVar(&expert, "expert", false, "Set game mode to Expert")
	flag.BoolVar(&expert, "e", false, "Alias for --expert")

	flag.Parse()

	count := 0
	if beginner {
		count++
	}
	if intermediate {
		count++
	}
	if expert {
		count++
	}

	if count > 1 {
		return Settings{}, errors.New("can choose only one game mode")
	}
	if count == 0 {
		return Settings{}, errors.New("need to set game mode")
	}

	switch {
	case beginner:
		settings = Settings{
			Mode: Beginner, Height: 9, Width: 9, Mines: 10,
		}
	case intermediate:
		settings = Settings{
			Mode: Intermediate, Height: 16, Width: 16, Mines: 40,
		}
	case expert:
		settings = Settings{
			Mode: Expert, Height: 16, Width: 30, Mines: 99,
		}
	}

	return settings, nil
}
