package shared

import (
	"net"
	"server/util"
)

type Player struct {
	Pos util.Vector2I
	// might not need these will see how design pans out
	ChunkPos util.Vector2I
	Name     string
	Conn     net.Conn
}
