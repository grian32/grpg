package s2c

import (
	"client/shared"
	"client/util"
	"cmp"
	"grpg/data-go/gbuf"
	"log"
)

type ObjUpdate struct {
}

func (o *ObjUpdate) Handle(buf *gbuf.GBuf, game *shared.Game) {
	rebuild, err := buf.ReadBool()
	if err != nil {
		log.Printf("failed to read rebuild bool in obj update: %v\n", err)
	}
	objLen, err := buf.ReadUint16()
	if err != nil {
		log.Printf("failed to read uint16 len in obj update: %v\n", err)
		return
	}

	if rebuild {
		objMap := make(map[util.Vector2I]*shared.GameObj)

		for range objLen {
			x, err1 := buf.ReadUint32()
			y, err2 := buf.ReadUint32()
			objId, err3 := buf.ReadUint16()
			state, err4 := buf.ReadByte()
			if err := cmp.Or(err1, err2, err3, err4); err != nil {
				log.Printf("failed to read obj in obj update: %v", err)
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
			state, err3 := buf.ReadByte()

			if err := cmp.Or(err1, err2, err3); err != nil {
				log.Printf("failed to read obj in obj update: %v", err)
				return
			}

			pos := util.Vector2I{
				X: int32(x),
				Y: int32(y),
			}

			obj, ok := game.TrackedObjs[pos]
			if !ok {
				log.Printf("tried to modify obj on rebuild that did not exist in tracked objs")
				return
			}
			obj.State = state
		}
	}
}
