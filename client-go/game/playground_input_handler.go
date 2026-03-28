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

	minInvX, maxInvX, minInvY, maxInvY int
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
	h.minInvX = RightGameframeX + TileSize
	h.maxInvX = h.minInvX + TileSize*4
	h.minInvY = TileSize
	h.maxInvY = h.minInvY + TileSize*6
}

func (h *PgInputHandler) UpdateItemMove(renderType RenderType, outlineInvSpot *int) {
	if renderType == Inventory && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mouseX, mouseY := ebiten.CursorPosition()
		if mouseX >= h.minInvX && mouseX < h.maxInvX && mouseY >= h.minInvY && mouseY < h.maxInvY {
			col := (mouseX - h.minInvX) / TileSize
			row := (mouseY - h.minInvY) / TileSize
			idx := row*ItemsPerRow + col

			if *outlineInvSpot != -1 {
				if *outlineInvSpot == idx || h.Player.Inventory[idx].ItemId != 0 {
					// deselect behaviour, basically
					*outlineInvSpot = -1
					return
				}

				shared.SendPacket(h.Game.Conn, &c2s.InvSwap{
					From: byte(*outlineInvSpot),
					To:   byte(idx),
				})
				*outlineInvSpot = -1

				return
			}

			if h.Player.Inventory[idx].ItemId != 0 {
				*outlineInvSpot = idx
			}
		}
	}
}
