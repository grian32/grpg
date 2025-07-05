package s2c

import (
	"client/game"
	"client/shared"
	"grpg/data-go/gbuf"
)

type LoginAccepted struct{}

func (l *LoginAccepted) Handle(buf *gbuf.GBuf, g *shared.Game) {
	initialX, _ := buf.ReadInt32()
	initialY, _ := buf.ReadInt32()

	g.Player.Move(initialX, initialY, g)
	g.SceneManager.SwitchTo(&game.Playground{Game: g})
	// this isn't *really* necessary since login screen switches off once login is accepted but doing it for
	// consistency or if i allow log out later or something else.
	g.ShowFailedLogin = false
}
