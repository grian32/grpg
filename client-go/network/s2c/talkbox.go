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

	var msg, name string

	switch tbType {
	case PLAYER:
		content, err := buf.ReadString()
		if err != nil {
			log.Printf("couldn't read talkbox packet: %v\n", err)
		}
		msg = content
		name = game.Player.Name
	case NPC:
		npcId, err1 := buf.ReadUint16()
		content, err2 := buf.ReadString()

		if err := cmp.Or(err1, err2); err != nil {
			log.Printf("couldn't read talkbox packet: %v\n", err)
		}
		msg = content
		name = game.Npcs[npcId].Name
	case CLEAR:
		game.Talkbox.CurrentName = ""
		game.Talkbox.CurrentMessage = ""
		game.Talkbox.Active = false
		return
	}

	game.Talkbox.CurrentName = name
	game.Talkbox.CurrentMessage = msg
	game.Talkbox.Active = true
}
