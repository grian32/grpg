package shared

import (
	"client/network/c2s"
	"client/util"
	"grpg/data-go/gbuf"
	"grpg/data-go/grpgobj"
	"grpg/data-go/grpgtile"
	"net"
)

type Game struct {
	ScreenWidth  int32
	ScreenHeight int32
	MaxX         uint16
	MaxY         uint16
	CollisionMap map[util.Vector2I]struct{}
	Objs         map[uint16]*grpgobj.Obj
	// this is literally only needed to send the right obj id with the interact packet, only stores stateful packets
	ObjIdByLoc      map[util.Vector2I]uint16
	Tiles           map[uint16]*grpgtile.Tile
	TrackedObjs     map[util.Vector2I]GameObj
	TileSize        int32
	SceneManager    *GSceneManager
	Player          *LocalPlayer
	OtherPlayers    map[string]*RemotePlayer
	Conn            net.Conn
	ShowFailedLogin bool
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
