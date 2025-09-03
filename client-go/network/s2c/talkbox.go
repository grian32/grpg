package s2c

import (
	"client/shared"
	"cmp"
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

	if tbType == PLAYER {
		content, err := buf.ReadString()
		if err != nil {
			log.Printf("couldn't read talkbox packet: %v\n", err)
		}
		game.Talkbox.CurrentMessage = content
		game.Talkbox.CurrentName = game.Player.Name
		game.Talkbox.Active = true
	} else if tbType == NPC {
		npcId, err1 := buf.ReadUint16()
		content, err2 := buf.ReadString()

		if err := cmp.Or(err1, err2); err != nil {
			log.Printf("couldn't read talkbox packet: %v\n", err)
		}

		game.Talkbox.CurrentName = game.Npcs[npcId].Name
		game.Talkbox.CurrentMessage = content
		game.Talkbox.Active = true
	} else {
		game.Talkbox.CurrentName = ""
		game.Talkbox.CurrentMessage = ""
		game.Talkbox.Active = false
	}
}
