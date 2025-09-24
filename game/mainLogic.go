package game

import (
	"errors"
	"fmt"
	"image"
	"main/field"
	"main/settings"
	"math/rand"
	"sync"
	"sync/atomic"
)

var (
	Bomb      = make(chan image.Point, 100)
	EmptyEl   = make(chan image.Point, 100)
	BombCount atomic.Int32
)

func saveBomb(y, x int) {
	field.Field.Mu.Lock()
	defer field.Field.Mu.Unlock()

	if field.Field.Data[y][x] != -2 {
		field.Field.Data[y][x] = -2
		BombCount.Add(-1)
	}
}

func Start(mode settings.Settings) error {
	BombCount.Store(int32(mode.Mines))

	// Click somewhere in the middle
	y := rand.Intn(mode.Height*3/4-mode.Height/4+1) + mode.Height/4
	x := rand.Intn(mode.Width*3/4-mode.Width/4+1) + mode.Width/4
	field.Click(y, x)

	err := field.Parse(mode)
	if err != nil {
		return errors.New("Start:" + err.Error())
	}

	for BombCount.Load() != 0 {
		wg := sync.WaitGroup{}

		wg.Add(3)
		go findBombAround(mode, &wg)
		go findSafeAround(mode, &wg)
		go find1221Pattern(mode, &wg)

		wg.Wait()
		changed := false
		for {
			select {
			case p := <-Bomb:
				saveBomb(p.Y, p.X)
				changed = true
			case p := <-EmptyEl:
				field.Click(p.Y, p.X)
				changed = true
			default:
				if !changed {
					selectRandomCell(mode)
				}
				goto done
			}
		}
	done:
		err := field.Parse(mode)
		if err != nil {
			return errors.New("Loop:" + err.Error())
		}
		changed = false
	}

	return nil
}

func selectRandomCell(mode settings.Settings) {
	msg := "Random cell click:"
	point, err := findBorderUnknownCells(mode)
	if err != nil {
		msg = "Full random cell click:"
		point = fullRandomCells(mode)
	}
	fmt.Println(msg, point)
	field.Click(point.Y, point.X)
}

func findBorderUnknownCells(mode settings.Settings) (image.Point, error) {
	field.Field.Mu.Lock()
	defer field.Field.Mu.Unlock()

	height := mode.Height
	width := mode.Width

	var borderCells []image.Point
	visited := make(map[image.Point]bool)

	directions := []image.Point{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			if field.Field.Data[i][j] >= 0 {
				for _, dir := range directions {
					ni, nj := i+dir.Y, j+dir.X
					if ni >= 0 && ni < height && nj >= 0 && nj < width {
						if field.Field.Data[ni][nj] == -1 {
							point := image.Point{X: nj, Y: ni}
							if !visited[point] {
								borderCells = append(borderCells, point)
								visited[point] = true
							}
						}
					}
				}
			}
		}
	}
	if len(borderCells) == 0 {
		return image.Point{}, errors.New("no border cells found")
	}
	randomCell := borderCells[rand.Intn(len(borderCells))]

	return randomCell, nil
}

func fullRandomCells(mode settings.Settings) image.Point {
	field.Field.Mu.RLock()
	defer field.Field.Mu.RUnlock()

	var unknownCells []image.Point
	for i := 1; i < mode.Height-1; i++ {
		for j := 1; j < mode.Width-1; j++ {
			if field.Field.Data[i][j] == -1 {
				unknownCells = append(unknownCells, image.Point{X: j, Y: i})
			}
		}
	}
	randomCell := unknownCells[rand.Intn(len(unknownCells))]
	return randomCell
}

func findBombAround(mode settings.Settings, wg *sync.WaitGroup) {
	field.Field.Mu.RLock()
	defer field.Field.Mu.RUnlock()

	height := mode.Height
	width := mode.Width

	dy := []int{-1, -1, -1, 0, 0, 1, 1, 1}
	dx := []int{-1, 0, 1, -1, 1, -1, 0, 1}

	uniqueSuspects := make(map[image.Point]struct{})

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			value := field.Field.Data[i][j]
			if value > 0 && value <= 8 {
				flagCount := 0
				unknownCount := 0
				var suspects []image.Point
				for k := 0; k < 8; k++ {
					ny := i + dy[k]
					nx := j + dx[k]
					if ny >= 0 && ny < height && nx >= 0 && nx < width {
						cell := field.Field.Data[ny][nx]
						if cell == -1 {
							unknownCount++
							suspects = append(suspects, image.Point{X: nx, Y: ny})
						} else if cell == -2 {
							flagCount++
						}
					}
				}
				if unknownCount > 0 && unknownCount == value-flagCount {
					for _, p := range suspects {
						uniqueSuspects[p] = struct{}{}
					}
				}
			}
		}
	}

	for p := range uniqueSuspects {
		Bomb <- p
	}
	wg.Done()
}

