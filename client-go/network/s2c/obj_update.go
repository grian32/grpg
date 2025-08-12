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

	p := game.Player

	// either just logged in and no prev x/y or crossed a zone
	rebuild := p.PrevX == 0 && p.PrevY == 0 || (p.PrevX/16 != p.ChunkX || p.PrevY/16 != p.ChunkY)

	if rebuild {
		objMap := make(map[util.Vector2I]*shared.GameObj)

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

			objMap[pos] = &shared.GameObj{
				DataObj: game.Objs[objId],
				State:   state,
			}
		}

		game.TrackedObjs = objMap
	} else {
		for range objLen {
			x, err1 := buf.ReadUint32()
			y, err2 := buf.ReadUint32()
			_, err3 := buf.ReadUint16()
			state, err4 := buf.ReadByte()
			if err := cmp.Or(err1, err2, err3, err4); err != nil {
				fmt.Printf("failed to read obj in obj update: %v", err)
				return
			}

			pos := util.Vector2I{X: int32(x), Y: int32(y)}

			obj, exists := game.TrackedObjs[pos]
			if !exists {
				fmt.Printf("didn't require rebuild but obj doesn't exist in obj update, aborting packet")
				return
			}
			obj.State = state
		}
	}
}
