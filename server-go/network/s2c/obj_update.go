package s2c

import (
	"grpg/data-go/gbuf"
	"server/shared"
	"server/util"
)

type ObjUpdate struct {
	ChunkPos util.Vector2I
	Rebuild  bool
}

func (o *ObjUpdate) Opcode() byte {
	return 0x04
}

func (o *ObjUpdate) Handle(buf *gbuf.GBuf, game *shared.Game) {
	packetLen := 1 + 2 // rebuild boolean + amount of objs
	objLen := 0

	for _, obj := range game.TrackedObjs {
		if obj.ChunkPos == o.ChunkPos {
			packetLen += 4 + 4 + 1 // x, y, state
			if o.Rebuild {
				packetLen += 2 // objid
			}
			objLen++
		}
	}

	buf.WriteUint16(uint16(packetLen))
	buf.WriteBool(o.Rebuild)
	buf.WriteUint16(uint16(objLen))

	for pos, obj := range game.TrackedObjs {
		if obj.ChunkPos == o.ChunkPos {
			buf.WriteUint32(pos.X)
			buf.WriteUint32(pos.Y)
			if o.Rebuild {
				buf.WriteUint16(obj.ObjData.ObjId)
			}
			buf.WriteByte(obj.State)
		}
	}
}
