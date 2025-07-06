package shared

import "client/network/c2s"

type Player struct {
	X      int32
	Y      int32
	RealX  int32
	RealY  int32
	ChunkX int32
	ChunkY int32
	Name   string
}

func (p *Player) Move(newX, newY int32, game *Game) {
	// TODO: send packet
	// TODO: bounds after map loading
	// basic bounds checking for testing purposes
	if newX > 15 || newX < 0 || newY > 31 || newY < 0 {
		return
	}

	p.X = newX
	p.Y = newY

	p.RealX = (p.X % int32(game.ChunkSize)) * int32(game.TileSize)
	p.RealY = (p.Y % int32(game.ChunkSize)) * int32(game.TileSize)

	p.ChunkX = p.X / int32(game.ChunkSize)
	p.ChunkY = p.Y / int32(game.ChunkSize)
}

func (p *Player) SendMovePacket(game *Game) {
	SendPacket(game.Conn, &c2s.MovePacket{
		X: p.X,
		Y: p.Y,
	})
}
