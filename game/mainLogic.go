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

func Start(mode settings.Settings) error {
	BombCount.Store(int32(mode.Mines))

	x := rand.Intn(mode.Height*3/4-mode.Height/4+1) + mode.Height/4
	y := rand.Intn(mode.Width*3/4-mode.Width/4+1) + mode.Width/4
	field.Click(y, x)
	err := field.Parse(mode)
	if err != nil {
		return errors.New("Start:" + err.Error())
	}

	for BombCount.Load() != 0 {
		wg := sync.WaitGroup{}

		wg.Add(2)
		go findBombAround(mode, &wg)
		go findSafeAround(mode, &wg)

		wg.Wait()
		changed := false
		for {
			select {
			case p := <-Bomb:
				field.Bomb(p.Y, p.X)
				fmt.Printf("%d %d bomb\n", p.Y, p.X)
				BombCount.Add(-1)
				changed = true
			case p := <-EmptyEl:
				field.Click(p.Y, p.X)
				fmt.Printf("%d %d empty\n", p.Y, p.X)
				changed = true
			default:
				// Каналы пусты, выходим из внутреннего цикла
				goto done
			}
		}
	done:
		if changed {
			err := field.Parse(mode)
			if err != nil {
				return errors.New("Loop:" + err.Error())
			}
		} else {
			field.PrintField(mode)
			break
		}
	}

	return nil
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
