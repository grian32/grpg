package shared

import (
	"database/sql"
	"net"
	"server/util"
)

type Game struct {
	Players      map[*Player]struct{}
	Connections  map[net.Conn]*Player
	MaxX         uint32
	MaxY         uint32
	Database     *sql.DB
	TrackedObjs  map[util.Vector2I]*GameObj
	TrackedNpcs  map[util.Vector2I]*GameNpc
	CollisionMap map[util.Vector2I]struct{}
	CurrentTick  uint32
}
