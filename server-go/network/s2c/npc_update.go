package s2c

import (
	"grpg/data-go/gbuf"
	"server/shared"
	"server/util"
)

type NpcUpdate struct {
	ChunkPos util.Vector2I
}

func (n *NpcUpdate) Opcode() byte {
	return 0x06
}

func (n *NpcUpdate) Handle(buf *gbuf.GBuf, game *shared.Game) {
	packetLen := 2
	npcLen := 0

	for _, npc := range game.TrackedNpcs {
		if npc.ChunkPos == n.ChunkPos {
			packetLen += 4 + 4 + 2 // x y id
		}
		npcLen++
	}

	buf.WriteUint16(uint16(packetLen))
	buf.WriteUint16(uint16(npcLen))

	for pos, npc := range game.TrackedNpcs {
		if npc.ChunkPos == n.ChunkPos {
			buf.WriteUint32(pos.X)
			buf.WriteUint32(pos.Y)
			buf.WriteUint16(npc.NpcData.NpcId)
		}
	}
}
