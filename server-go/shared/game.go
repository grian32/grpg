package shared

import "server/util"

type Game struct {
	Players      []*Player
	MaxX         uint32
	MaxY         uint32
	CollisionMap map[util.Vector2I]struct{}
}
