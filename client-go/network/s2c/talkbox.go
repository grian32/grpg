package s2c

import (
	"client/shared"
	"fmt"
	"grpg/data-go/gbuf"
	"log"
)

type TalkboxType byte

const (
	PLAYER TalkboxType = iota
	NPC
	CLEAR
)

type Talkbox struct{}

func (t *Talkbox) Handle(buf *gbuf.GBuf, game *shared.Game) {
	tType, err := buf.ReadByte()

	if err != nil {
		log.Printf("couldn't read talkbox packet: %v\n", err)
	}
	tbType := TalkboxType(tType)

	if tbType != CLEAR {
		content, err := buf.ReadString()
		if err != nil {
			log.Printf("couldn't read talkbox packet: %v\n", err)
		}
		game.Talkbox.CurrentMessage = content
		game.Talkbox.Active = true
		fmt.Printf("got talkbox pkt: %v\n", game.Talkbox)
	} else {
		game.Talkbox.CurrentMessage = ""
		game.Talkbox.Active = false
	}

	// TODO: i need to transfer the npc id also
	_ = tType
}
