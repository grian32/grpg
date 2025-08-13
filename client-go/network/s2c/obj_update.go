package s2c

import (
	"client/shared"
	"client/util"
	"cmp"
	"fmt"
	"grpg/data-go/gbuf"
)

type ObjUpdate struct {
}

func (o *ObjUpdate) Handle(buf *gbuf.GBuf, game *shared.Game) {
	objLen, err := buf.ReadUint16()
	if err != nil {
		fmt.Printf("failed to read uint16 len in obj update: %v", err)
		return
	}

	// just forced to rebuild as i don't really have a way of determining if the player has crossed a zone at packet
	// processing time
	objMap := make(map[util.Vector2I]shared.GameObj)

	for range objLen {
		x, err1 := buf.ReadUint32()
		y, err2 := buf.ReadUint32()
		objId, err3 := buf.ReadUint16()
		state, err4 := buf.ReadByte()
		if err := cmp.Or(err1, err2, err3, err4); err != nil {
			fmt.Printf("failed to read obj in obj update: %v", err)
			return
		}

		pos := util.Vector2I{X: int32(x), Y: int32(y)}

		objMap[pos] = shared.GameObj{
			DataObj: game.Objs[objId],
			State:   state,
		}
	}

	game.TrackedObjs = objMap
}
