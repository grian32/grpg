package shared

import "fmt"

type Direction byte

const (
	UP Direction = iota
	RIGHT
	DOWN
	LEFT
)

func (d Direction) String() string {
	switch d {
	case DOWN:
		return "Down"
	case LEFT:
		return "Left"
	case RIGHT:
		return "Right"
	case UP:
		return "Up"
	default:
		panic(fmt.Sprintf("unexpected shared.Direction: %#v", d))
	}
}
