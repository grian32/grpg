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
	TrackedObjs   map[util.Vector2I]*GameObj
	CollisionMap  map[util.Vector2I]struct{}
	ScriptManager *scripts.ScriptManager
	CurrentTick   uint32
	TimedScripts  map[uint32][]scripts.TimedScript
}

func (g *Game) AddTimedScript(tick uint32, script scripts.TimedScript) {
	_, ok := g.TimedScripts[tick]
	if !ok {
		g.TimedScripts[tick] = []scripts.TimedScript{script}
	} else {
		g.TimedScripts[tick] = append(g.TimedScripts[tick], script)
	}
}
