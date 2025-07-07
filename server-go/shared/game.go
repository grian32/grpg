package shared

type Game struct {
	Players []*Player
	// these will be dynamic once map loading is done and as such will be needed
	// for bounds checks.
	MaxX uint32
	MaxY uint32
}
