package shared

import (
	"database/sql"
	"server/util"
)

type Game struct {
	Players      []*Player
	MaxX         uint32
	MaxY         uint32
	Database     *sql.DB
	CollisionMap map[util.Vector2I]struct{}
}
