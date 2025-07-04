package game

type Player struct {
	X     int32
	Y     int32
	RealX int32
	RealY int32
	Name  string
}

func (p *Player) Move(newX, newY int32, game *Game) {
	// TODO: send packet
	// TODO: bounds after map loading
	// basic bounds checking for testing purposes
	if newX > 15 || newX < 0 || newY > 15 || newY < 0 {
		return
	}

	p.X = newX
	p.Y = newY

	p.RealX = p.X * game.TileSize
	p.RealY = p.Y * game.TileSize
}
