package shared

type Direction byte

const (
	UP Direction = iota
	RIGHT
	DOWN
	LEFT
)

func DirectionString(direction Direction) string {
	switch direction {
	case UP:
		return "UP"
	case RIGHT:
		return "RIGHT"
	case DOWN:
		return "DOWN"
	case LEFT:
		return "LEFT"
	}

	return "UNKNOWN"
}
