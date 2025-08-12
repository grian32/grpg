package s2c

import (
	"grpg/data-go/gbuf"
	"server/shared"
	"server/util"
)

type ObjUpdate struct {
	ChunkPos util.Vector2I
}

func (o *ObjUpdate) Opcode() byte {
	return 0x04
}

func (o *ObjUpdate) Handle(buf *gbuf.GBuf, game *shared.Game) {
	packetLen := 2 // amount of objs
	objLen := 0

	for _, obj := range game.TrackedObjs {
		if obj.ChunkPos == o.ChunkPos {
			packetLen += 4 + 4 + 2 + 1 // x, y, objid, state
			objLen++
		}
	}

	buf.WriteUint16(uint16(packetLen))
	buf.WriteUint16(uint16(objLen))

	// TODO: look into if it's possible to only send objid in the case of a rebuild, seems dodgy tho, ngl :S

	for pos, obj := range game.TrackedObjs {
		if obj.ChunkPos == o.ChunkPos {
			buf.WriteUint32(pos.X)
			buf.WriteUint32(pos.Y)
			buf.WriteUint16(obj.ObjData.ObjId)
			buf.WriteByte(obj.State)
		}
	}
}