func findSafeAround(mode settings.Settings, wg *sync.WaitGroup) {
	field.Field.Mu.RLock()
	defer field.Field.Mu.RUnlock()

	height := mode.Height
	width := mode.Width

	dy := []int{-1, -1, -1, 0, 0, 1, 1, 1}
	dx := []int{-1, 0, 1, -1, 1, -1, 0, 1}

	uniqueSafes := make(map[image.Point]struct{})

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			value := field.Field.Data[i][j]
			if value >= 0 {
				flagCount := 0
				unknownCount := 0
				var safes []image.Point
				for k := 0; k < 8; k++ {
					ny := i + dy[k]
					nx := j + dx[k]
					if ny >= 0 && ny < height && nx >= 0 && nx < width {
						cell := field.Field.Data[ny][nx]
						if cell == -1 {
							unknownCount++
							safes = append(safes, image.Point{X: nx, Y: ny})
						} else if cell == -2 {
							flagCount++
						}
					}
				}
				if unknownCount > 0 && flagCount == value {
					for _, p := range safes {
						uniqueSafes[p] = struct{}{}
					}
				}
			}
		}
	}

	for p := range uniqueSafes {
		EmptyEl <- p
	}
	wg.Done()
}

func find1221Pattern(mode settings.Settings, wg *sync.WaitGroup) {
	field.Field.Mu.RLock()
	defer field.Field.Mu.RUnlock()

	height := mode.Height
	width := mode.Width

	uniqueBombs := make(map[image.Point]struct{})
	uniqueSafes := make(map[image.Point]struct{})

	for i := 0; i < height; i++ {
		for j := 0; j < width-3; j++ {
			if field.Field.Data[i][j] == 1 && field.Field.Data[i][j+1] == 2 && field.Field.Data[i][j+2] == 2 && field.Field.Data[i][j+3] == 1 {
				if i > 0 {
					if field.Field.Data[i-1][j] == -1 && field.Field.Data[i-1][j+1] == -1 && field.Field.Data[i-1][j+2] == -1 && field.Field.Data[i-1][j+3] == -1 {
						uniqueSafes[image.Point{Y: i - 1, X: j}] = struct{}{}
						uniqueBombs[image.Point{Y: i - 1, X: j + 1}] = struct{}{}
						uniqueBombs[image.Point{Y: i - 1, X: j + 2}] = struct{}{}
						uniqueSafes[image.Point{Y: i - 1, X: j + 3}] = struct{}{}
					}
				}
				if i < height-1 {
					if field.Field.Data[i+1][j] == -1 && field.Field.Data[i+1][j+1] == -1 && field.Field.Data[i+1][j+2] == -1 && field.Field.Data[i+1][j+3] == -1 {
						uniqueSafes[image.Point{Y: i + 1, X: j}] = struct{}{}
						uniqueBombs[image.Point{Y: i + 1, X: j + 1}] = struct{}{}
						uniqueBombs[image.Point{Y: i + 1, X: j + 2}] = struct{}{}
						uniqueSafes[image.Point{Y: i + 1, X: j + 3}] = struct{}{}
					}
				}
			}
		}
	}

	for j := 0; j < width; j++ {
		for i := 0; i < height-3; i++ {
			if field.Field.Data[i][j] == 1 && field.Field.Data[i+1][j] == 2 && field.Field.Data[i+2][j] == 2 && field.Field.Data[i+3][j] == 1 {
				if j > 0 {
					if field.Field.Data[i][j-1] == -1 && field.Field.Data[i+1][j-1] == -1 && field.Field.Data[i+2][j-1] == -1 && field.Field.Data[i+3][j-1] == -1 {
						uniqueSafes[image.Point{Y: i, X: j - 1}] = struct{}{}
						uniqueBombs[image.Point{Y: i + 1, X: j - 1}] = struct{}{}
						uniqueBombs[image.Point{Y: i + 2, X: j - 1}] = struct{}{}
						uniqueSafes[image.Point{Y: i + 3, X: j - 1}] = struct{}{}
					}
				}
				if j < width-1 {
					if field.Field.Data[i][j+1] == -1 && field.Field.Data[i+1][j+1] == -1 && field.Field.Data[i+2][j+1] == -1 && field.Field.Data[i+3][j+1] == -1 {
						uniqueSafes[image.Point{Y: i, X: j + 1}] = struct{}{}
						uniqueBombs[image.Point{Y: i + 1, X: j + 1}] = struct{}{}
						uniqueBombs[image.Point{Y: i + 2, X: j + 1}] = struct{}{}
						uniqueSafes[image.Point{Y: i + 3, X: j + 1}] = struct{}{}
					}
				}
			}
		}
	}

	for p := range uniqueBombs {
		Bomb <- p
	}
	for p := range uniqueSafes {
		EmptyEl <- p
	}
	wg.Done()
}
