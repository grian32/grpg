package shared

import (
	"client/network/c2s"
	"client/util"
	"log"
	"time"
)

type LocalPlayer struct {
	X, Y           int32
	PrevX, PrevY   int32
	RealX, RealY   int32
	ChunkX, ChunkY int32
	Facing         Direction
	CurrFrame      uint8
	FrameCounter   float64
	Inventory      [24]InventoryItem
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

func (lp *LocalPlayer) SendCmdPacket(game *Game, cmd string) {
	if cmd == "debug" {
		game.DebugMode = !game.DebugMode
		return
	}

	SendPacket(game.Conn, &c2s.Command{
		Msg: cmd,
	})
}

func (lp *LocalPlayer) SendItemUsePacket(game *Game, invIdx uint8) {
	if invIdx > 23 {
		return
	}
	item := lp.Inventory[invIdx]
	if item.ItemId == 0 {
		return
	}
	SendPacket(game.Conn, &c2s.ItemUse{
		InvIdx: invIdx,
	})
}

// SendInteractPacket TODO: maybe bad place for this?
func (lp *LocalPlayer) SendInteractPacket(game *Game) {
	facing := lp.GetFacingCoord()
	pos := util.Vector2I{X: facing.X, Y: facing.Y}

	if objId, ok := game.ObjIdByLoc[pos]; ok {
		SendPacket(game.Conn, &c2s.InteractPacket{
			ObjId: objId,
			X:     uint32(facing.X),
			Y:     uint32(facing.Y),
		})
		return
	}

	if npc, ok := game.NpcsByPos[pos]; ok {
		SendPacket(game.Conn, &c2s.TalkPacket{
			NpcId: npc.NpcData.NpcId,
			Uid:   npc.Uid,
		})
		return
	}

	log.Printf("warn: interact/talk packet tried to find obj/npc that does not exist @ %d, %d\n", facing.X, facing.Y)
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
		log.Fatalf("unexpected shared.Direction: %#v", lp.Facing)
	}

	return util.Vector2I{}
}

func (lp *LocalPlayer) Update(game *Game, crossedZone bool, movementHeld bool) {
	targetX := (lp.X % 16) * game.TileSize
	targetY := (lp.Y % 16) * game.TileSize

	const speed = 16.0

	if lp.PrevX == 0 && lp.PrevY == 0 {
		lp.RealX = targetX
		lp.RealY = targetY

		lp.PrevX = lp.X
		lp.PrevY = lp.Y

		lp.CurrFrame = 0
		lp.FrameCounter = 0
		return
	}

	if crossedZone {
		lp.RealX = targetX
		lp.RealY = targetY

		lp.CurrFrame = 0
		lp.FrameCounter = 0
	} else if lp.RealX != targetX || lp.RealY != targetY {
		lp.FrameCounter++
		lp.CurrFrame = uint8(lp.FrameCounter/8) % 4

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
	} else if movementHeld {
		lp.FrameCounter++
		lp.CurrFrame = uint8(lp.FrameCounter/8) % 4
	} else {
		lp.CurrFrame = 0
		lp.FrameCounter = 0
	}

	lp.PrevX = lp.X
	lp.PrevY = lp.Y
}

type RemotePlayer struct {
	X, Y         int32
	PrevX, PrevY int32
	RealX, RealY int32
	Facing       Direction
	CurrFrame    uint8
	FrameCounter uint64
	LastMoveTime time.Time
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
	// FIXME: really dubious hack to be honest, im not too happy with it, but can't really think of another good way to do it
	if rp.X != newX || rp.Y != newY {
		rp.LastMoveTime = time.Now()
	}

	rp.X = newX
	rp.Y = newY
	rp.Facing = facing
}

func (rp *RemotePlayer) Update(game *Game) {
	targetX := (rp.X % 16) * game.TileSize
	targetY := (rp.Y % 16) * game.TileSize
	isMoving := time.Since(rp.LastMoveTime) < 200*time.Millisecond

	// just logged in, basically.
	if rp.PrevX == 0 && rp.PrevY == 0 {
		rp.RealX = targetX
		rp.RealY = targetY

		rp.PrevX = rp.X
		rp.PrevY = rp.Y

		rp.CurrFrame = 0
		rp.FrameCounter = 0
		return
	}

	const speed = 16.0

	if rp.RealX != targetX || rp.RealY != targetY {
		rp.FrameCounter++
		rp.CurrFrame = uint8(rp.FrameCounter/8) % 4

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
	} else if isMoving {
		rp.FrameCounter++
		rp.CurrFrame = uint8(rp.FrameCounter/8) % 4
	} else {
		rp.CurrFrame = 0
		rp.FrameCounter = 0
	}

	rp.PrevX = rp.X
	rp.PrevY = rp.Y
}
