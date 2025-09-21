package screen

import (
	"github.com/kbinani/screenshot"
	"image"
	"main/settings"
)

var (
	// FirstPoint indicates the upper left corner of the game field, which is the constant for the left position
	FirstPoint = image.Point{X: 313, Y: 210}
)

func GetScreen(mode settings.Settings) (*image.RGBA, error) {
	width := mode.Width * settings.BlockSize
	height := mode.Height * settings.BlockSize

	img, err := screenshot.Capture(FirstPoint.X, FirstPoint.Y, width, height)
	if err != nil {
		return nil, err
	}
	return img, err
}
