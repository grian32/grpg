package s2c

import (
	"client/shared"
	"grpg/data-go/gbuf"
)

type LoginRejected struct{}

func (l LoginRejected) Handle(buf *gbuf.GBuf, game *shared.Game) {
}
