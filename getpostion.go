package translate4game

import (
	"github.com/go-vgo/robotgo"
)

type Position struct {
	X int
	Y int
}

func GetCurrentPosition() *Position {
	p := &Position{}
	p.X, p.Y = robotgo.GetMousePos()
	return p
}
