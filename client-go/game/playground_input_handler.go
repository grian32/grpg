package game

import (
	"client/network/c2s"
	"client/shared"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type PgInputHandler struct {
	Player *shared.LocalPlayer
	Game   *shared.Game

	IsTypingCommand bool
	CommandString   string
}

func NewPgInputHandler(g *shared.Game) *PgInputHandler {
	return &PgInputHandler{Game: g, Player: g.Player}
}

func (h *PgInputHandler) Update() {
	if h.IsTypingCommand {
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			h.IsTypingCommand = false
			h.CommandString = ""
		} else if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			h.IsTypingCommand = false
			h.Player.SendCmdPacket(h.Game, h.CommandString)
			h.CommandString = ""
		} else if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) && len(h.CommandString) != 0 {
			h.CommandString = h.CommandString[:len(h.CommandString)-1]
		} else {
			h.CommandString = string(ebiten.AppendInputChars([]rune(h.CommandString)))
		}
	} else {
		if inpututil.IsKeyJustPressed(ebiten.KeyW) {
			h.Player.SendMovePacket(h.Game, h.Player.X, h.Player.Y-1, shared.UP)
		} else if inpututil.IsKeyJustPressed(ebiten.KeyS) {
			h.Player.SendMovePacket(h.Game, h.Player.X, h.Player.Y+1, shared.DOWN)
		} else if inpututil.IsKeyJustPressed(ebiten.KeyA) {
			h.Player.SendMovePacket(h.Game, h.Player.X-1, h.Player.Y, shared.LEFT)
		} else if inpututil.IsKeyJustPressed(ebiten.KeyD) {
			h.Player.SendMovePacket(h.Game, h.Player.X+1, h.Player.Y, shared.RIGHT)
		} else if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
			h.Player.SendInteractPacket(h.Game)
		} else if h.Game.Talkbox.Active && inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			shared.SendPacket(h.Game.Conn, &c2s.Continue{})
		} else if inpututil.IsKeyJustPressed(ebiten.KeyC) {
			h.IsTypingCommand = true
		}
	}
}
