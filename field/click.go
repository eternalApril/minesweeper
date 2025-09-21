package field

import (
	"github.com/go-vgo/robotgo"
	"main/screen"
	"time"
)

func Click(y, x int) {
	Field.Mu.Lock()
	defer Field.Mu.Unlock()
	point := screen.FirstPoint
	add := ElemToPixel(y, x)
	point = point.Add(add)
	robotgo.Move(point.X, point.Y)
	robotgo.Click("left")
	time.Sleep(30 * time.Millisecond)
}
