package main

// https://minesweeperonline.com

import (
	"errors"
	"fmt"
	"github.com/kbinani/screenshot"
	"github.com/lxn/win"
	"image"
	"os"
	"time"
	"unsafe"
)

const (
	blockLength  = 20
	beginner     = 0
	intermediate = 1
	expert       = 2
	mine         = 10
)

type mode int

var (
	Mode       mode
	Cols       int
	Rows       int
	firstPoint image.Point
	field      [][]int8
)

func initMode(mode string) error {
	if mode == "beginner" || mode == "b" {
		Mode = beginner
		Cols = 9
		Rows = 9
		firstPoint = image.Point{X: 668, Y: 210}
	} else if mode == "intermediate" || mode == "i" {
		Mode = intermediate
		Cols = 16
		Rows = 16
		firstPoint = image.Point{X: 668, Y: 210}
	} else if mode == "expert" || mode == "e" {
		Mode = expert
		Cols = 30
		Rows = 16
		firstPoint = image.Point{X: 655, Y: 210}
	} else {
		return errors.New("invalid mode")
	}
	return nil
}

func takeScreen() (*image.RGBA, error) {
	var secondPoint image.Point
	switch Mode {
	case beginner:
		secondPoint = image.Point{X: firstPoint.X + blockLength*9, Y: firstPoint.Y + blockLength*9}
	case intermediate:
		secondPoint = image.Point{X: firstPoint.X + blockLength*16, Y: firstPoint.Y + blockLength*16}
	case expert:
		secondPoint = image.Point{X: firstPoint.X + blockLength*30, Y: firstPoint.Y + blockLength*16}
	}
	var bounds = image.Rectangle{
		Min: firstPoint,
		Max: secondPoint}
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		return nil, err
	}
	return img, err
}

func readGame(rgba *image.RGBA) {
	for i := 0; i < Rows; i++ {
		for j := 0; j < Cols; j++ {
			if field[i][j] == mine {
				continue
			}
			r, _, _, _ := rgba.At(blockLength*j, blockLength*i).RGBA()
			if r == 65535 {
				field[i][j] = -1
			} else {
				r, g, b, _ := rgba.At(10+blockLength*j, 10+blockLength*i).RGBA()
				if g == 48573 {
					field[i][j] = 0
				} else if b == 65535 {
					field[i][j] = 1
				} else if r == 514 {
					field[i][j] = 2
				} else if r == 65535 {
					field[i][j] = 3
				} else if b == 31611 && g == 0 {
					field[i][j] = 4
				} else if r == 31611 {
					field[i][j] = 5
				} else if b == 31611 && g == 31611 {
					field[i][j] = 6
				}
			}
		}
	}
}

func clickElem(x, y int) {
	blockLength := blockLength / 5 * 4 // scale monitor 125% on windows
	var absX, absY int32
	if Mode == expert {
		absX, absY = 532, 176
	} else {
		absX, absY = 542, 176
	}
	absX += int32(x * blockLength)
	absY += int32(y * blockLength)
	if !win.SetCursorPos(absX, absY) {
		fmt.Printf("Error moving mouse to (%d, %d)\n", absX, absY)
		return
	}
	type MOUSEINPUT struct {
		Dx          int32
		Dy          int32
		MouseData   uint32
		DwFlags     uint32
		Time        uint32
		DwExtraInfo uintptr
	}

	type INPUT struct {
		Type uint32
		Mi   MOUSEINPUT
	}

	const (
		INPUT_MOUSE          = 0
		MOUSEEVENTF_LEFTDOWN = 0x0002
		MOUSEEVENTF_LEFTUP   = 0x0004
	)
	mouseDown := INPUT{
		Type: INPUT_MOUSE,
		Mi: MOUSEINPUT{
			DwFlags: MOUSEEVENTF_LEFTDOWN,
		},
	}
	mouseUp := INPUT{
		Type: INPUT_MOUSE,
		Mi: MOUSEINPUT{
			DwFlags: MOUSEEVENTF_LEFTUP,
		},
	}
	inputs := []INPUT{mouseDown, mouseUp}
	ret := win.SendInput(uint32(len(inputs)), unsafe.Pointer(&inputs[0]), int32(unsafe.Sizeof(mouseDown)))
	if ret != uint32(len(inputs)) {
		fmt.Println("SendInput failed, events sent:", ret)
	}
	time.Sleep(time.Millisecond * 50)
}

func checkEnd() bool {
	for i := 0; i < Rows; i++ {
		for j := 0; j < Cols; j++ {
			if field[i][j] == -1 {
				return false
			}
		}
	}
	return true
}

func game() {
	clickElem(Cols/2, Rows/2)
	img, _ := takeScreen()
	readGame(img)
	for !checkEnd() {
		findBomb()
		safeCells := findSafeCells()
		if len(safeCells) == 0 {
			break
		}
		for _, point := range safeCells {
			clickElem(point.J, point.I)
		}
		img, _ = takeScreen()
		readGame(img)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Choose type of game")
		os.Exit(1)
	}
	err := initMode(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	field = make([][]int8, Rows)
	for i := 0; i < Rows; i++ {
		field[i] = make([]int8, Cols)
	}
	game()
}
