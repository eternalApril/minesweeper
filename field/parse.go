package field

import (
	"image"
	"main/settings"
	"sync"
)

type Matrix struct {
	data *[][]int
	mu   sync.RWMutex
}

func InitField(mod settings.Settings) Matrix {
	data := make([][]int, mod.Height)
	for i := range data {
		data[i] = make([]int, mod.Width)
	}
	return Matrix{&data, sync.RWMutex{}}
}

func Parse(matrix *Matrix, img *image.RGBA) {
	matrix.mu.Lock()
	defer matrix.mu.Unlock()
	// parse screenshoot
}

func CheckEnd(matrix *[][]int) bool {
	return false
}
