package field

import (
	"errors"
	"fmt"
	"image"
	"main/screen"
	"main/settings"
	"sync"
)

var (
	Field = Matrix{}
)

type Matrix struct {
	Data [][]int
	Mu   sync.RWMutex
}

func InitField(mod settings.Settings) {
	Field.Mu.Lock()
	defer Field.Mu.Unlock()

	Field.Data = make([][]int, mod.Height)
	for i := range Field.Data {
		Field.Data[i] = make([]int, mod.Width)
	}
}

// elemToUpLeftPixel return coordinates of up left corner of element field
func elemToUpLeftPixel(y, x int) image.Point {
	var point image.Point
	point.X += x * settings.BlockSize
	point.Y += y * settings.BlockSize
	return point
}

// ElemToPixel return coordinates of center element
func ElemToPixel(y, x int) image.Point {
	var point image.Point
	point.X += x*settings.BlockSize + settings.BlockSize/2
	point.Y += y*settings.BlockSize + settings.BlockSize/2
	return point
}

func getInfoBlock(y, x int, img *image.RGBA) int {
	point := elemToUpLeftPixel(y, x)
	color := img.RGBAAt(point.X, point.Y)
	if color == settings.White {
		return -1
	}
	point = ElemToPixel(y, x)
	color = img.RGBAAt(point.X, point.Y)
	switch color {
	case settings.Gray:
		return 0
	case settings.Blue:
		return 1
	case settings.Green:
		return 2
	case settings.Red:
		return 3
	case settings.DarkBlue:
		return 4
	case settings.DarkRed:
		return 5
	case settings.Navy:
		return 6
	case settings.Black:
		return -2
	default:
		fmt.Println(color)
	}
	return 69
}

func Parse(mode settings.Settings) error {
	Field.Mu.Lock()
	defer Field.Mu.Unlock()

	img, err := screen.GetScreen(mode)
	if err != nil {
		return errors.New("error getting screen")
	}

	for i := 0; i != mode.Height; i++ {
		for j := 0; j != mode.Width; j++ {
			if Field.Data[i][j] != -2 { // check if element is a bomb
				Field.Data[i][j] = getInfoBlock(i, j, img)
				if Field.Data[i][j] == -2 {
					return errors.New("bomb click")
				}
			}
		}
	}
	return nil
}
