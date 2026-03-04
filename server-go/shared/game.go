package shared

import (
	"database/sql"
	"net"
	"server/util"
	"sync"
)

type Game struct {
	Players      map[*Player]struct{}
	Connections  map[net.Conn]*Player
	MaxX         uint32
	MaxY         uint32
	Database     *sql.DB
	TrackedObjs  map[util.Vector2I]*GameObj
	TrackedNpcs  map[util.Vector2I]*GameNpc
	Mu           sync.RWMutex
	CollisionMap map[util.Vector2I]struct{}
	// currently practical duplicate of CollisionMap but CollisionMap will presumably have other sutff in the future, otherwise easy to remove
	Objs         map[util.Vector2I]struct{}
	TimedScripts map[uint32][]func()
	// so this is : chunk to path and then each path actually contains the moves themselves, so for each chunk itll pop one from each list and send it with the packet
	NpcMoves    map[util.Vector2I][]NpcPath
	CurrentTick uint32
}
