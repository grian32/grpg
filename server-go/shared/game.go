package shared

import (
	"database/sql"
	"net"
	"server/scripts"
	"server/util"
)

type Game struct {
	Players       map[*Player]struct{}
	Connections   map[net.Conn]*Player
	MaxX          uint32
	MaxY          uint32
	Database      *sql.DB
	CollisionMap  map[util.Vector2I]struct{}
	ScriptManager *scripts.ScriptManager
}
