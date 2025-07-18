package s2c

import (
	"grpg/data-go/gbuf"
	"log"
	"server/shared"
	"server/util"
)

type PlayersUpdate struct {
	ChunkPos util.Vector2I
}

func (p *PlayersUpdate) Opcode() byte {
	return 0x03
}

func (p *PlayersUpdate) Handle(buf *gbuf.GBuf, game *shared.Game) {
	packetLen := 2 // player len
	playerLen := 0

	// TODO: might be able to do this without iterating twice but would require a way to modify at pos in gbuf
	for _, player := range game.Players {
		if player.ChunkPos == p.ChunkPos {
			packetLen += 4 + len(player.Name) + 4 + 4 // len name, name, x, y
			playerLen++
		}
	}

	// this will also catch name len being > uint32 & len(players) being > uint16 since packetlen includes them
	if packetLen > 65535 {
		log.Printf("Couldn't send update packet due to too high packet length. %v", p.ChunkPos)
	}

	buf.WriteUint16(uint16(packetLen))

	buf.WriteUint16(uint16(playerLen))

	for _, player := range game.Players {
		if player.ChunkPos == p.ChunkPos {
			buf.WriteString(player.Name)
			buf.WriteUint32(player.Pos.X)
			buf.WriteUint32(player.Pos.Y)
		}
	}
}
