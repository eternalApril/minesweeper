package screen

import (
	"github.com/kbinani/screenshot"
	"image"
	"main/settings"
)

func GetScreen(mode settings.Settings) (*image.RGBA, error) {
	// firstPoint indicates the upper left corner of the game field, which is the constant for the left position
	firstPoint := image.Point{X: 313, Y: 210}

	width := mode.Width * settings.BlockSize
	height := mode.Height * settings.BlockSize

	img, err := screenshot.Capture(firstPoint.X, firstPoint.Y, width, height)
	if err != nil {
		return nil, err
	}
	return img, err
}
