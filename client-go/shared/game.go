package shared

import (
	"client/constants"
	"client/network/c2s"
	"client/util"
	"grpg/data-go/gbuf"
	"grpg/data-go/grpgitem"
	"grpg/data-go/grpgnpc"
	"grpg/data-go/grpgobj"
	"grpg/data-go/grpgtile"
	"net"
)

type Game struct {
	ScreenWidth  int32
	ScreenHeight int32
	ScreenRatio  float64
	MaxX         uint16
	MaxY         uint16
	CollisionMap map[util.Vector2I]struct{}
	Objs         map[uint16]*grpgobj.Obj
	Npcs         map[uint16]*grpgnpc.Npc
	// this is literally only needed to send the right obj id with the interact packet, only stores stateful packets
	ObjIdByLoc  map[util.Vector2I]uint16
	Tiles       map[uint16]*grpgtile.Tile
	Items       map[uint16]grpgitem.Item
	TrackedObjs map[util.Vector2I]*GameObj
	NpcsByPos   map[util.Vector2I]*GameNpc
	// uid to gamenpc
	TrackedNpcs         map[uint32]*GameNpc
	Skills              map[Skill]*SkillInfo
	SkillHoverMsgs      map[Skill]*string
	TileSize            int32
	SceneManager        *GSceneManager
	Player              *LocalPlayer
	Talkbox             Talkbox
	OtherPlayers        map[string]*RemotePlayer
	PlayerVars          map[constants.PlayerVarId]uint16
	PlayerVarHandlers   map[constants.PlayerVarId]PlayerVarHandlerFunc
	Conn                net.Conn
	OutlineInvSpot      int
	ShowFailedLogin     bool
	DebugMode           bool
	RenderExclamOnGuide bool
}

// i think this would make sense as a function on game but er.. cyclical lol!
func SendPacket(conn net.Conn, packet c2s.Packet) {
	buf := gbuf.NewEmptyGBuf()
	buf.WriteByte(packet.Opcode())
	packet.Handle(buf)
	_, err := conn.Write(buf.Bytes())
	if err != nil {
		return
	}
}
