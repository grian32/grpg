package shared

import (
	"client/network/c2s"
	"client/util"
	"fmt"
)

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
	p.X = newX
	p.Y = newY

	p.RealX = (p.X % 16) * game.TileSize
	p.RealY = (p.Y % 16) * game.TileSize

	p.ChunkX = p.X / 16
	p.ChunkY = p.Y / 16
}

func (p *Player) SendMovePacket(game *Game, x, y int32) {
	fmt.Println(game.CollisionMap)
	_, exists := game.CollisionMap[util.Vector2I{X: x, Y: y}]
	if x > int32(game.MaxX) || x < 0 || y > int32(game.MaxY) || y < 0 || exists {
		return
	}

	SendPacket(game.Conn, &c2s.MovePacket{
		X: x,
		Y: y,
	})
}
