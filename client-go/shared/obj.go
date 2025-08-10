package shared

import "grpg/data-go/grpgobj"

type GameObj struct {
	DataObj grpgobj.Obj
	State   uint16 // if applicable
}
