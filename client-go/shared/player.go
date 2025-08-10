package shared

import (
	"client/network/c2s"
	"client/util"
	"fmt"
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

func (lp *LocalPlayer) Move(newX, newY int32, facing Direction) {
	lp.X = newX
	lp.Y = newY
	lp.Facing = facing

	lp.ChunkX = lp.X / 16
	lp.ChunkY = lp.Y / 16
}

func (lp *LocalPlayer) SendMovePacket(game *Game, x, y int32, facing Direction) {
	if facing > 3 {
		return
	}

	_, exists := game.CollisionMap[util.Vector2I{X: x, Y: y}]
	if x > int32(game.MaxX) || x < 0 || y > int32(game.MaxY) || y < 0 || exists {
		if facing != lp.Facing {
			SendPacket(game.Conn, &c2s.MovePacket{
				X:      uint32(lp.X),
				Y:      uint32(lp.Y),
				Facing: byte(facing),
			})
		}

		return
	}

	SendPacket(game.Conn, &c2s.MovePacket{
		X:      uint32(x),
		Y:      uint32(y),
		Facing: byte(facing),
	})
}

func (lp *LocalPlayer) GetFacingCoord() util.Vector2I {
	switch lp.Facing {
	case DOWN:
		return util.Vector2I{X: lp.X, Y: lp.Y + 1}
	case LEFT:
		return util.Vector2I{X: lp.X - 1, Y: lp.Y}
	case RIGHT:
		return util.Vector2I{X: lp.X + 1, Y: lp.Y}
	case UP:
		return util.Vector2I{X: lp.X, Y: lp.Y - 1}
	default:
		panic(fmt.Sprintf("unexpected shared.Direction: %#v", lp.Facing))
	}
}

func (lp *LocalPlayer) Update(game *Game, crossedZone bool) {
	targetX := (lp.X % 16) * game.TileSize
	targetY := (lp.Y % 16) * game.TileSize

	const speed = 16.0

	if crossedZone {
		lp.RealX = targetX
		lp.RealY = targetY
	} else {
		if lp.RealX < targetX {
			lp.RealX += speed
		} else if lp.RealX > targetX {
			lp.RealX -= speed
		}

		if lp.RealY < targetY {
			lp.RealY += speed
		} else if lp.RealY > targetY {
			lp.RealY -= speed
		}
	}

	lp.PrevX = lp.X
	lp.PrevY = lp.Y
}

type RemotePlayer struct {
	X, Y         int32
	PrevX, PrevY int32
	RealX, RealY int32
	Facing       Direction
	Name         string
}

func NewRemotePlayer(x, y int32, facing Direction, name string) *RemotePlayer {
	return &RemotePlayer{
		X:      x,
		Y:      y,
		Facing: facing,
		Name:   name,
	}
}

func (rp *RemotePlayer) Move(newX, newY int32, facing Direction) {
	rp.X = newX
	rp.Y = newY
	rp.Facing = facing
}

func (rp *RemotePlayer) Update(game *Game) {
	targetX := (rp.X % 16) * game.TileSize
	targetY := (rp.Y % 16) * game.TileSize

	// just logged in, basically.
	if rp.PrevX == 0 && rp.PrevY == 0 {
		rp.RealX = targetX
		rp.RealY = targetY

		rp.PrevX = rp.X
		rp.PrevY = rp.Y
		return
	}

	const speed = 16.0

	if rp.RealX < targetX {
		rp.RealX += speed
	} else if rp.RealX > targetX {
		rp.RealX -= speed
	}

	if rp.RealY < targetY {
		rp.RealY += speed
	} else if rp.RealY > targetY {
		rp.RealY -= speed
	}

	rp.PrevX = rp.X
	rp.PrevY = rp.Y
}
