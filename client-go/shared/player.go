package shared

import (
	"client/network/c2s"
	"client/util"
)

// TODO: separate this to a local/remote player :S, also, figure out wtf im doing with that move function lol, sob. its only used in one place
type LocalPlayer struct {
	X, Y           int32
	PrevX, PrevY   int32
	RealX, RealY   int32
	ChunkX, ChunkY int32
	Facing         Direction
	Name           string
}

func (lp *LocalPlayer) Move(newX, newY int32) {
	lp.X = newX
	lp.Y = newY

	lp.ChunkX = lp.X / 16
	lp.ChunkY = lp.Y / 16
}

func (lp *LocalPlayer) SendMovePacket(game *Game, x, y int32) {
	_, exists := game.CollisionMap[util.Vector2I{X: x, Y: y}]
	if x > int32(game.MaxX) || x < 0 || y > int32(game.MaxY) || y < 0 || exists {
		return
	}

	SendPacket(game.Conn, &c2s.MovePacket{
		X: uint32(x),
		Y: uint32(y),
	})
}

func (lp *LocalPlayer) Update(game *Game) {
}

type RemotePlayer struct {
	X, Y         int32
	PrevX, PrevY int32
	RealX, RealY int32
	Facing       Direction
	Name         string
}

func NewRemotePlayer(x, y int32, facing Direction, name string, game *Game) RemotePlayer {
	return RemotePlayer{
		X:      x,
		Y:      y,
		RealX:  (x % 16) * game.TileSize,
		RealY:  (y % 16) * game.TileSize,
		Facing: facing,
		Name:   name,
	}
}
