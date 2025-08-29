package grpgnpc

import (
	"cmp"
	"grpg/data-go/gbuf"
)

type Header struct {
	Magic [8]byte
}

type Npc struct {
	NpcId     uint16
	Name      string
	TextureId uint16
}

func WriteHeader(buf *gbuf.GBuf) {
	buf.WriteBytes([]byte("GRPGNPC\x00"))
}

func ReadHeader(buf *gbuf.GBuf) (Header, error) {
	bytes, err := buf.ReadBytes(8)
	if err != nil {
		return Header{}, err
	}

	return Header{
		Magic: [8]byte(bytes),
	}, nil
}

func WriteNpcs(buf *gbuf.GBuf, npcs []Npc) {
	buf.WriteUint16(uint16(len(npcs)))

	for _, npc := range npcs {
		buf.WriteUint16(npc.NpcId)
		buf.WriteString(npc.Name)
		buf.WriteUint16(npc.TextureId)
	}
}

func ReadNpcs(buf *gbuf.GBuf) ([]Npc, error) {
	npcLen, err := buf.ReadUint16()
	if err != nil {
		return nil, err
	}

	npcs := make([]Npc, npcLen)

	for idx := range npcLen {
		id, err1 := buf.ReadUint16()
		name, err2 := buf.ReadString()
		textureId, err3 := buf.ReadUint16()

		if err := cmp.Or(err1, err2, err3); err != nil {
			return nil, err
		}
		npcs[idx] = Npc{
			NpcId:     id,
			Name:      name,
			TextureId: textureId,
		}
	}

	return npcs, nil
}
